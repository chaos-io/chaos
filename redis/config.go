package redis

import "github.com/redis/go-redis/v9"

type Config struct {
	*redis.Options
}
