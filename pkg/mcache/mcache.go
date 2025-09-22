package mcache

import "time"

type IByteCache interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte, expire time.Duration) error
}
