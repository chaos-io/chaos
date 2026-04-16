package messaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/chaos-io/chaos/config"
)

type Client struct {
	Config   *Config
	Provider Queue
}

func NewClient() (*Client, error) {
	cfg := &Config{}
	if err := config.ScanFrom(cfg, "messaging"); err != nil {
		return nil, err
	}
	return NewClientWith(cfg)
}

func NewClientWith(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	cfg.Provider = normalizeProvider(cfg.Provider)
	if cfg.Provider == "" {
		return nil, ErrInvalidProviderName
	}

	constructor, ok := getConstructor(cfg.Provider)
	if !ok {
		return nil, fmt.Errorf("%w: %q (register provider first, e.g. messaging/providers.RegisterDefaults)", ErrProviderNotRegistered, cfg.Provider)
	}

	provider, err := constructor(cfg)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrNilProvider
	}

	return &Client{Provider: provider, Config: cfg}, nil
}

func (c *Client) GetConfig() *Config {
	if c != nil {
		return c.Config
	}
	return nil
}

func (c *Client) Shutdown() {
	if c == nil || c.Provider == nil {
		return
	}
	c.Provider.Shutdown()
}

func (c *Client) Subscribe(subscription *Subscription, handler Handler) error {
	if c == nil {
		return ErrNilClient
	}
	if c.Provider == nil {
		return ErrNilProvider
	}
	if subscription == nil {
		return ErrNilSubscription
	}
	if strings.TrimSpace(subscription.Topic) == "" {
		return ErrEmptyTopic
	}
	if handler == nil {
		return ErrNilHandler
	}
	return c.Provider.Subscribe(subscription, handler)
}

func (c *Client) Publish(ctx context.Context, topic string, messages ...*Message) error {
	if c == nil {
		return ErrNilClient
	}
	if c.Provider == nil {
		return ErrNilProvider
	}
	if strings.TrimSpace(topic) == "" {
		return ErrEmptyTopic
	}
	return c.Provider.Publish(ctx, topic, messages...)
}
