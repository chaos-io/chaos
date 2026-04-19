package messaging

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Provider interface {
	Publish(ctx context.Context, topic string, messages ...*Message) error
	Subscribe(subscription *Subscription, handler Handler) error
	Shutdown()
}

var (
	initializers     map[string]initializer
	initializersOnce sync.Once
	initializersMu   sync.RWMutex
)

type initializer func(cfg *Config) (Queue, error)

var _ Provider = (*Client)(nil)

func Register(name string, init initializer) {
	if init == nil {
		panic("messaging: initializer is nil")
	}

	initializersOnce.Do(func() {
		initializers = make(map[string]initializer)
	})

	initializersMu.Lock()
	initializers[providerName(name)] = init
	initializersMu.Unlock()
}

func buildQueue(cfg *Config) (Queue, error) {
	normalized, err := cfg.normalized()
	if err != nil {
		return nil, err
	}

	initializersMu.RLock()
	init, ok := initializers[providerName(normalized.Driver)]
	initializersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedDriver, normalized.Driver)
	}

	return init(normalized)
}

func providerName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
