package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/chaos-io/chaos/messaging"
)

func TestInitNats(t *testing.T) {
	n := InitNats()

	tmpMsg := messaging.Message{
		Data: "aaa",
	}

	err := n.Publish("events.1", tmpMsg)
	if err != nil {
		t.Logf("publish error: %v", err)
		return
	}

	fmt.Println("published 1 messages")

	for _, sub := range n.Config.Subscriptions {
		n.Subscribe(sub, func(ctx context.Context, subscription *messaging.Subscription, m *messaging.SubMessage) error {
			if m.Data == "aaa" {
				m.Ack()
				fmt.Println("m.Ack")
			} else {
				m.Nak()
				fmt.Println("m.Nak")
			}
			return nil
		})
	}
}
