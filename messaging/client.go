package messaging

import (
	"context"
	"strings"
)

type Client struct {
	queue         Queue
	subscriptions []*Subscription
}

func NewClient(queue Queue) (*Client, error) {
	if queue == nil {
		return nil, ErrNilQueue
	}
	return &Client{queue: queue}, nil
}

func (c *Client) Shutdown() {
	if c == nil || c.queue == nil {
		return
	}
	c.queue.Shutdown()
}

func (c *Client) Subscriptions() []*Subscription {
	if c == nil {
		return nil
	}
	return c.subscriptions
}

func (c *Client) Subscribe(subscription *Subscription, handler Handler) error {
	if c == nil {
		return ErrNilClient
	}
	if c.queue == nil {
		return ErrNilQueue
	}
	if err := subscription.Validate(); err != nil {
		return err
	}
	if handler == nil {
		return ErrNilHandler
	}
	return c.queue.Subscribe(subscription, handler)
}

func (c *Client) Publish(ctx context.Context, topic string, messages ...*Message) error {
	if c == nil {
		return ErrNilClient
	}
	if c.queue == nil {
		return ErrNilQueue
	}
	if strings.TrimSpace(topic) == "" {
		return ErrEmptyTopic
	}
	return c.queue.Publish(ctx, topic, messages...)
}
