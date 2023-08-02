package messaging

import "github.com/nats-io/nats.go"

type Config struct {
	Url           string `json:"Url" default:"nats://127.0.0.1:4222"`
	StreamName    string
	Subjects      []string
	subscriptions map[string]*nats.Subscription
	// jetstream.StreamInfo
}
