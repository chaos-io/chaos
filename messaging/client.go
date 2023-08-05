package messaging

import (
	"encoding/json"
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
		subjects:   cfg.Subjects,
		shutdownCh: make(chan struct{}),
	}
	if len(n.subjects) == 0 {
		n.subjects = []string{n.streamName + ".*"}
	}

	info, err := n.js.StreamInfo(n.streamName)
	if err != nil {
		logs.Warnw("failed to get the stream info", "error", err)
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
		logs.Warnw("create stream failed", "error", err)
		return err
	}
	logs.Info("create stream", "name", name, "subjects", subjects)
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
		logs.Warnw("update stream failed", "error", err)
		return err
	}

	logs.Infow("update stream", "name", name, "subjects", allSubjects)

	return nil
}

func (n *Nats) Publish(subject string, messages ...Message) error {
	for _, message := range messages {
		bytes, err := json.Marshal(message)
		if err != nil {
			return logs.NewErrorw("marshal message error", "message", message, "error", err)
		}

		if err = n.publish(subject, bytes); err != nil {
			return logs.NewErrorw("publish message error", "subject", subject, "error", err)
		}
	}

	logs.Infow("published message", "subject", subject)

	return nil
}

func (n *Nats) publish(subject string, msg []byte) error {
	_, err := n.js.Publish(subject, msg)
	return err
}

func (n *Nats) Subscribe(s *Subscribe, h nats.Handler) error {
	if len(s.Queue) > 0 {
		_, err := n.queueSubscribe(s, nil)
		if err != nil {
			logs.Warnw("failed to to subscribe with queue", "subject", s.Subject, "queue", s.Queue)
		}
	} else {
		_, err := n.subscribe(s, nil)
		if err != nil {
			logs.Warnw("failed to to subscribe with queue", "subject", s.Subject, "queue", s.Queue)
		}
	}
	return nil
}

func (n *Nats) subscribe(s *Subscribe, msgH nats.MsgHandler) (*nats.Subscription, error) {
	return n.js.Subscribe(s.Subject, msgH)
}

func (n *Nats) queueSubscribe(s *Subscribe, msgH nats.MsgHandler) (*nats.Subscription, error) {
	if s.Pull {
		return n.pullSubscribe(s, msgH)
	} else {
		return n.js.QueueSubscribe(s.Subject, s.Queue, msgH)
	}
}

func (n *Nats) pullSubscribe(s *Subscribe, msgH nats.MsgHandler) (*nats.Subscription, error) {
	subOpts := []nats.SubOpt{nats.PullMaxWaiting(PullMaxWait)}
	if s.AckWait > 0 {
		subOpts = append(subOpts, nats.SubOpt(nats.AckWait(time.Duration(s.AckWait)*time.Second)))
	}

	subs, err := n.js.PullSubscribe(s.Subject, s.Durable)
	if err != nil {
		logs.Warnw("failed to pull subscribe", "subject", s.Subject, "error", err)
		return nil, err
	}

	go n.Fetch(subs, s.Durable, msgH)

	return subs, nil
}

func (n *Nats) Fetch(subs *nats.Subscription, durable string, msgH nats.MsgHandler) {
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
			msgH(msg)
		}
	}
}
