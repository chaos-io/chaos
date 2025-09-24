package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/chaos-io/chaos/pkg/logs"
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

// Expire 让键在指定的秒后过期
//
// Returns:
// - 1: 键的过期时间成功设置
// - 0: 键不存在，或者无法设置过期时间（例如，键已经有过期时间）
func Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return GetRedis().Expire(ctx, key, duration).Result()
}

// ExpireAt 将给定键的过期时间设置为给定的UNIX时间戳
func ExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	return GetRedis().ExpireAt(ctx, key, time).Result()
}

// PExpire 让键在指定的毫秒后过期
func PExpire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return GetRedis().PExpire(ctx, key, duration).Result()
}

func PExpireAt(ctx context.Context, key string, time time.Time) (bool, error) {
	return GetRedis().PExpireAt(ctx, key, time).Result()
}

// TTL 查看键的剩余生存时间(Time To Live，单位：秒)
//
// Returns:
// - -1: 表示该键没有过期时间
// - -2: 表示该键已经不存在
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return GetRedis().TTL(ctx, key).Result()
}

// PTTL 查看键的剩余生存时间(Time To Live，单位：毫秒)
func PTTL(ctx context.Context, key string) (time.Duration, error) {
	return GetRedis().PTTL(ctx, key).Result()
}

// Persist 移除键的过期时间
func Persist(ctx context.Context, key string) (bool, error) {
	return GetRedis().Persist(ctx, key).Result()
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
	return GetRedis().Sort(ctx, key, &goredis.Sort{
		By:     by,
		Offset: offset,
		Count:  count,
		Get:    get,
		Order:  order,
		Alpha:  alpha,
	}).Result()
}

// Eval 直接执行 Lua 代码，适用于一次性脚本
func Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	return GetRedis().Eval(ctx, script, keys, args...).Result()
}

// ScriptLoad 预加载脚本，并返回 SHA1，配合 EvalSha 使用
func ScriptLoad(ctx context.Context, script string) (string, error) {
	return GetRedis().ScriptLoad(ctx, script).Result()
}

// EvalSha 使用 Lua 脚本的 SHA1 哈希值来执行脚本，避免重复解析脚本
func EvalSha(ctx context.Context, sha1 string, keys []string, args ...any) (any, error) {
	return GetRedis().EvalSha(ctx, sha1, keys, args...).Result()
}

// AcquireLock 使用 SET 命令的 NX 选项和 PX 选项实现原子加锁
func AcquireLock(ctx context.Context, lockKey, requestId string, ttl time.Duration) bool {
	luaScript := `
		if redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then
			return 1
		else
			return 0
		end`

	res, err := Int(Eval(ctx, luaScript, []string{lockKey}, requestId, ttl.Milliseconds()))
	if err != nil {
		logs.Warnw("failed to acquire lock", "lockKey", lockKey, "requestId", requestId, "error", err)
		return false
	}

	return res == 1
}

// ReleaseLock 释放锁（确保只释放自己加的锁）
func ReleaseLock(ctx context.Context, lockKey, requestId string) bool {
	luaScript := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end`

	res, err := Int(Eval(ctx, luaScript, []string{lockKey}, requestId))
	if err != nil {
		logs.Warnw("failed to release lock", "lockKey", lockKey, "requestId", requestId, "error", err)
		return false
	}

	return res == 1
}

// RenewLock 续约锁（防止锁过期导致业务未完成）
func RenewLock(ctx context.Context, lockKey, requestId string, ttl time.Duration) bool {
	luaScript := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end`

	res, err := Int(Eval(ctx, luaScript, []string{lockKey}, requestId, ttl.Milliseconds()))
	if err != nil {
		logs.Warnw("failed to renew lock", "lockKey", lockKey, "requestId", requestId, "error", err)
		return false
	}

	return res == 1
}
