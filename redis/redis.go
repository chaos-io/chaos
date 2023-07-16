package redis

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

func New(cfg *Config) *redis.Client {
	if cfg == nil {
		cfg = &Config{&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		}}
	}

	Rdb = redis.NewClient(cfg.Options)

	return Rdb
}
