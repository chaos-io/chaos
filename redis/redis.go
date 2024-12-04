package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/logs"
)

const ImplementorGoRedis = "go-redis"

var (
	plugins     map[string]func(*Config) Redis
	pluginsOnce sync.Once
)

type Redis interface {
	Do(ctx context.Context, cmd string, args ...any) (any, error)
	Close() error
}

func RegisterPlugin(name string, plugin func(config *Config) Redis) {
	pluginsOnce.Do(func() {
		plugins = make(map[string]func(config *Config) Redis)
	})

	plugins[name] = plugin
}

func New(cfg *Config) Redis {
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

	if len(cfg.Implementor) == 0 {
		cfg.Implementor = ImplementorGoRedis
	}

	if p, ok := plugins[cfg.Implementor]; ok {
		return p(cfg)
	}
	panic(fmt.Sprintf("the redis client implementor (%s) not found", cfg.Implementor))
}
