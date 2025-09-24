//go:build local
// +build local

package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	streamName   = "schedule:task:stream"
	groupName    = "schedule:task:group"
	consumerName = "schedule:task:consumer"
)

var ctx = context.Background()

func Test_XAdd(t *testing.T) {
	for i := 0; i < 3; i++ {
		res, err := XAdd(ctx, streamName, "number", i)
		if err != nil {
			fmt.Printf("i=%d, error=%v\n", i, err)
			return
		}

		fmt.Printf("xadd(%d), %v\n", i, res)
	}
}

func Test_XRead(t *testing.T) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	streams, err := XRead(ctx, streamName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	for _, stream := range streams {
		fmt.Printf("stream(%s): %v\n", stream.Stream, stream.Messages)
	}
}

func Test_XGroupCreate(t *testing.T) {
	create, err := XGroupCreate(ctx, streamName, groupName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	fmt.Printf("create: %v\n", create)
}

func Test_XReadGroup(t *testing.T) {
	xStreams, err := XReadGroup(ctx, streamName, groupName, consumerName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	for _, stream := range xStreams {
		fmt.Printf("stream(%s): %v\n", stream.Stream, stream.Messages)
	}
}

func Test_XPending(t *testing.T) {
	pending, err := XPending(ctx, streamName, groupName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	fmt.Printf("pending: %v\n", pending)
}
