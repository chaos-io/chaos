package goredis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/chaos-io/chaos/redis"
)

func init() {
	redis.RegisterPlugin(redis.ImplementorGoRedis, NewRedis)
}

type Redis struct {
	Client goredis.UniversalClient
}

func NewRedis(cfg *redis.Config) redis.Redis {
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
		option := &goredis.Options{
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
		return &Redis{Client: goredis.NewClient(option)}
	} else {
		option := &goredis.ClusterOptions{
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
		return &Redis{Client: goredis.NewClusterClient(option)}
	}
}

func (r *Redis) Do(ctx context.Context, cmd string, arguments ...any) (any, error) {
	args := make([]any, 0, len(arguments)+1)
	args = append(args, cmd)
	args = append(args, arguments...)
	return r.Client.Do(ctx, args...).Result()
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
