package redis

import (
	"context"
	"testing"
	"time"

	"github.com/chaos-io/chaos/logs"
)

func Test_Trans(t *testing.T) {
	key := "transKey"

	for i := 0; i < 3; i++ {
		go trans(ctx, key)
	}

	time.Sleep(500 * time.Millisecond)
	_ = Del(ctx, key)
}

// trans 等trans结束时，val才被赋值
func trans(ctx context.Context, key string) {
	pipeline := Pipeline()
	val, _ := pipeline.Incr(ctx, key).Result()
	logs.Debugw("trans incr value", "key", key, "value", val)
	time.Sleep(100 * time.Millisecond)
	val, _ = pipeline.Decr(ctx, key).Result()
	logs.Debugw("trans decr value", "key", key, "value", val)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		logs.Warnw("trans exec err", "error", err)
	}
}

func Test_NotTrans(t *testing.T) {
	key := "notTransKey"

	for i := 0; i < 3; i++ {
		go notTrans(ctx, key)
	}

	time.Sleep(500 * time.Millisecond)
	_ = Del(ctx, key)
}

func notTrans(ctx context.Context, key string) {
	val, _ := Incr(ctx, key)
	logs.Debugw("notTrans incr value", "key", key, "value", val)
	time.Sleep(100 * time.Millisecond)
	val, _ = Decr(ctx, key)
	logs.Debugw("notTrans decr value", "key", key, "value", val)
}
