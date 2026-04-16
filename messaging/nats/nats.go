package nats

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"

	"github.com/chaos-io/chaos/logs"
	"github.com/chaos-io/chaos/messaging"
)

const defaultPullMaxWaiting = 128

var (
	registerOnce sync.Once
	registerErr  error
)

func Register() error {
	registerOnce.Do(func() {
		registerErr = messaging.Register("nats", func(cfg *messaging.Config) (messaging.Queue, error) {
			return New(cfg)
		})
	})
	return registerErr
}

func MustRegister() {
	if err := Register(); err != nil {
		panic(err)
	}
}

type Nats struct {
	id            string
	conn          *nats.Conn
	js            nats.JetStreamContext
	subscriptions map[string]*nats.Subscription
	shutdownCh    chan struct{}
	shutdownOnce  sync.Once
	subMu         sync.Mutex
	streamMu      sync.Mutex

	config     *messaging.Config
	pullNumber atomic.Int32
}

func New(cfg *messaging.Config) (*Nats, error) {
	if cfg == nil {
		return nil, messaging.ErrNilConfig
	}
	if len(strings.TrimSpace(cfg.ServiceName)) == 0 {
		return nil, errors.New("messaging nats: serviceName is empty")
	}

	nc, err := nats.Connect(cfg.ServiceName)
	if err != nil {
		logs.Errorw("failed to connect nats", "error", err)
		return nil, err
	}

	n := &Nats{
		id:            "",
		conn:          nc,
		config:        cfg,
		subscriptions: map[string]*nats.Subscription{},
		shutdownCh:    make(chan struct{}),
	}

	if cfg.Nats != nil && len(cfg.Nats.JetStream) > 0 {
		js, err := nc.JetStream()
		if err != nil {
			return nil, err
		}

		n.js = js

		if len(n.SubjectNames()) == 0 {
			n.SetSubjectNames(n.StreamName() + ".*")
		}

		n.createStream(n.StreamName(), n.SubjectNames())
	}

	return n, nil
}

func (n *Nats) createStream(name string, subjects []string) {
	if n == nil || n.js == nil || len(name) == 0 {
		return
	}

	streamInfo, err := n.js.StreamInfo(name)
	if err != nil {
		logs.Warnw("failed to get the stream info", "name", name, "error", err)
	}

	if streamInfo == nil {
		streamCfg := n.NewStreamConfig(name, subjects)
		if _, err := n.js.AddStream(streamCfg); err != nil {
			logs.Warnw("failed to create stream", "name", name, "error", err)
		} else {
			logs.Infof("creating stream %q and subject %q", name, subjects)
		}
	} else {
		if len(subjects) > 0 {
			allSubjects := mergeSubjects(streamInfo.Config.Subjects, subjects)
			streamConfig := n.NewStreamConfig(name, allSubjects)
			if _, err := n.js.UpdateStream(streamConfig); err != nil {
				logs.Warnw("failed to update stream", "name", name, "error", err)
			}
		}
	}
}

func mergeSubjects(existing []string, subjects []string) []string {
	if len(existing) == 0 {
		return append([]string{}, subjects...)
	}
	if len(subjects) == 0 {
		return append([]string{}, existing...)
	}

	set := make(map[string]struct{}, len(existing)+len(subjects))
	merged := make([]string, 0, len(existing)+len(subjects))
	for _, s := range existing {
		if len(s) == 0 {
			continue
		}
		if _, ok := set[s]; ok {
			continue
		}
		set[s] = struct{}{}
		merged = append(merged, s)
	}
	for _, s := range subjects {
		if len(s) == 0 {
			continue
		}
		if _, ok := set[s]; ok {
			continue
		}
		set[s] = struct{}{}
		merged = append(merged, s)
	}
	return merged
}

func (n *Nats) Publish(ctx context.Context, topic string, messages ...*messaging.Message) error {
	for _, message := range messages {
		bytes, err := jsoniter.Marshal(message)
		if err != nil {
			logs.Errorw("marshal message error", "message", message, "error", err)
			return err
		}

		if err = n.publish(ctx, topic, bytes); err != nil {
			return logs.NewErrorw("publish message error", "topic", topic, "error", err)
		}
	}

	logs.Infow("Nats: Published to queue", "topic", topic, "id", n.id)
	return nil
}

func (n *Nats) publish(ctx context.Context, topic string, msg []byte) error {
	_ = ctx

	if n.js != nil {
		n.updateSubjectName(topic)
		if _, err := n.js.Publish(topic, msg); err != nil {
			return err
		}
	} else {
		if err := n.conn.Publish(topic, msg); err != nil {
			return err
		}
	}

	return nil
}

