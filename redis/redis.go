package redis

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient     Cmdable
	redisClientOnce sync.Once
)

func Redis() Cmdable {
	redisClientOnce.Do(func() {
		redisClient = NewClient(NewConfig())
	})
	return redisClient
}

func InitRedis(cfg *Config) {
	redisClientOnce.Do(func() {
		redisClient = NewClient(cfg)
	})
}

func NewClient(cfg *Config) Cmdable {
	if cfg == nil {
		cfg = NewConfig()
	}

	normalizeConfig(cfg)

	if len(cfg.Connections) == 1 {
		return NewProvider(newSingle(cfg))
	}

	return NewProvider(newCluster(cfg))
}

func newSingle(cfg *Config) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:            cfg.Connections[0],
		Password:        cfg.Password,
		DB:              cfg.DB,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxRetries:      cfg.MaxRetries,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		MinRetryBackoff: cfg.MinRetryBackoff,
	})
}

func newCluster(cfg *Config) redis.UniversalClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           cfg.Connections,
		Password:        cfg.Password,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxRetries:      cfg.MaxRetries,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		MinRetryBackoff: cfg.MinRetryBackoff,
	})
}

func normalizeConfig(cfg *Config) {
	if len(cfg.Connections) == 0 {
		cfg.Connections = []string{":6379"}
	}

	if cfg.MinIdleConns == 0 {
		if cfg.MaxIdleConns > 0 {
			cfg.MinIdleConns = cfg.MaxIdleConns
		} else {
			cfg.MinIdleConns = 100
		}
	}

	if cfg.PoolSize == 0 {
		cfg.PoolSize = 300
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Second
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = cfg.ReadTimeout
	}
}
