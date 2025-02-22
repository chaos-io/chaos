package messaging

import (
	"context"
	"errors"
	"fmt"

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

func NewClientWith(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("there is no config for the messaging client")
	}

	if queues != nil {
		if c, ok := queues[config.Provider]; ok {
			if q, err := c(config); err != nil {
				return nil, err
			} else {
				return &Client{Provider: q, Config: config}, nil
			}
		}
	}

	return nil, fmt.Errorf("has no proper queue provider for the %q", config.Provider)
}

func (c *Client) GetConfig() *Config {
	if c != nil {
		return c.Config
	}
	return nil
}

func (c *Client) Shutdown() {
	if c != nil {
		c.Provider.Shutdown()
	}
}

func (c *Client) Subscribe(subscription *Subscription, handler Handler) {
	if c != nil {
		c.Provider.Subscribe(subscription, handler)
	}
}

func (c *Client) Publish(ctx context.Context, topic string, messages ...*Message) error {
	if c != nil {
		return c.Provider.Publish(ctx, topic, messages...)
	}

	return nil
}