func (n *Nats) updateSubjectName(name string) {
	if n == nil || n.js == nil || len(name) == 0 || len(n.StreamName()) == 0 {
		return
	}

	n.streamMu.Lock()
	defer n.streamMu.Unlock()

	subjects := n.SubjectNames()
	for _, s := range subjects {
		if s == name || s == n.StreamName()+".*" {
			return
		}
	}

	updatedSubjects := append(append([]string{}, subjects...), name)
	streamCfg := n.NewStreamConfig(n.StreamName(), updatedSubjects)
	if _, err := n.js.UpdateStream(streamCfg); err != nil {
		logs.Warnw("failed to update stream", "name", n.StreamName(), "subject", name, "error", err.Error())
		return
	}

	n.SetSubjectNames(updatedSubjects...)
}

func (n *Nats) Subscribe(s *messaging.Subscription, h messaging.Handler) error {
	if s == nil {
		return messaging.ErrNilSubscription
	}
	if h == nil {
		return messaging.ErrNilHandler
	}
	if len(strings.TrimSpace(s.Topic)) == 0 {
		return messaging.ErrEmptyTopic
	}

	callback := func(m *nats.Msg) {
		msg := &messaging.SubMessage{}
		if err := jsoniter.Unmarshal(m.Data, msg); err != nil {
			logs.Warnw("Nats: failed to unmarshal data form topic", "topic", s.Topic, "error", err)
			return
		}

		msg.SetAck(func() {
			if err := m.Ack(); err != nil {
				logs.Warnw("Nats: failed to ack", "topic", s.Topic, "error", err)
				return
			}
		})
		msg.SetNak(func() {
			if err := m.Nak(); err != nil {
				logs.Warnw("Nats: failed to nak", "topic", s.Topic, "error", err)
				return
			}
		})
		msg.SetTerm(func() {
			if err := m.Term(); err != nil {
				logs.Warnw("Nats: failed to term", "topic", s.Topic, "error", err)
				return
			}
		})
		msg.SetInProgress(func() {
			if err := m.InProgress(); err != nil {
				logs.Warnw("Nats: failed to inProgress", "topic", s.Topic, "error", err)
				return
			}
		})

		ctx := messaging.WithTopic(context.Background(), s.Topic)
		ctx = messaging.WithMessage(ctx, msg)
		if err := h(ctx, s, msg); err != nil {
			return
		}

		if s.AutoAck {
			msg.Ack()
		}
	}

	if len(s.Group) > 0 {
		sub, err := n.queueSubscribe(s, callback)
		if err != nil {
			logs.Warnw("Nats: failed to subscribe with group", "topic", s.Topic, "group", s.Group, "error", err)
			return err
		}
		n.setSubscription(s, sub)
	} else {
		sub, err := n.subscribe(s, callback)
		if err != nil {
			logs.Warnw("Nats: failed to subscribe", "topic", s.Topic, "error", err)
			return err
		}
		n.setSubscription(s, sub)
	}

	return nil
}

func subscriptionKey(s *messaging.Subscription) string {
	if s == nil {
		return ""
	}
	return strings.Join([]string{s.Name, s.Topic, s.Group}, "|")
}

func (n *Nats) setSubscription(s *messaging.Subscription, sub *nats.Subscription) {
	if n == nil || s == nil || sub == nil {
		return
	}
	key := subscriptionKey(s)
	if len(key) == 0 {
		return
	}

	n.subMu.Lock()
	n.subscriptions[key] = sub
	n.subMu.Unlock()
}

func (n *Nats) queueSubscribe(s *messaging.Subscription, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	if n.js != nil {
		n.updateSubjectName(s.Topic)
		if s.Pull {
			sub, err = n.pullSubscribe(s, cb)
		} else {
			sub, err = n.js.QueueSubscribe(s.Topic, s.Group, cb)
		}
	} else {
		sub, err = n.conn.QueueSubscribe(s.Topic, s.Group, cb)
	}
	if err != nil {
		return nil, err
	}

	err = setPendingLimit(sub, s)
	return
}

func (n *Nats) subscribe(s *messaging.Subscription, cb nats.MsgHandler) (sub *nats.Subscription, err error) {
	if n.js != nil {
		n.updateSubjectName(s.Topic)
		sub, err = n.js.Subscribe(s.Topic, cb)
	} else {
		sub, err = n.conn.Subscribe(s.Topic, cb)
	}
	if err != nil {
		return nil, err
	}

	err = setPendingLimit(sub, s)
	return
}

var durableNameRegex = regexp.MustCompile("[a-zA-Z0-9_-]+")

