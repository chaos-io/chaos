package mcache

import "time"

type MCache interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte, expire time.Duration) error
}
