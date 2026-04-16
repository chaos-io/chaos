package nats

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"

	"github.com/chaos-io/chaos/logs"
	"github.com/chaos-io/chaos/messaging"
)

const defaultFetchWait = time.Second

type Nats struct {
	conn          *nats.Conn
	js            nats.JetStreamContext
	subscriptions map[string]*nats.Subscription
	shutdownCh    chan struct{}
	shutdownOnce  sync.Once
	subMu         sync.Mutex
}

func New(cfg *Config) (*Nats, error) {
	if cfg == nil {
		return nil, messaging.ErrNilConfig
	}
	if len(strings.TrimSpace(cfg.URL)) == 0 {
		return nil, errors.New("messaging nats: url is empty")
	}

	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		logs.Errorw("failed to connect nats", "error", err)
		return nil, err
	}

	n := &Nats{
		conn:          nc,
		subscriptions: map[string]*nats.Subscription{},
		shutdownCh:    make(chan struct{}),
	}

	if cfg.JetStream {
		js, err := nc.JetStream()
		if err != nil {
			return nil, err
		}
		n.js = js
	}

	return n, nil
}

func (n *Nats) Publish(ctx context.Context, topic string, messages ...*messaging.Message) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

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

	return nil
}

func (n *Nats) publish(ctx context.Context, topic string, msg []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	if n.js != nil {
		if _, err := n.js.Publish(topic, msg); err != nil {
			return err
		}
		return nil
	}

	return n.conn.Publish(topic, msg)
}

func (n *Nats) Subscribe(s *messaging.Subscription, h messaging.Handler) error {
	if err := s.Validate(); err != nil {
		return err
	}
	return n.SubscribeConsumer(Consumer{Subscription: *s}, h)
}

func (n *Nats) SubscribeConsumer(consumer Consumer, h messaging.Handler) error {
	s := &consumer.Subscription
	if h == nil {
		return messaging.ErrNilHandler
	}
	if err := consumer.Validate(); err != nil {
		return err
	}

	callback := func(raw *nats.Msg) {
		msg := &messaging.SubMessage{}
		if err := jsoniter.Unmarshal(raw.Data, msg); err != nil {
			logs.Warnw("Nats: failed to unmarshal data form topic", "topic", s.Topic, "error", err)
			if n.js != nil {
				if termErr := raw.Term(); termErr != nil {
					logs.Warnw("Nats: failed to term invalid message", "topic", s.Topic, "error", termErr)
				}
			}
			return
		}

		msg.SetAck(func() {
			if err := raw.Ack(); err != nil {
				logs.Warnw("Nats: failed to ack", "topic", s.Topic, "error", err)
			}
		})
		msg.SetNak(func() {
			if err := raw.Nak(); err != nil {
				logs.Warnw("Nats: failed to nak", "topic", s.Topic, "error", err)
			}
		})
		msg.SetTerm(func() {
			if err := raw.Term(); err != nil {
				logs.Warnw("Nats: failed to term", "topic", s.Topic, "error", err)
			}
		})
		msg.SetInProgress(func() {
			if err := raw.InProgress(); err != nil {
				logs.Warnw("Nats: failed to inProgress", "topic", s.Topic, "error", err)
			}
		})

		ctx := messaging.WithTopic(context.Background(), s.Topic)
		ctx = messaging.WithMessage(ctx, msg)
		if err := h(ctx, s, msg); err != nil {
			if !msg.Done() {
				msg.Nak()
			}
			return
		}

		if s.AutoAck && !msg.Done() {
			msg.Ack()
		}
	}

	var (
		sub *nats.Subscription
		err error
	)
	if len(s.Group) > 0 {
		sub, err = n.queueSubscribe(consumer, callback)
	} else {
		sub, err = n.subscribe(consumer, callback)
	}
	if err != nil {
		return err
	}

	n.setSubscription(s, sub)
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
	if key == "" {
		return
	}

	n.subMu.Lock()
	n.subscriptions[key] = sub
	n.subMu.Unlock()
}

func (n *Nats) queueSubscribe(consumer Consumer, cb nats.MsgHandler) (*nats.Subscription, error) {
	s := consumer.Subscription
	var (
		sub *nats.Subscription
		err error
	)
	if n.js != nil {
		if consumer.Pull {
			sub, err = n.pullSubscribe(consumer, cb)
		} else {
			sub, err = n.js.QueueSubscribe(s.Topic, s.Group, cb)
		}
	} else {
		sub, err = n.conn.QueueSubscribe(s.Topic, s.Group, cb)
	}
	if err != nil {
		return nil, err
	}

	if err := setPendingLimit(sub, consumer); err != nil {
		return nil, err
	}
	return sub, nil
}

func (n *Nats) subscribe(consumer Consumer, cb nats.MsgHandler) (*nats.Subscription, error) {
	s := consumer.Subscription
	var (
		sub *nats.Subscription
		err error
	)
	if n.js != nil {
		sub, err = n.js.Subscribe(s.Topic, cb)
	} else {
		sub, err = n.conn.Subscribe(s.Topic, cb)
	}
	if err != nil {
		return nil, err
	}

	if err := setPendingLimit(sub, consumer); err != nil {
		return nil, err
	}
	return sub, nil
}

func (n *Nats) pullSubscribe(consumer Consumer, cb nats.MsgHandler) (*nats.Subscription, error) {
	s := consumer.Subscription
	durableName := strings.Join([]string{s.Group, s.Name}, "-")
	if durableName == "-" || durableName == "" {
		durableName = strings.ReplaceAll(s.Topic, ".", "-")
	}

	subOpts := []nats.SubOpt{}
	if consumer.PullMaxWaiting > 0 {
		subOpts = append(subOpts, nats.PullMaxWaiting(consumer.PullMaxWaiting))
	}
	if consumer.AckWait > 0 {
		subOpts = append(subOpts, nats.AckWait(consumer.AckWait))
	}

	sub, err := n.js.PullSubscribe(s.Topic, durableName, subOpts...)
	if err != nil {
		logs.Warnw("failed to pull subscribe the topic", "topic", s.Topic, "error", err)
		return nil, err
	}

	go func() {
		for {
			select {
			case <-n.shutdownCh:
				return
			default:
			}

			ms, err := sub.Fetch(1, nats.MaxWait(defaultFetchWait))
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, nats.ErrBadSubscription) {
					return
				}
				logs.Warnw("failed to fetch the next message", "subscribe", durableName, "error", err)
				continue
			}

			for _, msg := range ms {
				cb(msg)
			}
		}
	}()

	return sub, nil
}

func setPendingLimit(sub *nats.Subscription, consumer Consumer) error {
	if sub == nil {
		return nil
	}

	s := consumer.Subscription
	msgLimit := consumer.PendingMsgLimit
	bytesLimit := consumer.PendingBytesLimit
	if msgLimit != 0 || bytesLimit != 0 {
		if msgLimit == 0 {
			msgLimit = nats.DefaultSubPendingMsgsLimit
		}
		if bytesLimit == 0 {
			bytesLimit = nats.DefaultSubPendingBytesLimit
		}
		if err := sub.SetPendingLimits(msgLimit, bytesLimit); err != nil {
			logs.Warnw("failed to set pending limits", "name", s.Name, "topic", s.Topic,
				"msgLimit", msgLimit, "bytesLimit", bytesLimit, "error", err)
			return err
		}
	}
	return nil
}

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
