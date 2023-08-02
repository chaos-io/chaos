package example

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func TestInitClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nats := InitNats()

	stream, err := nats.JetStream.CreateStream(ctx, jetstream.StreamConfig{
		Name:     nats.Config.StreamName,
		Subjects: nats.Config.Subjects,
	})
	if err != nil {
		t.Errorf("create stream error: %v", err)
	}

	// Publish some messages
	for i := 0; i < 100; i++ {
		nats.JetStream.Publish(ctx, "ORDERS.new", []byte("hello message "+strconv.Itoa(i)))
		fmt.Printf("Published hello message %d\n", i)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		fmt.Printf("create consumer error: %v", err)
		return
	}

	// Get 10 messages from the consumer
	messageCounter := 0
	msgs, err := consumer.Fetch(10)
	if err != nil {
		fmt.Printf("consumer fetch error: %v", err)
		return
	}
	for msg := range msgs.Messages() {
		msg.Ack()
		fmt.Printf("Received a JetStream message via fetch: %s\n", string(msg.Data()))
		messageCounter++
	}
	if msgs.Error() != nil {
		fmt.Println("Error during Fetch(): ", msgs.Error())
	}

	// Receive messages continuously in a callback
	cons, _ := consumer.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		fmt.Printf("Received a JetStream message via callback: %s\n", string(msg.Data()))
		messageCounter++
	})
	defer cons.Stop()

	// Iterate over messages continuously
	it, _ := consumer.Messages()
	for i := 0; i < 10; i++ {
		msg, _ := it.Next()
		msg.Ack()
		fmt.Printf("Received a JetStream message via iterator: %s\n", string(msg.Data()))
		messageCounter++
	}
	it.Stop()

	// block until all 100 published messages have been processed
	for messageCounter < 100 {
		time.Sleep(10 * time.Millisecond)
	}
}
