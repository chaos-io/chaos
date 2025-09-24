package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func IsNilError(err error) bool {
	return errors.Is(err, redis.Nil)
}
