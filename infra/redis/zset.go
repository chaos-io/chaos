package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Z = goredis.Z

func ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	return GetRedis().ZAdd(ctx, key, members...).Result()
}

// ZRange 返回指定区间内的元素，按照分数（score）从低到高的顺序排列
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 返回指定区间内的元素，按照分数（score）从高到低的顺序排列
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().ZRevRange(ctx, key, start, stop).Result()
}

func ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]goredis.Z, error) {
	return GetRedis().ZRangeWithScores(ctx, key, start, stop).Result()
}

func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]goredis.Z, error) {
	return GetRedis().ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func ZCard(ctx context.Context, key string) (int64, error) {
	return GetRedis().ZCard(ctx, key).Result()
}

func ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return GetRedis().ZRem(ctx, key, members...).Result()
}

// ZRank 获取有序集合中某个成员的排名（rank），从小到大排名，从0开始计数。
func ZRank(ctx context.Context, key string, member string) (int64, error) {
	return GetRedis().ZRank(ctx, key, member).Result()
}

func ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return GetRedis().ZRevRank(ctx, key, member).Result()
}

// ZRemRangeByRank 删除有序集合（sorted set）中指定排名范围内的元素。
//
// start 和 stop 是按排名位置来指定的，排名从 0 开始。
// start 和 stop 可以是负数，表示从集合的末尾开始计数，-1 是最后一个元素，-2 是倒数第二个元素，以此类推。
func ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return GetRedis().ZRemRangeByRank(ctx, key, start, stop).Result()
}

func ZScore(ctx context.Context, key, member string) (float64, error) {
	return GetRedis().ZScore(ctx, key, member).Result()
}

func ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return GetRedis().ZIncrBy(ctx, key, increment, member).Result()
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
func ZInterStore(ctx context.Context, dest string, keys []string, weights []float64, aggregate string) (int64, error) {
	return GetRedis().ZInterStore(ctx, dest, &goredis.ZStore{
		Keys:      keys,
		Weights:   weights,
		Aggregate: aggregate,
	}).Result()
}

func ZUnionStore(ctx context.Context, dest string, keys []string, weights []float64, aggregate string) (int64, error) {
	return GetRedis().ZUnionStore(ctx, dest, &goredis.ZStore{
		Keys:      keys,
		Weights:   weights,
		Aggregate: aggregate,
	}).Result()
}
