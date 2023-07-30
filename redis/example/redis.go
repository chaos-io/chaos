package example

import (
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/redis"
)

var rdb *redis.Redis
var rdbOnce sync.Once

func InitRedis() *redis.Redis {
	rdbOnce.Do(func() {
		cfg := &redis.Config{}
		if err := config.ScanFrom(cfg, "redis"); err != nil {
			panic(fmt.Errorf("failed to get the redis config, error: %v", err))
		}

		if rdb = redis.New(cfg); rdb == nil {
			panic("created db is nil")
		}
	})

	return rdb
}
