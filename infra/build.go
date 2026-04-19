package infra

import (
	"context"
	"errors"
	"sync"

	"github.com/chaos-io/chaos/db"
	"github.com/chaos-io/chaos/idgen"
	idredis "github.com/chaos-io/chaos/idgen/redis"
	"github.com/chaos-io/chaos/messaging"
	"github.com/chaos-io/chaos/messaging/nats"
	chaosredis "github.com/chaos-io/chaos/redis"
	"github.com/chaos-io/chaos/storage"
	"github.com/chaos-io/chaos/storage/minio"
	"github.com/chaos-io/chaos/storage/s3"
)

type Components struct {
	OSS       storage.Storage
	IDGen     idgen.IDGenerator
	DB        db.Provider
	Redis     chaosredis.Provider
	Messaging messaging.Provider

	owned ownedComponents
}

type BuildOptions struct {
	DB        *db.Config
	Redis     *chaosredis.Config
	IDGen     *idredis.Config
	Storage   *storage.Config
	Messaging *messaging.Config

	Registrars []func()
}

type ownedComponents struct {
	OSS       storage.Storage
	IDGen     idgen.IDGenerator
	DB        db.Provider
	Redis     chaosredis.Provider
	Messaging messaging.Provider
}

var defaultRegistrarsOnce sync.Once

func Build(c Components, opts BuildOptions) (Components, error) {
	registerDrivers(opts)

	var owned ownedComponents
	fail := func(err error) (Components, error) {
		_ = closeOwned(context.Background(), owned)
		return c, err
	}

	if c.DB == nil {
		component, err := buildDB(opts.DB)
		if err != nil {
			return fail(err)
		}
		c.DB = component
		owned.DB = component
	}

	if c.Redis == nil {
		component, err := buildRedis(opts.Redis)
		if err != nil {
			return fail(err)
		}
		c.Redis = component
		owned.Redis = component
	}

	if c.IDGen == nil {
		component, err := buildIDGen(c.Redis, opts.IDGen)
		if err != nil {
			return fail(err)
		}
		c.IDGen = component
		owned.IDGen = component
	}

	if c.OSS == nil {
		component, err := buildStorage(opts.Storage)
		if err != nil {
			return fail(err)
		}
		c.OSS = component
		owned.OSS = component
	}

	if c.Messaging == nil {
		component, err := buildMessaging(opts.Messaging)
		if err != nil {
			return fail(err)
		}
		c.Messaging = component
		owned.Messaging = component
	}

	c.owned = owned
	return c, nil
}

func (c Components) Close(ctx context.Context) error {
	return closeOwned(ctx, c.owned)
}

func closeOwned(ctx context.Context, owned ownedComponents) error {
	var errs []error

	if owned.Messaging != nil {
		owned.Messaging.Shutdown()
	}

	if closer, ok := owned.IDGen.(interface{ Close() error }); ok {
		errs = append(errs, closer.Close())
	}

	if owned.Redis != nil {
		errs = append(errs, owned.Redis.Close(ctx))
	}

	if owned.DB != nil {
		errs = append(errs, owned.DB.Close(ctx))
	}

	return errors.Join(errs...)
}

func registerDrivers(opts BuildOptions) {
	defaultRegistrarsOnce.Do(func() {
		minio.Register()
		s3.Register()
		nats.Register()
	})

	runRegistrars(opts.Registrars)
}

func runRegistrars(registers []func()) {
	for _, register := range registers {
		if register != nil {
			register()
		}
	}
}

func buildDB(cfg *db.Config) (db.Provider, error) {
	if cfg != nil {
		return db.NewWithConfig(cfg)
	}
	return db.New()
}

func buildRedis(cfg *chaosredis.Config) (chaosredis.Provider, error) {
	if cfg != nil {
		return chaosredis.NewWithConfig(cfg)
	}
	return chaosredis.New()
}

func buildIDGen(redisProvider chaosredis.Provider, cfg *idredis.Config) (idgen.IDGenerator, error) {
	if redisProvider == nil {
		if cfg == nil {
			return idredis.New()
		}
		return idredis.NewWithConfig(cfg)
	}
	return idredis.NewWithProvider(cfg, redisProvider)
}

func buildStorage(cfg *storage.Config) (storage.Storage, error) {
	if cfg != nil {
		return storage.NewWithConfig(cfg)
	}
	return storage.New()
}

func buildMessaging(cfg *messaging.Config) (messaging.Provider, error) {
	if cfg != nil {
		return messaging.NewWithConfig(cfg)
	}
	return messaging.New()
}
