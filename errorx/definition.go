package errorx

import "strings"

const (
	DefaultMessage    = "Service Internal Error"
	DefaultCountInSLA = true
)

// Definition describes one stable business error code.
type Definition struct {
	Code       int32
	Message    string
	CountInSLA bool
}

type DefineOption func(*Definition)

func Define(code int32, message string, opts ...DefineOption) Definition {
	def := Definition{
		Code:       code,
		Message:    message,
		CountInSLA: DefaultCountInSLA,
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

func CountInSLA(countInSLA bool) DefineOption {
	return func(def *Definition) {
		def.CountInSLA = countInSLA
	}
}

func (def Definition) New(opts ...Option) error {
	err := newError(def, nil)
	applyOptions(err, opts...)
	return err
}

func (def Definition) Wrap(cause error, opts ...Option) error {
	if cause == nil {
		return nil
	}
	err := newError(def, cause)
	applyOptions(err, opts...)
	return err
}

func (def Definition) Is(err error) bool {
	return Is(err, def.Code)
}

func (def Definition) normalized() Definition {
	if strings.TrimSpace(def.Message) == "" {
		def.Message = DefaultMessage
	}
	return def
}
