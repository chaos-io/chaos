package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

type stubQueue struct {
	publishErr   error
	subscribeErr error
}

func (s *stubQueue) Publish(ctx context.Context, topic string, messages ...*Message) error {
	_ = ctx
	_ = topic
	_ = messages
	return s.publishErr
}

func (s *stubQueue) Subscribe(subscription *Subscription, handler Handler) error {
	_ = subscription
	_ = handler
	return s.subscribeErr
}

func (s *stubQueue) Shutdown() {}

var providerSeq atomic.Uint64

func uniqueProvider(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, providerSeq.Add(1))
}

func TestRegister(t *testing.T) {
	t.Run("invalid provider name", func(t *testing.T) {
		err := Register("  ", func(*Config) (Queue, error) { return &stubQueue{}, nil })
		if !errors.Is(err, ErrInvalidProviderName) {
			t.Fatalf("expect ErrInvalidProviderName, got %v", err)
		}
	})

	t.Run("nil constructor", func(t *testing.T) {
		err := Register(uniqueProvider("nil-constructor"), nil)
		if !errors.Is(err, ErrNilQueueConstructor) {
			t.Fatalf("expect ErrNilQueueConstructor, got %v", err)
		}
	})

	t.Run("duplicate provider", func(t *testing.T) {
		name := uniqueProvider("dup-provider")
		if err := Register(name, func(*Config) (Queue, error) { return &stubQueue{}, nil }); err != nil {
			t.Fatalf("first register failed: %v", err)
		}

		err := Register(" "+name+" ", func(*Config) (Queue, error) { return &stubQueue{}, nil })
		if !errors.Is(err, ErrProviderAlreadyRegistered) {
			t.Fatalf("expect ErrProviderAlreadyRegistered, got %v", err)
		}
	})
}

func TestNewClientWith(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		_, err := NewClientWith(nil)
		if !errors.Is(err, ErrNilConfig) {
			t.Fatalf("expect ErrNilConfig, got %v", err)
		}
	})

	t.Run("unknown provider", func(t *testing.T) {
		_, err := NewClientWith(&Config{Provider: uniqueProvider("not-found")})
		if !errors.Is(err, ErrProviderNotRegistered) {
			t.Fatalf("expect ErrProviderNotRegistered, got %v", err)
		}
	})

	t.Run("normalize provider and init", func(t *testing.T) {
		name := uniqueProvider("normalized")
		if err := Register(name, func(*Config) (Queue, error) { return &stubQueue{}, nil }); err != nil {
			t.Fatalf("register provider failed: %v", err)
		}

		cfg := &Config{Provider: " " + name + " "}
		client, err := NewClientWith(cfg)
		if err != nil {
			t.Fatalf("new client failed: %v", err)
		}
		if cfg.Provider != name {
			t.Fatalf("expect provider %q, got %q", name, cfg.Provider)
		}
		if client.Provider == nil {
			t.Fatal("expect non-nil provider")
		}
	})

	t.Run("constructor returns nil provider", func(t *testing.T) {
		name := uniqueProvider("nil-provider")
		if err := Register(name, func(*Config) (Queue, error) { return nil, nil }); err != nil {
			t.Fatalf("register provider failed: %v", err)
		}

		_, err := NewClientWith(&Config{Provider: name})
		if !errors.Is(err, ErrNilProvider) {
			t.Fatalf("expect ErrNilProvider, got %v", err)
		}
	})
}

func TestClientValidation(t *testing.T) {
	var nilClient *Client
	if err := nilClient.Publish(context.Background(), "topic", &Message{}); !errors.Is(err, ErrNilClient) {
		t.Fatalf("expect ErrNilClient, got %v", err)
	}
	if err := nilClient.Subscribe(&Subscription{Topic: "topic"}, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, ErrNilClient) {
		t.Fatalf("expect ErrNilClient, got %v", err)
	}

	client := &Client{}
	if err := client.Publish(context.Background(), "topic", &Message{}); !errors.Is(err, ErrNilProvider) {
		t.Fatalf("expect ErrNilProvider, got %v", err)
	}
	if err := client.Subscribe(&Subscription{Topic: "topic"}, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, ErrNilProvider) {
		t.Fatalf("expect ErrNilProvider, got %v", err)
	}

	client.Provider = &stubQueue{}
	if err := client.Publish(context.Background(), " ", &Message{}); !errors.Is(err, ErrEmptyTopic) {
		t.Fatalf("expect ErrEmptyTopic, got %v", err)
	}
	if err := client.Subscribe(nil, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, ErrNilSubscription) {
		t.Fatalf("expect ErrNilSubscription, got %v", err)
	}
	if err := client.Subscribe(&Subscription{Topic: ""}, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, ErrEmptyTopic) {
		t.Fatalf("expect ErrEmptyTopic, got %v", err)
	}
	if err := client.Subscribe(&Subscription{Topic: "topic"}, nil); !errors.Is(err, ErrNilHandler) {
		t.Fatalf("expect ErrNilHandler, got %v", err)
	}

	publishErr := errors.New("publish err")
	client.Provider = &stubQueue{publishErr: publishErr}
	if err := client.Publish(context.Background(), "topic", &Message{}); !errors.Is(err, publishErr) {
		t.Fatalf("expect publish error passthrough, got %v", err)
	}

	subscribeErr := errors.New("subscribe err")
	client.Provider = &stubQueue{subscribeErr: subscribeErr}
	if err := client.Subscribe(&Subscription{Topic: "topic"}, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, subscribeErr) {
		t.Fatalf("expect subscribe error passthrough, got %v", err)
	}
}

func TestContextHelpers(t *testing.T) {
	msg := &SubMessage{
		Message: Message{
			Id:         "id-1",
			Attributes: map[string]any{"k": "v"},
		},
	}

	ctx := WithTopic(nil, "topic-1")
	ctx = WithMessage(ctx, msg)
	if topic := GetContextTopic(ctx); topic != "topic-1" {
		t.Fatalf("expect topic-1, got %q", topic)
	}
	if got := GetContextMessage(ctx); got != msg {
		t.Fatalf("expect same message pointer")
	}
	if id := GetContextMessageId(ctx); id != "id-1" {
		t.Fatalf("expect id-1, got %q", id)
	}
	attrs := GetContextMessageAttributes(ctx)
	if attrs["k"] != "v" {
		t.Fatalf("expect attrs[k]=v, got %#v", attrs["k"])
	}
}

func TestSubMessageAckOnlyOnce(t *testing.T) {
	var ack, nak, term, inProgress atomic.Int64
	msg := &SubMessage{}
	msg.SetAck(func() { ack.Add(1) })
	msg.SetNak(func() { nak.Add(1) })
	msg.SetTerm(func() { term.Add(1) })
	msg.SetInProgress(func() { inProgress.Add(1) })

	msg.Ack()
	msg.Nak()
	msg.Term()
	msg.InProgress()

	if ack.Load() != 1 {
		t.Fatalf("expect ack 1, got %d", ack.Load())
	}
	if nak.Load() != 0 {
		t.Fatalf("expect nak 0, got %d", nak.Load())
	}
	if term.Load() != 0 {
		t.Fatalf("expect term 0, got %d", term.Load())
	}
	if inProgress.Load() != 1 {
		t.Fatalf("expect inProgress 1, got %d", inProgress.Load())
	}
}
