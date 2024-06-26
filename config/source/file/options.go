package file

import (
	"context"
	"io/fs"

	"github.com/chaos-io/chaos/config/source"
)

type (
	filePathKey struct{}
	fsKey       struct{}
)

// WithPath sets the path to file.
func WithPath(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, filePathKey{}, p)
	}
}

// WithFS sets the underlying filesystem to lookup file from  (default os.FS).
func WithFS(fs fs.FS) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, fsKey{}, fs)
	}
}
