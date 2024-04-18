package storage

import (
	"errors"
	"sync"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/core"
	"github.com/chaos-io/chaos/logs"
)

var (
	initializers     map[string]initializer
	initializersOnce sync.Once
)

type initializer func(cfg *Config) Storage

func Register(name string, init initializer) {
	initializersOnce.Do(func() {
		initializers = make(map[string]initializer)
	})

	initializers[name] = init
}

type Storage interface {
	BucketName() string
	SetBucket(name string) error

	Read(key string, options core.Options) (*Object, error)
	Write(object *Object, options core.Options) error

	Download(key string, path string, options core.Options) error
	Upload(localFile string, key string, options core.Options) error
}

func NewStorage(cfg *Config) Storage {
	if init, ok := initializers[cfg.Vendor]; ok {
		return init(cfg)
	}
	return nil
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

func Read(key string, options core.Options) (*Object, error) {
	return GetStorage().Read(key, options)
}

func Write(object *Object, options core.Options) error {
	return GetStorage().Write(object, options)
}

func Download(key string, path string, options core.Options) error {
	return GetStorage().Download(key, path, options)
}

func Upload(localFile string, key string, options core.Options) error {
	return GetStorage().Upload(localFile, key, options)
}

type DummyStorage struct {
	err error
}

func NewDummyStorage() Storage                                      { return &DummyStorage{err: errors.New("DummyStorage: not implement")} }
func (s *DummyStorage) BucketName() string                          { return "dummy" }
func (s *DummyStorage) SetBucket(string) error                      { return s.err }
func (s *DummyStorage) Read(string, core.Options) (*Object, error)  { return nil, s.err }
func (s *DummyStorage) Write(*Object, core.Options) error           { return s.err }
func (s *DummyStorage) Download(string, string, core.Options) error { return s.err }
func (s *DummyStorage) Upload(string, string, core.Options) error   { return s.err }
