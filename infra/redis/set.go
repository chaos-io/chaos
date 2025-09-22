package redis

import "context"

func SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return GetRedis().SAdd(ctx, key, members...).Result()
}

func SMembers(ctx context.Context, key string) ([]string, error) {
	return GetRedis().SMembers(ctx, key).Result()
}

func SCard(ctx context.Context, key string) (int64, error) {
	return GetRedis().SCard(ctx, key).Result()
}

func SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return GetRedis().SIsMember(ctx, key, member).Result()
}

func SPop(ctx context.Context, key string) (string, error) {
	return GetRedis().SPop(ctx, key).Result()
}

func SRandMember(ctx context.Context, key string) (string, error) {
	return GetRedis().SRandMember(ctx, key).Result()
}

func SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return GetRedis().SRem(ctx, key, members...).Result()
}

func SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return GetRedis().SDiff(ctx, keys...).Result()
}

func SDiffStore(ctx context.Context, dst string, keys ...string) (int64, error) {
	return GetRedis().SDiffStore(ctx, dst, keys...).Result()
}

func SInter(ctx context.Context, keys ...string) ([]string, error) {
	return GetRedis().SInter(ctx, keys...).Result()
}

func SInterStore(ctx context.Context, dst string, keys ...string) (int64, error) {
	return GetRedis().SInterStore(ctx, dst, keys...).Result()
}

func SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return GetRedis().SUnion(ctx, keys...).Result()
}

func SUnionStore(ctx context.Context, dst string, keys ...string) (int64, error) {
	return GetRedis().SUnionStore(ctx, dst, keys...).Result()
}
