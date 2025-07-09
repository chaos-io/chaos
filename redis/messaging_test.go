//go:build local
// +build local

package redis

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chaos-io/chaos/logs"
)

func Test_Publish(t *testing.T) {
	t.Run("Publish and subscribe", func(t *testing.T) {
		go runPublish()
		publisher(6)
		time.Sleep(1 * time.Second)
	})
}

func publisher(n int) {
	time.Sleep(1 * time.Second)

	for n > 0 {
		Publish(ctx, "channel", n)
		n--
	}
}

func runPublish() {
	pubSub := Subscribe(ctx, "channel")
	defer pubSub.Close()

	var count int32
	for item := range pubSub.Channel() {
		fmt.Println(item)
		atomic.AddInt32(&count, 1)
		fmt.Println(count)

		switch count {
		case 4:
			if err := pubSub.Unsubscribe(ctx, "channel"); err != nil {
				logs.Warnw("unsubscribe error", "error", err)
			} else {
				logs.Infow("unsubscribe success")
			}
		case 5:
			break
		default:
		}
	}
}
