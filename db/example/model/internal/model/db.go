package model

import (
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/db"
)

var d *db.DB
var dOnce sync.Once

func initDB() *db.DB {
	dOnce.Do(func() {
		cfg := &db.Config{}
		if err := config.ScanFrom(cfg, "db"); err != nil {
			panic(fmt.Errorf("failed to get the db config, error: %v", err))
		}

		if d = db.New(cfg); d == nil {
			panic("created db is nil")
		}
	})

	return d
}
