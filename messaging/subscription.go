package messaging

import (
	"strings"
	"time"
)

type Subscription struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
	Group string `json:"group"`

	Pull              bool          `json:"pull"`
	AckWait           time.Duration `json:"ackWait"`
	PullMaxWaiting    int           `json:"pullMaxWaiting"`
	PendingMsgLimit   int           `json:"pendingMsgLimit"`
	PendingBytesLimit int           `json:"pendingBytesLimit"`

	Endpoint Endpoint `json:"endpoint"`
}

type Endpoint struct {
	Service string `json:"service"`
	Method  string `json:"method"`
}

func (s *Subscription) Validate() error {
	if s == nil {
		return ErrNilSubscription
	}
	if strings.TrimSpace(s.Topic) == "" {
		return ErrEmptyTopic
	}
	return nil
}
