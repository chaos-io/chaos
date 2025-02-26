package redis

import (
	"github.com/redis/go-redis/v9"
)

func Pipeline() redis.Pipeliner {
	return GetRedis().Pipeline()
}
