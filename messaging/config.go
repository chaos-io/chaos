package messaging

type Config struct {
	Provider      string
	ServiceName   string `json:"ServiceName" default:"nats://127.0.0.1:4222"`
	Subscriptions []*Subscription
	Nats          *Nats
}

type Nats struct {
	JetStream  string
	TopicNames []string
	MaxMsgs    int64
	MaxAge     int64
}
