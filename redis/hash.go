package redis

import "context"

func HSet(ctx context.Context, key string, field string, value any) (int64, error) {
	return GetRedis().HSet(ctx, key, field, value).Result()
}

func HGet(ctx context.Context, key string, field string) (string, error) {
	return GetRedis().HGet(ctx, key, field).Result()
}

func HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return GetRedis().HMGet(ctx, key, fields...).Result()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return GetRedis().HGetAll(ctx, key).Result()
}

func HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	return GetRedis().HIncrBy(ctx, key, field, incr).Result()
}
