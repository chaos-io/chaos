package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Provider interface {
	Raw() goredis.UniversalClient
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

var _ Provider = (*Service)(nil)
