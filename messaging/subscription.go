package messaging

import "strings"

type Subscription struct {
	Name    string
	Topic   string
	Group   string
	AutoAck bool
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
