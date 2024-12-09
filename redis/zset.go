package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ZAdd(ctx context.Context, key string, score float64, member ...string) (int64, error) {
	return GetRedis().ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Result()
}

func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().ZRange(ctx, key, start, stop).Result()
}

func ZCard(ctx context.Context, key string) (int64, error) {
	return GetRedis().ZCard(ctx, key).Result()
}

func ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return GetRedis().ZRem(ctx, key, members...).Result()
}
