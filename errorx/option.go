package errorx

import (
	"fmt"
	"strings"
)

type Option func(*Error)

func WithExtra(extra map[string]string) Option {
	return func(err *Error) {
		if err == nil || extra == nil {
			return
		}
		if err.extra == nil {
			err.extra = make(map[string]string, len(extra))
		}
		for k, v := range extra {
			err.extra[k] = v
		}
	}
}

func WithMessageParam(key, value string) Option {
	return func(err *Error) {
		if err == nil || key == "" {
			return
		}
		err.message = strings.ReplaceAll(err.message, fmt.Sprintf("{%s}", key), value)
	}
}

func WithExtraMessage(message string) Option {
	return func(err *Error) {
		if err == nil || message == "" {
			return
		}
		err.message = fmt.Sprintf("%s,%s", err.message, message)
	}
}

func WithoutStack() Option {
	return func(err *Error) {
		if err != nil {
			err.stack = ""
		}
	}
}

func applyOptions(err *Error, opts ...Option) {
	for _, opt := range opts {
		opt(err)
	}
}
