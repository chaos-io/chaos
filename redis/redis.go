package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
	Config *Config
}

func New(cfg *Config) *Redis {
	if cfg == nil {
		cfg = &Config{&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		}}
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Second
	}

	if cfg.PoolSize == 0 {
		cfg.PoolSize = 300
	}

	if cfg.MinIdleConns == 0 {
		if cfg.MaxIdleConns > 0 {
			cfg.MinIdleConns = cfg.MaxIdleConns
		} else {
			cfg.MinIdleConns = 100
		}
	}

	return &Redis{
		redis.NewClient(cfg.Options),
		cfg,
	}
}
