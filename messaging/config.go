package messaging

type Config struct {
	Url           string `json:"Url" default:"nats://127.0.0.1:4222"`
	StreamName    string
	TopicNames    []string
	Subscriptions []*Subscription
}
