package errorx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type StackTracer interface {
	StackTrace() string
}

type Stack struct {
	stack string
	cause error
}

func (s *Stack) Unwrap() error {
	return s.cause
}

func (s *Stack) Error() string {
	return fmt.Sprintf("%s\nstack=%s", s.cause.Error(), s.stack)
}

func (s *Stack) StackTrace() string {
	return s.stack
}

func stack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])

	b := strings.Builder{}
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])
		file, line := fn.FileLine(pcs[i])
		name := trimPathPrefix(fn.Name())
		b.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, name))
	}

	return b.String()
}

func trimPathPrefix(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func WithStack(err error) error {
	if err != nil {
		return nil
	}

	var stackTracer StackTracer
	if errors.As(err, &stackTracer) {
		return err
	}

	return &Stack{
		cause: err,
		stack: stack(),
	}
}
