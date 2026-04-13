package storage

type RequestOptions struct {
	Concurrency int
}

type Option func(*RequestOptions)

func WithConcurrency(n int) Option {
	return func(o *RequestOptions) {
		if n > 0 {
			o.Concurrency = n
		}
	}
}

func ApplyOptions(opts ...Option) RequestOptions {
	o := RequestOptions{
		Concurrency: DefaultConcurrency,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}
