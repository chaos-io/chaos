package messaging

import (
	"fmt"
	"strings"
	"sync"
)

var (
	queues   = make(map[string]Constructor)
	queuesMu sync.RWMutex
)

//go:generate mockgen -destination=mocks/queue.go -package=mocks . Queue
type Queue interface {
	Publisher
	Subscriber

	Shutdown()
}

type Constructor func(config *Config) (Queue, error)

func Register(name string, constructor Constructor) error {
	name = normalizeProvider(name)
	if name == "" {
		return ErrInvalidProviderName
	}
	if constructor == nil {
		return ErrNilQueueConstructor
	}

	queuesMu.Lock()
	defer queuesMu.Unlock()
	if _, exists := queues[name]; exists {
		return fmt.Errorf("%w: %q", ErrProviderAlreadyRegistered, name)
	}
	queues[name] = constructor
	return nil
}

func MustRegister(name string, constructor Constructor) {
	if err := Register(name, constructor); err != nil {
		panic(err)
	}
}

func getConstructor(name string) (Constructor, bool) {
	name = normalizeProvider(name)
	if name == "" {
		return nil, false
	}

	queuesMu.RLock()
	constructor, ok := queues[name]
	queuesMu.RUnlock()
	return constructor, ok
}

func normalizeProvider(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
