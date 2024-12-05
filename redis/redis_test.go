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

func TestNew(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "rdb",
			args:    args{cfg: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !tt.wantErr {
				t.Errorf("New() = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestXAdd(t *testing.T) {
	for i := 0; i < 3; i++ {
		res, err := XAdd(ctx, streamName, "number", i)
		if err != nil {
			fmt.Printf("i=%d, error=%v\n", i, err)
			return
		}

		fmt.Printf("xadd(%d), %v\n", i, res)
	}
}

func TestXRead(t *testing.T) {
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

func TestXGroupCreate(t *testing.T) {
	create, err := XGroupCreate(ctx, streamName, groupName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	fmt.Printf("create: %v\n", create)
}

func TestXReadGroup(t *testing.T) {
	xStreams, err := XReadGroup(ctx, streamName, groupName, consumerName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	for _, stream := range xStreams {
		fmt.Printf("stream(%s): %v\n", stream.Stream, stream.Messages)
	}
}

func TestXPending(t *testing.T) {
	pending, err := XPending(ctx, streamName, groupName)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}

	fmt.Printf("pending: %v\n", pending)
}
