package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ZAdd(ctx context.Context, key string, score float64, member string) (int64, error) {
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

func ZScore(ctx context.Context, key, member string) (float64, error) {
	return GetRedis().ZScore(ctx, key, member).Result()
}

func ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return GetRedis().ZIncrBy(ctx, key, increment, member).Result()
}

// ZRevRange 返回指定区间内的元素，按照分数（score）从高到低的顺序排列
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().ZRevRange(ctx, key, start, stop).Result()
}

// ZInterStore 用于计算多个有序集合（Sorted Sets）交集并将结果存储到一个新有序集合的命令。
//
// destination：存储结果的目标有序集合的键。
// key [key ...]：参与计算的多个有序集合的键。
// WEIGHTS（可选）：给每个有序集合指定一个权重，用来加权计算每个元素的分数。如果没有指定，默认所有集合的权重为 1。
// AGGREGATE（可选）：指定如何计算交集的分数，默认为 sum。可以选择：
//   - sum：对交集元素的分数进行求和（默认）。
//   - min：取交集元素的最小分数。
//   - max：取交集元素的最大分数。
func ZInterStore(ctx context.Context, destination string, keys []string, weights []float64, aggregate string) (int64, error) {
	return GetRedis().ZInterStore(ctx, destination, &redis.ZStore{
		Keys:      keys,
		Weights:   weights,
		Aggregate: aggregate,
	}).Result()
}
