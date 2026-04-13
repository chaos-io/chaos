package storage

import "time"

type RequestOptions struct {
	Bucket      string
	Concurrency int
	TTL         time.Duration
}

type Option func(*RequestOptions)

func WithConcurrency(n int) Option {
	return func(o *RequestOptions) {
		if n > 0 {
			o.Concurrency = n
		}
	}
}

func WithSignTTL(ttl time.Duration) Option {
	return func(o *RequestOptions) {
		if ttl > 0 {
			o.TTL = ttl
		}
	}
}

func WithSignBucket(bucket string) Option {
	return func(o *RequestOptions) {
		if bucket != "" {
			o.Bucket = bucket
		}
	}
}

func ApplyOptions(opts ...Option) RequestOptions {
	o := RequestOptions{
		Concurrency: DefaultConcurrency,
		TTL:         DefaultSignTTL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}
