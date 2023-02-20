package reader

import (
	"github.com/chaos-io/chaos/config/encoder"
	"github.com/chaos-io/chaos/config/encoder/json"
	"github.com/chaos-io/chaos/config/encoder/toml"
	"github.com/chaos-io/chaos/config/encoder/xml"
	"github.com/chaos-io/chaos/config/encoder/yaml"
)

type Options struct {
	Encoding map[string]encoder.Encoder
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]encoder.Encoder{
			"json": json.NewEncoder(),
			"toml": toml.NewEncoder(),
			"xml":  xml.NewEncoder(),
			"yaml": yaml.NewEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		if o.Encoding == nil {
			o.Encoding = make(map[string]encoder.Encoder)
		}
		o.Encoding[e.String()] = e
	}
}
