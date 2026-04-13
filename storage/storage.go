package storage

import (
	"context"
	"errors"
	"sync"

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
	Read(ctx context.Context, key string, opts ...Option) (*Object, error)
	Write(ctx context.Context, object *Object, opts ...Option) error

	Download(ctx context.Context, key, path string, opts ...Option) error
	Upload(ctx context.Context, localFile, key string, opts ...Option) error

	PresignedDownloadURL(ctx context.Context, key string, opts ...Option) (string, error)
	PresignedUploadURL(ctx context.Context, key string, opts ...Option) (string, error)
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

func Read(ctx context.Context, key string, opts ...Option) (*Object, error) {
	return GetStorage().Read(ctx, key, opts...)
}

func Write(ctx context.Context, object *Object, opts ...Option) error {
	return GetStorage().Write(ctx, object, opts...)
}

func Download(ctx context.Context, key string, path string, opts ...Option) error {
	return GetStorage().Download(ctx, key, path, opts...)
}

func Upload(ctx context.Context, localFile string, key string, opts ...Option) error {
	return GetStorage().Upload(ctx, localFile, key, opts...)
}

func PresignedDownloadURL(ctx context.Context, key string, opts ...Option) (string, error) {
	return GetStorage().PresignedDownloadURL(ctx, key, opts...)
}

func PresignedUploadURL(ctx context.Context, key string, opts ...Option) (string, error) {
	return GetStorage().PresignedUploadURL(ctx, key, opts...)
}

type DummyStorage struct {
	err error
}

func NewDummyStorage() Storage                                                    { return &DummyStorage{err: errors.New("DummyStorage: not implement")} }
func (s *DummyStorage) Read(context.Context, string, ...Option) (*Object, error)  { return nil, s.err }
func (s *DummyStorage) Write(context.Context, *Object, ...Option) error           { return s.err }
func (s *DummyStorage) Download(context.Context, string, string, ...Option) error { return s.err }
func (s *DummyStorage) Upload(context.Context, string, string, ...Option) error   { return s.err }
func (s *DummyStorage) PresignedDownloadURL(context.Context, string, ...Option) (string, error) {
	return "", s.err
}
func (s *DummyStorage) PresignedUploadURL(context.Context, string, ...Option) (string, error) {
	return "", s.err
}
