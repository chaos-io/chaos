package messaging

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/chaos-io/chaos/logs"
)

type Client struct {
	jetstream.JetStream
	Config *Config
}

func New(cfg *Config) *Client {
	// nats.DefaultURL: "nats://127.0.0.1:4222"
	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		panic(fmt.Errorf("failed to connect nats, error: %v", err))
	}

	// Create a JetStream management interface
	js, err := jetstream.New(nc)
	if err != nil {
		panic(fmt.Errorf("failed to create jetstream, error: %v", err))
	}

	return &Client{
		JetStream: js,
		Config:    cfg,
	}
}

func (c *Client) CreateStream(ctx context.Context) {
	if _, err := c.JetStream.CreateStream(ctx, c.Config.StreamConfig); err != nil {
		logs.Warnw("create stream error", "error", err)
	}
}

func (c *Client) UpdateStream(ctx context.Context, streamConfig jetstream.StreamConfig) {
	if _, err := c.JetStream.UpdateStream(ctx, streamConfig); err != nil {
		logs.Warnw("update stream error", "error", err)
	}
}

func (c *Client) DeleteStream(ctx context.Context, stream string) {
	if err := c.JetStream.DeleteStream(ctx, stream); err != nil {
		logs.Warnw("delete stream error", "stream", stream, "error", err)
	}
}

func (c *Client) ListStream(ctx context.Context, opts []jetstream.StreamListOpt) {
	if err := c.JetStream.ListStreams(ctx, opts...); err != nil {
		logs.Warnw("list stream error", "error", err)
	}
}
