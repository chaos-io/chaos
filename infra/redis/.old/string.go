package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Set 设置键（key）对应的值（value）
//
// Parameters:
//   - EX seconds：设置键的过期时间（秒）。
//   - PX milliseconds：设置键的过期时间（毫秒）。
//   - NX：只有当键不存在时才设置值（适用于分布式锁）。
//   - XX：只有当键已经存在时才设置值。
//   - KEEPTTL：保留键的现有 TTL（不修改过期时间）。
//
// 示例	SET lock "token" NX PX 5000  # 仅当 lock 不存在时设置，且 5000 毫秒后过期
func Set(ctx context.Context, key string, value any, expire time.Duration) (string, error) {
	return GetRedis().Set(ctx, key, value, expire).Result()
}

func Get(ctx context.Context, key string) (string, error) {
	return GetRedis().Get(ctx, key).Result()
}

// SetNX 设置一个键（key）只有在它不存在的情况下才会被设置（SET if Not Exists）
func SetNX(ctx context.Context, key string, value any, expire time.Duration) (bool, error) {
	return GetRedis().SetNX(ctx, key, value, expire).Result()
}

func Incr(ctx context.Context, key string) (int64, error) {
	return GetRedis().Incr(ctx, key).Result()
}

func IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return GetRedis().IncrBy(ctx, key, increment).Result()
}

func IncrByFloat(ctx context.Context, key string, increment float64) (float64, error) {
	return GetRedis().IncrByFloat(ctx, key, increment).Result()
}

func Decr(ctx context.Context, key string) (int64, error) {
	return GetRedis().Decr(ctx, key).Result()
}

func DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return GetRedis().DecrBy(ctx, key, decrement).Result()
}

// Append 将值value追加到key当前值的末尾
func Append(ctx context.Context, key, value string) (int64, error) {
	return GetRedis().Append(ctx, key, value).Result()
}

// GetRange 获取一个由偏移量start至偏移量end范围内所有字符组成的子串，包含start和end在内
func GetRange(ctx context.Context, key string, start, stop int64) (string, error) {
	return GetRedis().GetRange(ctx, key, start, stop).Result()
}

// SetRange 将从start偏移量为offset的子串设置为给定值
func SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	return GetRedis().SetRange(ctx, key, offset, value).Result()
}

// SetBit 将字节串看作是二进制串，并将位串中偏移量为offset的二进制位的值设置为value
//
// Parameters:
// - offset: 偏移量（位的位置），从左往右数
// - value: 要设置的值，0 或 1
//
// Returns:
// - 返回设置之前的位值（0 或 1）
func SetBit(ctx context.Context, key string, offset int64, value int) (int64, error) {
	return GetRedis().SetBit(ctx, key, offset, value).Result()
}

// GetBit 将字节串看作是二进制串，并返回位串中偏移量为offset的二进制位的值
//
// Parameters:
// - offset: 偏移量（位的位置）
//
// Returns:
// - 返回当前位的位值
func GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	return GetRedis().GetBit(ctx, key, offset).Result()
}

// BitCount 统计二进制位串里面值为1的二进制位的数量，如何给定了可选的start和end偏移量，那么只对该偏移量氛围内的二进制位进行统计
//
// Parameters:
// - start(可选): 偏移量开始的位的位置，包含
// - end(可选): 偏移量结束的位的位置，包含
//
// Returns:
// - 二进制位为1的数量
func BitCount(ctx context.Context, key string, start, end int64) (int64, error) {
	return GetRedis().BitCount(ctx, key, &goredis.BitCount{
		Start: start,
		End:   end,
	}).Result()
}

// TODO: bit operation: and, or, xor, not
