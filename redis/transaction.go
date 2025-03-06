package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Tx = goredis.Tx
type Pipeliner = goredis.Pipeliner

func Pipeline() goredis.Pipeliner {
	return GetRedis().Pipeline()
}

func Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error {
	return GetRedis().Watch(ctx, fn, keys...)
}
