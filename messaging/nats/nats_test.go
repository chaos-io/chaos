package nats

import (
	"errors"
	"testing"

	"github.com/chaos-io/chaos/messaging"
)

func TestConfigNormalized(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var cfg *messaging.NatsConfig
		_, err := normalizeConfig(cfg)
		if !errors.Is(err, messaging.ErrNilConfig) {
			t.Fatalf("expect ErrNilConfig, got %v", err)
		}
	})

	t.Run("empty url", func(t *testing.T) {
		_, err := normalizeConfig(&messaging.NatsConfig{})
		if !errors.Is(err, ErrEmptyURL) {
			t.Fatalf("expect ErrEmptyURL, got %v", err)
		}
	})

	t.Run("trim url", func(t *testing.T) {
		cfg, err := normalizeConfig(&messaging.NatsConfig{URL: " nats://127.0.0.1:4222 "})
		if err != nil {
			t.Fatalf("normalized() failed: %v", err)
		}
		if cfg.URL != "nats://127.0.0.1:4222" {
			t.Fatalf("expect trimmed url, got %q", cfg.URL)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("new with config validate args", func(t *testing.T) {
		client, err := NewWithConfig(nil)
		if !errors.Is(err, messaging.ErrNilConfig) {
			t.Fatalf("expect ErrNilConfig, got %v", err)
		}
		if client != nil {
			t.Fatalf("expect nil client, got %#v", client)
		}

		client, err = NewWithConfig(&messaging.NatsConfig{})
		if !errors.Is(err, ErrEmptyURL) {
			t.Fatalf("expect ErrEmptyURL, got %v", err)
		}
		if client != nil {
			t.Fatalf("expect nil client, got %#v", client)
		}
	})
}
