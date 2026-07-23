//go:build local

package nats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chaos-io/chaos/messaging"
)

func TestConfiguredSubscriptionRoundTrip(t *testing.T) {
	Register()

	topic := fmt.Sprintf("fino2.messaging.%d", time.Now().UnixNano())
	config := &messaging.Config{
		Driver: messaging.DriverNATS,
		Nats: messaging.NatsConfig{
			URL: "nats://127.0.0.1:4222",
		},
		Subscriptions: []*messaging.Subscription{{
			Name:  "fino2-test",
			Topic: topic,
			Endpoint: messaging.Endpoint{
				Service: "Agent",
				Method:  "start_task",
			},
		}},
	}

	client, err := messaging.NewWithConfig(config)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	defer client.Shutdown()

	subscriptions := client.Subscriptions()
	if len(subscriptions) != 1 || subscriptions[0].Endpoint.Method != "start_task" {
		t.Fatalf("unexpected subscriptions: %#v", subscriptions)
	}

	received := make(chan *messaging.SubMessage, 1)
	if err := client.Subscribe(subscriptions[0], func(ctx context.Context, _ *messaging.Subscription, message *messaging.SubMessage) error {
		if got := messaging.GetContextTopic(ctx); got != topic {
			t.Errorf("context topic = %q, want %q", got, topic)
		}
		if got := messaging.GetContextMessageId(ctx); got != "message-1" {
			t.Errorf("context message id = %q, want message-1", got)
		}
		received <- message
		return nil
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := client.Publish(context.Background(), topic, &messaging.Message{Id: "message-1", Data: `{"task":"demo"}`}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case message := <-received:
		if message.Data != `{"task":"demo"}` {
			t.Fatalf("message data = %q", message.Data)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for NATS message")
	}
}
