package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Tx = redis.Tx
type Pipeliner = redis.Pipeliner

func Pipeline() redis.Pipeliner {
	return GetRedis().Pipeline()
}

func Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error {
	return GetRedis().Watch(ctx, fn, keys...)
}
