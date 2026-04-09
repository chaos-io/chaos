package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/chaos-io/core/go/chaos/core"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/logs"
)

const (
	VendorMinio = "minio"
	VendorS3    = "s3"
)

var (
	initializers     map[string]initializer
	initializersOnce sync.Once
	initializersMu   sync.RWMutex
)

type initializer func(cfg *Config) Storage

func Register(name string, init initializer) {
	initializersOnce.Do(func() {
		initializers = make(map[string]initializer)
	})

	initializersMu.Lock()
	initializers[name] = init
	initializersMu.Unlock()
}

//go:generate mockgen -destination=mocks/storage.go -package=mocks . Storage
type Storage interface {
	BucketName() string
	SetBucket(ctx context.Context, name string) error

	Read(ctx context.Context, key string, options core.Options) (*Object, error)
	Write(ctx context.Context, object *Object, options core.Options) error

	Download(ctx context.Context, key, path string, options core.Options) error
	Upload(ctx context.Context, localFile, key string, options core.Options) error
	PresignedUrl(ctx context.Context, key string) (string, error)
}

func NewStorage(cfg *Config) Storage {
	if cfg == nil {
		return nil
	}

	initializersMu.RLock()
	init, ok := initializers[cfg.Vendor]
	initializersMu.RUnlock()
	if !ok {
		return nil
	}

	return init(cfg)
}

var (
	storage     Storage
	storageOnce sync.Once
)

func GetStorage() Storage {
	storageOnce.Do(func() {
		conf := &Config{}
		if err := config.ScanFrom(conf, "storage"); err != nil {
			logs.Warnw("failed to get the storage config", "error", err.Error())
			storage = NewDummyStorage()
		} else {
			if storage = NewStorage(conf); storage == nil {
				storage = NewDummyStorage()
			}
		}
	})
	return storage
}

func Read(ctx context.Context, key string, options core.Options) (*Object, error) {
	return GetStorage().Read(ctx, key, options)
}

func Write(ctx context.Context, object *Object, options core.Options) error {
	return GetStorage().Write(ctx, object, options)
}

func Download(ctx context.Context, key string, path string, options core.Options) error {
	return GetStorage().Download(ctx, key, path, options)
}

func Upload(ctx context.Context, localFile string, key string, options core.Options) error {
	return GetStorage().Upload(ctx, localFile, key, options)
}

func PresignedUrl(ctx context.Context, key string) (string, error) {
	return GetStorage().PresignedUrl(ctx, key)
}

type DummyStorage struct {
	err error
}

func NewDummyStorage() Storage                                  { return &DummyStorage{err: errors.New("DummyStorage: not implement")} }
func (s *DummyStorage) BucketName() string                      { return "dummy" }
func (s *DummyStorage) SetBucket(context.Context, string) error { return s.err }
func (s *DummyStorage) Read(context.Context, string, core.Options) (*Object, error) {
	return nil, s.err
}
func (s *DummyStorage) Write(context.Context, *Object, core.Options) error           { return s.err }
func (s *DummyStorage) Download(context.Context, string, string, core.Options) error { return s.err }
func (s *DummyStorage) Upload(context.Context, string, string, core.Options) error   { return s.err }
func (s *DummyStorage) PresignedUrl(context.Context, string) (string, error)         { return "", s.err }
