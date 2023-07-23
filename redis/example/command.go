package example

import (
	"context"

	"github.com/chaos-io/chaos/redis"
)

func Set(redis *redis.Redis, key, value interface{}) (interface{}, error) {
	return redis.Do(context.Background(), "SET", key, value).Result()
}

func Get(redis *redis.Redis, key string) (interface{}, error) {
	return redis.Do(context.Background(), "GET", key).Result()
}

func Del(redis *redis.Redis, key interface{}) (interface{}, error) {
	return redis.Do(context.Background(), "DEL", key).Result()
}
