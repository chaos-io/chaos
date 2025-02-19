package redis

import (
	"context"
	"time"
)

func Do(ctx context.Context, args ...any) (any, error) {
	return GetRedis().Do(ctx, args...).Result()
}

func Del(ctx context.Context, keys ...string) error {
	return GetRedis().Del(ctx, keys...).Err()
}

func Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := GetRedis().Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func Type(ctx context.Context, key string) (string, error) {
	return GetRedis().Type(ctx, key).Result()
}

func TTL(ctx context.Context, key string) (time.Duration, error) {
	return GetRedis().TTL(ctx, key).Result()
}

func Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return GetRedis().Expire(ctx, key, duration).Result()
}
