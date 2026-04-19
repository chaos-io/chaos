package messaging

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	chaosconfig "github.com/chaos-io/chaos/config"
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

func TestNew(t *testing.T) {
	driver := "stub-client"
	defaultDriver := "stub-default"
	Register(driver, func(cfg *Config) (Queue, error) {
		_ = cfg
		return &stubQueue{}, nil
	})
	Register(defaultDriver, func(cfg *Config) (Queue, error) {
		_ = cfg
		return &stubQueue{}, nil
	})

	t.Run("new client", func(t *testing.T) {
		t.Run("nil queue", func(t *testing.T) {
			_, err := NewClient(nil)
			if !errors.Is(err, ErrNilQueue) {
				t.Fatalf("expect ErrNilQueue, got %v", err)
			}
		})

		t.Run("init", func(t *testing.T) {
			client, err := NewClient(&stubQueue{})
			if err != nil {
				t.Fatalf("new client failed: %v", err)
			}
			if client.queue == nil {
				t.Fatal("expect non-nil queue")
			}
		})
	})

	t.Run("new with config", func(t *testing.T) {
		client, err := NewWithConfig(nil)
		if !errors.Is(err, ErrNilConfig) {
			t.Fatalf("expect ErrNilConfig, got %v", err)
		}
		if client != nil {
			t.Fatalf("expect nil client, got %#v", client)
		}
	})

	t.Run("unsupported driver", func(t *testing.T) {
		client, err := NewWithConfig(&Config{Driver: "missing"})
		if !errors.Is(err, ErrUnsupportedDriver) {
			t.Fatalf("expect ErrUnsupportedDriver, got %v", err)
		}
		if client != nil {
			t.Fatalf("expect nil client, got %#v", client)
		}
	})

	t.Run("new with config init", func(t *testing.T) {
		client, err := NewWithConfig(&Config{Driver: driver})
		if err != nil {
			t.Fatalf("new client with config failed: %v", err)
		}
		if client == nil || client.queue == nil {
			t.Fatal("expect non-nil client queue")
		}
	})

	t.Run("loads config", func(t *testing.T) {
		loadTestConfig(t, "messaging.yaml", "messaging:\n  driver: "+defaultDriver+"\n")

		client, err := New()
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}
		if client == nil || client.queue == nil {
			t.Fatal("expect non-nil client queue")
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
	if err := client.Publish(context.Background(), "topic", &Message{}); !errors.Is(err, ErrNilQueue) {
		t.Fatalf("expect ErrNilQueue, got %v", err)
	}
	if err := client.Subscribe(&Subscription{Topic: "topic"}, func(context.Context, *Subscription, *SubMessage) error { return nil }); !errors.Is(err, ErrNilQueue) {
		t.Fatalf("expect ErrNilQueue, got %v", err)
	}

	client.queue = &stubQueue{}
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
	client.queue = &stubQueue{publishErr: publishErr}
	if err := client.Publish(context.Background(), "topic", &Message{}); !errors.Is(err, publishErr) {
		t.Fatalf("expect publish error passthrough, got %v", err)
	}

	subscribeErr := errors.New("subscribe err")
	client.queue = &stubQueue{subscribeErr: subscribeErr}
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
	if !msg.Done() {
		t.Fatal("expect message marked done after ack")
	}
}

func loadTestConfig(t *testing.T, filename, body string) {
	t.Helper()

	if err := chaosconfig.InitDefault(chaosconfig.WithWatcherDisabled()); err != nil {
		t.Fatalf("InitDefault() failed: %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if err := chaosconfig.LoadPath(filepath.Join(dir, "config")); err != nil {
		t.Fatalf("LoadPath() failed: %v", err)
	}
}
