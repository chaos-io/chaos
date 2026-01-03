package byted

import (
	"time"

	"github.com/coocood/freecache"

	"github.com/chaos-io/chaos/mcache"
)

func NewLRUCache(size int) mcache.MCache {
	return &lruCache{
		c: freecache.NewCache(size),
	}
}

type lruCache struct {
	c *freecache.Cache
}

// Get returns the value or not found error.
func (l *lruCache) Get(key []byte) ([]byte, error) {
	return l.c.Get(key)
}

// Set sets a key, value and expiration for a cache entry and stores it in the cache.
// If the key is larger than 65535 or value is larger than 1/1024 of the cache size,
// the entry will not be written to the cache. expireSeconds <= 0 means no expire,
// but it can be evicted when cache is full.
func (l *lruCache) Set(key, value []byte, expire time.Duration) error {
	return l.c.Set(key, value, int(expire.Seconds()))
}
