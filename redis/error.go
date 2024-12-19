package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func IsErrNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
