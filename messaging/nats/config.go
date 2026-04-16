package nats

import (
	"time"

	"github.com/chaos-io/chaos/messaging"
)

type Config struct {
	URL       string `json:"url" default:"nats://127.0.0.1:4222"`
	JetStream bool   `json:"jetStream"`
}

type Consumer struct {
	Subscription      messaging.Subscription
	Pull              bool
	AckWait           time.Duration
	PullMaxWaiting    int
	PendingMsgLimit   int
	PendingBytesLimit int
}

func (c Consumer) Validate() error {
	return c.Subscription.Validate()
}
