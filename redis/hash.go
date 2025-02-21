package redis

import (
	"context"
)

func HSet(ctx context.Context, key string, values ...any) (int64, error) {
	return GetRedis().HSet(ctx, key, values).Result()
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

func HIncrByFloat(ctx context.Context, key string, field string, incr float64) (float64, error) {
	return GetRedis().HIncrByFloat(ctx, key, field, incr).Result()
}

func HLen(ctx context.Context, key string) (int64, error) {
	return GetRedis().HLen(ctx, key).Result()
}

func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return GetRedis().HDel(ctx, key, fields...).Result()
}

// HExists 检查给定键是否存在于散列中
func HExists(ctx context.Context, key string, field string) (bool, error) {
	return GetRedis().HExists(ctx, key, field).Result()
}

// HKeys 获取散列包含的所有键
func HKeys(ctx context.Context, key string) ([]string, error) {
	return GetRedis().HKeys(ctx, key).Result()
}

// HVals 获取散列包含的所有值
func HVals(ctx context.Context, key string) ([]string, error) {
	return GetRedis().HVals(ctx, key).Result()
}
