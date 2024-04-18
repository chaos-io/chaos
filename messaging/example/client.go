package example

import (
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/messaging"
)

var (
	nats     *messaging.Nats
	natsOnce sync.Once
)

func InitNats() *messaging.Nats {
	natsOnce.Do(func() {
		cfg := &messaging.Config{}
		if err := config.ScanFrom(cfg, "messaging"); err != nil {
			panic(fmt.Errorf("failed to get the messaging config, error: %w", err))
		}

		if nats = messaging.New(cfg); nats == nil {
			panic("created messaging nats is nil")
		}
	})

	return nats
}
