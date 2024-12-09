package redis

import "context"

func RPush(ctx context.Context, key string, value string) (int64, error) {
	return GetRedis().RPush(ctx, key, value).Result()
}

func RPop(ctx context.Context, key string) (string, error) {
	return GetRedis().RPop(ctx, key).Result()
}

func LPush(ctx context.Context, key string, value string) (int64, error) {
	return GetRedis().LPush(ctx, key, value).Result()
}

func LPop(ctx context.Context, key string) (string, error) {
	return GetRedis().LPop(ctx, key).Result()
}

func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().LRange(ctx, key, start, stop).Result()
}

func LLen(ctx context.Context, key string) (int64, error) {
	return GetRedis().LLen(ctx, key).Result()
}
