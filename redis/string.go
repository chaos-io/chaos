package redis

import (
	"context"
	"time"
)

func Set(ctx context.Context, key string, value any, expire time.Duration) (string, error) {
	return GetRedis().Set(ctx, key, value, expire).Result()
}

func Get(ctx context.Context, key string) (string, error) {
	return GetRedis().Get(ctx, key).Result()
}

func SetNX(ctx context.Context, key string, value any, expire time.Duration) (bool, error) {
	return GetRedis().SetNX(ctx, key, value, expire).Result()
}

func Incr(ctx context.Context, key string) (int64, error) {
	return GetRedis().Incr(ctx, key).Result()
}

func IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return GetRedis().IncrBy(ctx, key, increment).Result()
}
