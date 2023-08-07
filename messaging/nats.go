package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/chaos-io/chaos/logs"
)

type Nats struct {
	conn       *nats.Conn
	js         nats.JetStreamContext
	Config     *Config
	streamName string
	subjects   []string
	shutdown   bool
	shutdownCh chan struct{}
}

func New(cfg *Config) *Nats {
	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		panic(fmt.Errorf("failed to connect nats, error: %v", err))
	}

	// Create a JetStream management interface
	js, err := nc.JetStream()
	if err != nil {
		panic(fmt.Errorf("failed to create jetstream, error: %v", err))
	}

	n := &Nats{
		js:         js,
		conn:       nc,
		Config:     cfg,
		streamName: cfg.StreamName,
		subjects:   cfg.TopicNames,
		shutdownCh: make(chan struct{}),
	}
	if len(n.subjects) == 0 {
		n.subjects = []string{n.streamName + ".*"}
	}

	info, err := n.js.StreamInfo(n.streamName)
	if err != nil && !errors.Is(err, nats.ErrStreamNotFound) {
		logs.Warnw("failed to get the stream info", "streamName", n.streamName, "error", err)
		return nil
	}

	if info == nil {
		err = n.createStream(n.streamName, n.subjects)
	} else {
		err = n.updateStream(n.streamName, info.Config.Subjects, n.subjects)
	}
	if err != nil {
		return nil
	}

	return n
}

func (n *Nats) createStream(name string, subjects []string) error {
	if _, err := n.js.AddStream(&nats.StreamConfig{Name: name, Subjects: subjects}); err != nil {
		logs.Warnw("failed to create stream", "error", err)
		return err
	}
	logs.Infow("create stream", "name", name, "subjects", subjects)
	return nil
}

func (n *Nats) updateStream(name string, originSubjects, subjects []string) error {
	mergeSubjects := append(originSubjects, subjects...)
	freqMap := make(map[string]struct{})
	allSubjects := make([]string, 0, len(originSubjects)+len(subjects))
	for _, s := range mergeSubjects {
		if _, ok := freqMap[s]; !ok {
			freqMap[s] = struct{}{}
			allSubjects = append(allSubjects, s)
		}
	}

	if _, err := n.js.UpdateStream(&nats.StreamConfig{Name: name, Subjects: allSubjects}); err != nil {
		logs.Warnw("failed to update stream", "streamName", name, "subjects", allSubjects, "error", err)
		return err
	}

	logs.Infow("update stream", "name", name, "subjects", allSubjects)

	return nil
}

func (n *Nats) Publish(subj string, messages ...Message) error {
	for _, message := range messages {
		bytes, err := json.Marshal(message)
		if err != nil {
			return logs.NewErrorw("marshal message error", "message", message, "error", err)
		}

		if err = n.publish(subj, bytes); err != nil {
			return logs.NewErrorw("publish message error", "subj", subj, "error", err)
		}
	}

	logs.Infow("published message", "subj", subj)

	return nil
}

func (n *Nats) publish(subject string, msg []byte) error {
	_, err := n.js.Publish(subject, msg)
	return err
}

type Handler func(ctx context.Context, subscription *Subscription, m *SubMessage) error

func (n *Nats) Subscribe(s *Subscription, h Handler) {
	natsCB := func(m *nats.Msg) {
		msg := &SubMessage{}
		if err := json.Unmarshal(m.Data, msg); err != nil {
			logs.Warnw("failed to unmarshal data form subject", "subject", s.Subject, "error", err)
			return
		}

		msg.SetAck(func() { m.Ack() })
		msg.SetNak(func() { m.Nak() })
		msg.SetTerm(func() { m.Term() })
		msg.SetInProgress(func() { m.InProgress() })

		if err := h(context.Background(), s, msg); err != nil {
			return
		}

		if s.AutoAck {
			msg.Ack()
		}
	}

	if len(s.Queue) > 0 {
		_, err := n.queueSubscribe(s, natsCB)
		if err != nil {
			logs.Warnw("failed to to subscribe with queue", "subject", s.Subject, "queue", s.Queue)
		}
	} else {
		_, err := n.subscribe(s, natsCB)
		if err != nil {
			logs.Warnw("failed to to subscribe with queue", "subject", s.Subject, "queue", s.Queue)
		}
	}
}

func (n *Nats) subscribe(s *Subscription, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	return n.js.Subscribe(s.Subject, cb, opts...)
}

func (n *Nats) queueSubscribe(s *Subscription, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	if s.Pull {
		return n.pullSubscribe(s, cb, opts...)
	} else {
		return n.js.QueueSubscribe(s.Subject, s.Queue, cb, opts...)
	}
}

func (n *Nats) pullSubscribe(s *Subscription, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	opts = append(opts, nats.PullMaxWaiting(PullMaxWait))
	if s.AckWait > 0 {
		opts = append(opts, nats.SubOpt(nats.AckWait(time.Duration(s.AckWait)*time.Second)))
	}

	subs, err := n.js.PullSubscribe(s.Subject, s.Durable, opts...)
	if err != nil {
		logs.Warnw("failed to pull subscribe", "subject", s.Subject, "durable", s.Durable, "error", err)
		return nil, err
	}

	go n.Fetch(subs, s.Durable, cb)

	return subs, nil
}

func (n *Nats) Fetch(subs *nats.Subscription, durable string, cb nats.MsgHandler) {
	for {
		select {
		case <-n.shutdownCh:
			logs.Infow("shutdown the nats client, will closed the pull subscribe", "durable", durable)
			return
		default:
		}

		msgs, err := subs.Fetch(1)
		if err != nil {
			logs.Debugw("fetch the message error", "durable", durable)
		}

		for _, msg := range msgs {
			logs.Infow("fetched the message", "durable", durable, "subject", msg.Subject)
			cb(msg)
		}
	}
}

// Shutdown shuts down all subscribers
func (n *Nats) Shutdown() {
	n.conn.Close()
	close(n.shutdownCh)
	n.shutdown = true
	return
}
