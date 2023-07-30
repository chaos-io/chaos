package messaging

import "github.com/nats-io/nats.go/jetstream"

type Config struct {
	Url          string `json:"Url" default:"nats://127.0.0.1:4222"`
	StreamConfig jetstream.StreamConfig
}
