package errorx

import (
	"errors"
	"fmt"
	"strings"
)

// Status is the public view of a registered status-bearing error.
// 如果有获取code或其他扩展字段的需求，再考虑对外暴露接口.
type Status interface {
	error
	Code() int32
}

type statusValue struct {
	code    int32
	message string

	// 稳定性标识 可用于SLA稳定的监测
	// true:会影响系统稳定性, 并体现在接口错误率中
	// false:不影响稳定性
	affectsStability bool

	extra map[string]string // 扩展信息
}

func (s *statusValue) Code() int32 {
	return s.code
}

func (s *statusValue) Error() string {
	return fmt.Sprintf("code=%d message=%s", s.code, s.message)
}

func (s *statusValue) Extra() map[string]string {
	return s.extra
}

func (s *statusValue) WithExtra(m map[string]string) {
	if s.extra == nil {
		s.extra = make(map[string]string)
	}

	for k, v := range m {
		s.extra[k] = v
	}
}

type codedError struct {
	status *statusValue

	// at intnal server
	stack string
	cause error // original error
}

func (w *codedError) Unwrap() error {
	return w.cause
}

func (w *codedError) Error() string {
	b := strings.Builder{}
	b.WriteString(w.status.Error())

	if w.cause != nil {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("cause=%s", w.cause.Error()))
	}
	if w.stack != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("stack=%s", w.stack))
	}

	return b.String()
}

func (w *codedError) Code() int32 {
	return w.status.Code()
}

func (w *codedError) Extra() map[string]string {
	return w.status.Extra()
}

func (w *codedError) WithExtra(m map[string]string) {
	w.status.WithExtra(m)
}

func (w *codedError) StackTrace() string {
	return w.stack
}

func (w *codedError) Is(target error) bool {
	var s Status
	if errors.As(target, &s) && w.status.Code() == s.Code() {
		return true
	}
	return false
}

func (w *codedError) As(target any) bool {
	return errors.As(w.status, target)
}
