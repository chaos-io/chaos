package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

// Subscribe 订阅一个或多个频道
func Subscribe(ctx context.Context, channels ...string) *goredis.PubSub {
	return GetRedis().Subscribe(ctx, channels...)
}

// Publish 将消息发布到一个频道
func Publish(ctx context.Context, channel string, message any) {
	GetRedis().Publish(ctx, channel, message)
}
