package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
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

// Sort 对列表、集合、或者有序集合（sorted set）进行排序
//
// Parameters:
// - key: 需要排序的数据的键
// - by: 按key中的这个元素来排序
// - get: 用于用于获取排序后元素的值
// - order: 排序的顺序，ASC升序（默认），DESC降序
// - offset: 对排序结果进行分页，offset为起始位置
// - count: count为返回的元素数量
// - alpha: 为true时表示按字母排序，默认为按数字排序
func Sort(ctx context.Context, key string, by string, get []string, order string, offset, count int64, alpha bool) ([]string, error) {
	return GetRedis().Sort(ctx, key, &redis.Sort{
		By:     by,
		Offset: offset,
		Count:  count,
		Get:    get,
		Order:  order,
		Alpha:  alpha,
	}).Result()
}
