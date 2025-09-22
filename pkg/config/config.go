// Package config is an interface for dynamic configuration.
package config

import (
	"context"
	"strings"

	"github.com/chaos-io/chaos/pkg/config/loader"
	"github.com/chaos-io/chaos/pkg/config/reader"
	"github.com/chaos-io/chaos/pkg/config/source"
	"github.com/chaos-io/chaos/pkg/config/source/file"
	"github.com/chaos-io/chaos/pkg/config/source/flag"
)

// Config is an interface abstraction for dynamic configuration.
type Config interface {
	// provide the reader.Values interface
	reader.Values
	// Init the config
	Init(opts ...Option) error
	// Options in the config
	Options() Options
	// Close Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Sync Force a source changeset sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}

// Watcher is the config watcher.
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	// for alternative data
	Context context.Context

	Source []source.Source

	WithWatcherDisabled bool
}

type Option func(o *Options)

// Default Config Manager.
var DefaultConfig, _ = NewConfig()

// NewConfig returns new config.
func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

// Bytes Return config as raw json.
func Bytes() []byte {
	return DefaultConfig.Bytes()
}

// Map Return config as a map.
func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

// Scan values to a go type.
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

// ScanFrom scan from the specifier keys to a go type
func ScanFrom(v interface{}, key string, alternatives ...string) error {
	val, err := Get(key)
	if err != nil {
		return err
	}

	for _, alter := range alternatives {
		if !val.Null() {
			break
		}

		val, err = Get(alter)
		if err != nil {
			return err
		}
	}

	return val.Scan(v)
}

// Sync Force a source changeset sync.
func Sync() error {
	return DefaultConfig.Sync()
}

// Get a value from the config.
func Get(path ...string) (reader.Value, error) {
	return DefaultConfig.Get(normalizePath(path...)...)
}

// Load config sources.
func Load(source ...source.Source) error {
	return DefaultConfig.Load(source...)
}

// Watch a value for changes.
func Watch(path ...string) (Watcher, error) {
	return DefaultConfig.Watch(normalizePath(path...)...)
}

// LoadFile is short hand for creating a file source and loading it.
func LoadFile(path string) error {
	return Load(file.NewSource(
		file.WithPath(path),
	))
}

func LoadPath(path string) error {
	return Load(newFileSources(path, "")...)
}

func LoadPathWithSuffix(path string, suffix string) error {
	return Load(newFileSources(path, suffix)...)
}

// LoadFlag load command-line parameters
func LoadFlag() error {
	return Load(flag.NewSource())
}

func normalizePath(path ...string) []string {
	var segments []string
	for _, p := range path {
		s := strings.Split(p, ".")
		segments = append(segments, s...)
	}
	return segments
}
