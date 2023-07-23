package example

import (
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/db"
)

var d *db.DB
var dOnce sync.Once

func InitDB() *db.DB {
	dOnce.Do(func() {
		cfg := &db.Config{}
		if err := config.Get("db").Scan(cfg); err != nil {
			panic(fmt.Errorf("failed to get the db config, error: %v", err))
		}

		if d = db.New(cfg); d == nil {
			panic("create db error, db is nil")
		}
	})

	return d
}
