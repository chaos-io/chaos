package errorx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type CodedError interface {
	error
	Code() int32
	Message() string
	Extra() map[string]string
}

type Error struct {
	def     Definition
	message string
	extra   map[string]string
	stack   string
	cause   error
}

func newError(def Definition, cause error) *Error {
	def = def.normalized()
	return &Error{
		def:     def,
		message: def.Message,
		stack:   captureStack(),
		cause:   cause,
	}
}

func From(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var coded *Error
	if errors.As(err, &coded) {
		return coded, true
	}
	return nil, false
}

func CodeOf(err error) (int32, bool) {
	coded, ok := From(err)
	if !ok {
		return 0, false
	}
	return coded.Code(), true
}

func Is(err error, code int32) bool {
	coded, ok := From(err)
	return ok && coded.Code() == code
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

func (e *Error) Error() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("code=%d message=%s", e.Code(), e.Message()))
	if e.cause != nil {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("cause=%s", e.cause.Error()))
	}
	if e.stack != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("stack=%s", e.stack))
	}
	return b.String()
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Code() int32 {
	return e.def.Code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Extra() map[string]string {
	if e.extra == nil {
		return nil
	}
	cp := make(map[string]string, len(e.extra))
	for k, v := range e.extra {
		cp[k] = v
	}
	return cp
}

func (e *Error) StackTrace() string {
	return e.stack
}

func (e *Error) Is(target error) bool {
	coded, ok := From(target)
	return ok && coded.Code() == e.Code()
}

func captureStack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])

	b := strings.Builder{}
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pcs[i])
		b.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, trimFuncName(fn.Name())))
	}
	return b.String()
}

func trimFuncName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	if i == -1 {
		return name
	}
	return name[i+1:]
}
