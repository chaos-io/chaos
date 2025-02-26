package redis

import (
	"context"
	"time"
)

func RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetRedis().RPush(ctx, key, values).Result()
}

func RPop(ctx context.Context, key string) (string, error) {
	return GetRedis().RPop(ctx, key).Result()
}

func LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetRedis().LPush(ctx, key, values).Result()
}

func LPop(ctx context.Context, key string) (string, error) {
	return GetRedis().LPop(ctx, key).Result()
}

func LIndex(ctx context.Context, key string, index int64) (string, error) {
	return GetRedis().LIndex(ctx, key, index).Result()
}

func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().LRange(ctx, key, start, stop).Result()
}

func LLen(ctx context.Context, key string) (int64, error) {
	return GetRedis().LLen(ctx, key).Result()
}

func LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	return GetRedis().LTrim(ctx, key, start, stop).Result()
}

// BRPop 从一个或多个列表中弹出元素。当列表为空时，BRPop会阻塞，直到列表中有元素可弹出或者超时。
//
// Parameters:
// - timeout: 阻塞的最大时间，单位为秒。为0时表示阻塞直到元素可弹出。
//
// Returns:
// - 返回一个包含2个元素的数组
// 1. 列表的名称（key）
// 2. 从该列表弹出的元素
func BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetRedis().BRPop(ctx, timeout, keys...).Result()
}

func BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetRedis().BLPop(ctx, timeout, keys...).Result()
}

// RPopLPush 从src列表中弹出位于最右端的元素，然后将这个元素推入dest列表的最左端，返回该元素
func RPopLPush(ctx context.Context, src, dest string) (string, error) {
	return GetRedis().RPopLPush(ctx, src, dest).Result()
}

func BRPopLPush(ctx context.Context, src, dest string, timeout time.Duration) (string, error) {
	return GetRedis().BRPopLPush(ctx, src, dest, timeout).Result()
}
