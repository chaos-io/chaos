package redis

import (
	"sync"
	"time"

	redis2 "github.com/redis/go-redis/v9"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/logs"
)

var (
	redisClient     redis2.UniversalClient
	redisClientOnce sync.Once
)

type Redis struct {
	Client redis2.UniversalClient
}

func GetRedis() redis2.UniversalClient {
	redisClientOnce.Do(func() {
		cfg := NewConfig()
		redisClient = New(cfg).Client
	})
	return redisClient
}

func New(cfg *Config) *Redis {
	if cfg == nil {
		cfg = &Config{}
		if err := config.ScanFrom(cfg, "redis"); err != nil {
			logs.Warnw("not set the config and can't read from the config file, will try to use the default config")
			cfg.Connections = []string{":6379"}
		}
	}

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

	if len(cfg.Connections) == 1 {
		option := &redis2.Options{
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
		}
		return &Redis{Client: redis2.NewClient(option)}
	}

	option := &redis2.ClusterOptions{
		Addrs:           cfg.Connections,
		Password:        cfg.Password,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxRetries:      cfg.MaxRetries,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		MinRetryBackoff: cfg.MinRetryBackoff,
	}

	return &Redis{Client: redis2.NewClusterClient(option)}
}
