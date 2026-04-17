package errorx

import (
	"errors"
	"fmt"
	"strings"
)

type Error interface {
	error
	Code() int32
	Extra() map[string]string
	WithExtra(map[string]string)
}

func New(format string, args ...any) error {
	return WithStack(fmt.Errorf(format, args...))
}

func NewByCode(code int32, opts ...Option) *withStatus {
	ws := &withStatus{
		status: GetStatusByCode(code),
		stack:  stack(),
		cause:  nil,
	}

	for _, opt := range opts {
		opt(ws)
	}

	return ws
}

// WrapByCode returns an error annotating err with a stack trace
// at the point WrapByCode is called, and the status code.
func WrapByCode(err error, code int32, opts ...Option) *withStatus {
	if err == nil {
		return nil
	}

	ws := &withStatus{
		status: GetStatusByCode(code),
		cause:  err,
	}

	for _, opt := range opts {
		opt(ws)
	}

	var stackTracker StackTracer
	if errors.As(err, &stackTracker) {
		return ws
	}

	ws.stack = stack()

	return ws
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return WrapfWithStack(err, format, args...)
}

// GetStatus 获取错误链中最顶层的 Status
// 如果有获取code或其他扩展字段的需求，再考虑对外暴露
func GetStatus(err error) *status {
	if err != nil {
		var s *status
		if errors.As(err, &s) {
			return s
		}
	}
	return nil
}

func GetStatusByCode(code int32) *status {
	if rs, ok := initializers[code]; ok {
		return &status{
			code:             rs.Code,
			message:          rs.Message,
			AffectsStability: rs.AffectsStability,
		}
	}
	return &status{
		code:             code,
		message:          DefaultErrorMsg,
		AffectsStability: DefaultAffectsStability,
	}
}

// FromStatus converts err to Status
// 解析RPC返回的error，若是Status，返回true
func FromStatus(err error) (*status, bool) {
	if err == nil {
		return nil, false
	}

	if sErr := GetStatus(err); sErr != nil {
		return sErr, true
	}

	return nil, false
}

func ErrorWithoutStack(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	index := strings.Index(msg, "stack=")
	if index == -1 {
		return msg[:index]
	}
	return msg
}
