package infra_test

import (
	"context"

	"github.com/chaos-io/chaos/db"
	idredis "github.com/chaos-io/chaos/idgen/redis"
	"github.com/chaos-io/chaos/infra"
	"github.com/chaos-io/chaos/messaging"
	"github.com/chaos-io/chaos/redis"
	"github.com/chaos-io/chaos/storage"
)

func ExampleBuild() {
	components, err := infra.Build(infra.Components{}, infra.BuildOptions{
		DB: &db.Config{
			Driver: db.SqliteDriver,
			DSN:    ":memory:",
		},
		Redis: &redis.Config{
			Addresses: []string{"127.0.0.1:6379"},
		},
		IDGen: &idredis.Config{
			ServerIDs: []int64{1},
		},
		Storage: &storage.Config{
			Vendor:     "s3",
			Endpoint:   "http://127.0.0.1:9000",
			Region:     "us-east-1",
			BucketName: "chaos",
			AccessKey:  "access",
			SecretKey:  "secret",
		},
		Messaging: &messaging.Config{
			Driver: messaging.DriverNATS,
			Nats: messaging.NatsConfig{
				URL: "nats://127.0.0.1:4222",
			},
		},
	})
	if err != nil {
		return
	}
	defer func() {
		_ = components.Close(context.Background())
	}()

	_ = components
}
