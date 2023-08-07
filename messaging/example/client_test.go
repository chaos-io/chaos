package example

import (
	"context"
	"testing"
	"time"

	"github.com/chaos-io/chaos/messaging"
)

func TestInitNats(t *testing.T) {
	n := InitNats()

	defer n.Shutdown()

	err := n.Publish("EVENTS.test", messaging.Message{Data: "aaa1"})
	if err != nil {
		t.Errorf("publish error: %v", err)
		return
	}

	n.Publish("EVENTS.test", messaging.Message{Data: "aaa"})

	t.Log("published 1 messages")

	for _, sub := range n.Config.Subscriptions {
		n.Subscribe(sub, func(ctx context.Context, subscription *messaging.Subscription, m *messaging.SubMessage) error {
			if m.Data == "aaa" {
				m.Ack()
				t.Log("m.Ack")
			} else {
				m.Nak()
				t.Log("m.Nak")
			}
			return nil
		})
	}

	time.Sleep(time.Millisecond)
}
