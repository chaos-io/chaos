package redis

import (
	"errors"

	goredis "github.com/redis/go-redis/v9"
)

func IsErrNil(err error) bool {
	return errors.Is(err, goredis.Nil)
}