func (n *Nats) pullSubscribe(s *messaging.Subscription, cb nats.MsgHandler) (*nats.Subscription, error) {
	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	durableName := s.Group
	pullNumber := n.pullNumber.Add(1)

	suffix := ""
	segments := strings.Split(s.Topic, ".")
	if len(segments) > 0 {
		suffix = segments[len(segments)-1]
	}
	if !durableNameRegex.MatchString(suffix) {
		suffix = strconv.Itoa(int(pullNumber))
	}

	durableName = strings.Join([]string{durableName, suffix}, "-")

	maxWaiting := s.PullMaxWaiting
	if maxWaiting == 0 {
		maxWaiting = defaultPullMaxWaiting
	}

	subOpts := []nats.SubOpt{nats.PullMaxWaiting(int(maxWaiting))}
	if len(s.AckTimeout) > 0 {
		ackTimeout, err := time.ParseDuration(s.AckTimeout)
		if err != nil {
			logs.Warnw("Nats: failed to parse ack timeout", "ackTimeout", s.AckTimeout, "error", err)
			return nil, err
		}
		subOpts = append(subOpts, nats.AckWait(ackTimeout))
	}

	subs, err := n.js.PullSubscribe(s.Topic, durableName, subOpts...)
	if err != nil {
		logs.Warnw("failed to pull subscribe the topic", "topic", s.Topic, "error", err)
		return nil, err
	}

	go func(s *messaging.Subscription, durableName string) {
		logs.Infow("nats subscribed", "topic", s.Topic, "subscribe", durableName)

		for {
			select {
			case <-n.shutdownCh:
				logs.Infow("shutdown the nats client, will closed the pull subscribe", "subscribe", durableName)
				return
			default:
			}

			logs.Debugw("fetching the next message", "subscribe", durableName)

			ms, err := subs.Fetch(1, nats.MaxWait(time.Second))
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, nats.ErrBadSubscription) {
					logs.Infow("pull subscribe is stopping", "subscribe", durableName, "error", err)
					return
				}
				logs.Warnw("failed to fetch the next message", "subscribe", durableName, "error", err)
				continue
			}

			for _, msg := range ms {
				logs.Debugw("fetched the message", "subscribe", durableName, "subject", msg.Subject)
				cb(msg)
			}
		}
	}(s, durableName)

	return subs, nil
}

func setPendingLimit(sub *nats.Subscription, opts *messaging.Subscription) error {
	if sub == nil || opts == nil {
		return nil
	}

	msgLimit := int(opts.PendingMsgLimit)
	bytesLimit := int(opts.PendingBytesLimit)
	if msgLimit != 0 || bytesLimit != 0 {
		if msgLimit == 0 {
			msgLimit = nats.DefaultSubPendingMsgsLimit
		}
		if bytesLimit == 0 {
			bytesLimit = nats.DefaultSubPendingBytesLimit
		}
		if err := sub.SetPendingLimits(msgLimit, bytesLimit); err != nil {
			logs.Warnw("failed to set pending limits", "name", opts.Name, "topic", opts.Topic,
				"msgLimit", msgLimit, "bytesLimit", bytesLimit, "error", err)
			return err
		}
	}
	return nil
}

// Shutdown shuts down all subscribers
func (n *Nats) Shutdown() {
	if n == nil {
		return
	}

	n.shutdownOnce.Do(func() {
		if n.shutdownCh != nil {
			close(n.shutdownCh)
		}

		n.subMu.Lock()
		for key, sub := range n.subscriptions {
			if sub != nil {
				if err := sub.Unsubscribe(); err != nil {
					logs.Warnw("failed to unsubscribe", "subscription", key, "error", err)
				}
			}
			delete(n.subscriptions, key)
		}
		n.subMu.Unlock()

		if n.conn != nil {
			n.conn.Close()
		}
	})
}

func (n *Nats) GetConn() *nats.Conn {
	if n != nil {
		return n.conn
	}
	return nil
}

func (n *Nats) StreamName() string {
	if n != nil && n.config != nil && n.config.Nats != nil {
		return n.config.Nats.JetStream
	}
	return ""
}

func (n *Nats) SubjectNames() []string {
	if n != nil && n.config != nil && n.config.Nats != nil {
		return n.config.Nats.TopicNames
	}
	return nil
}

func (n *Nats) AppendSubjectName(name string) {
	if n != nil && n.config != nil && n.config.Nats != nil {
		n.config.Nats.TopicNames = append(n.config.Nats.TopicNames, name)
	}
}

func (n *Nats) SetSubjectNames(names ...string) {
	if n != nil && n.config != nil && n.config.Nats != nil {
		n.config.Nats.TopicNames = append([]string{}, names...)
	}
}

func (n *Nats) NewStreamConfig(name string, subjects []string) *nats.StreamConfig {
	streamCfg := &nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	}

	if n.MaxMsgs() > 0 {
		streamCfg.MaxMsgs = n.MaxMsgs()
	}
	if n.MaxAge() > 0 {
		streamCfg.MaxAge = time.Duration(n.MaxAge()) * time.Second
	}

	return streamCfg
}

func (n *Nats) MaxMsgs() int64 {
	if n != nil && n.config != nil && n.config.Nats != nil {
		return n.config.Nats.MaxMsgs
	}
	return 0
}

func (n *Nats) MaxAge() int64 {
	if n != nil && n.config != nil && n.config.Nats != nil {
		return n.config.Nats.MaxAge
	}
	return 0
}

func (n *Nats) GetJetStream() nats.JetStreamContext {
	if n != nil {
		return n.js
	}
	return nil
}
