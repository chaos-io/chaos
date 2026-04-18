package errorx

import (
	"errors"
	"fmt"
	"strings"
)

type CodedError interface {
	error
	Code() int32
	Extra() map[string]string
	WithExtra(map[string]string)
}

type Error = CodedError

func New(format string, args ...any) error {
	return WithStack(fmt.Errorf(format, args...))
}

func NewByCode(code int32, opts ...Option) error {
	err := &codedError{
		status: newStatusValueByCode(code),
		stack:  stack(),
		cause:  nil,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// WrapByCode returns an error annotating err with a stack trace
// at the point WrapByCode is called, and the status code.
func WrapByCode(err error, code int32, opts ...Option) error {
	if err == nil {
		return nil
	}

	wrappedErr := &codedError{
		status: newStatusValueByCode(code),
		cause:  err,
	}

	for _, opt := range opts {
		opt(wrappedErr)
	}

	var stackTracker StackTracer
	if errors.As(err, &stackTracker) {
		return wrappedErr
	}

	wrappedErr.stack = stack()

	return wrappedErr
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
func GetStatus(err error) Status {
	if err != nil {
		var s Status
		if errors.As(err, &s) {
			return s
		}
	}
	return nil
}

func GetStatusByCode(code int32) Status {
	return newStatusValueByCode(code)
}

func newStatusValueByCode(code int32) *statusValue {
	registeredStatus := GetRegisteredStatus(code)
	return &statusValue{
		code:             registeredStatus.Code,
		message:          registeredStatus.Message,
		affectsStability: registeredStatus.AffectsStability,
	}
}

// FromStatus converts err to Status
// 解析RPC返回的error，若是Status，返回true
func FromStatus(err error) (Status, bool) {
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
		return msg
	}
	return strings.TrimRight(strings.TrimSpace(msg[:index]), "\n")
}
