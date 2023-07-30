package example

import (
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/messaging"
)

var client *messaging.Client
var clientOnce sync.Once

func InitClient() *messaging.Client {
	clientOnce.Do(func() {
		cfg := &messaging.Config{}
		if err := config.ScanFrom(cfg, "messaging"); err != nil {
			panic(fmt.Errorf("failed to get the messaging config, error: %v", err))
		}

		if c := messaging.New(cfg); c == nil {
			panic("created messaging client is nil")
		}
	})

	return client
}
