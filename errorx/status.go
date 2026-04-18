package errorx

import (
	"errors"
	"fmt"
	"strings"
)

// Status is interface for error with statusError.
// 如果有获取code或其他扩展字段的需求，再考虑对外暴露接口.
type Status interface {
	error
	Code() int32
}

type status struct {
	code    int32
	message string

	// 稳定性标识 可用于SLA稳定的监测
	// true:会影响系统稳定性, 并体现在接口错误率中
	// false:不影响稳定性
	AffectsStability bool

	extra map[string]string // 扩展信息
}

func (s *status) Code() int32 {
	return s.code
}

func (s *status) Error() string {
	return fmt.Sprintf("code=%d message=%s", s.code, s.message)
}

func (s *status) Extra() map[string]string {
	return s.extra
}

func (s *status) WithExtra(m map[string]string) {
	if s.extra == nil {
		s.extra = make(map[string]string)
	}

	for k, v := range m {
		s.extra[k] = v
	}
}

type withStatus struct {
	status *status

	// at intnal server
	stack string
	cause error // original error
}

func (w *withStatus) Unwrap() error {
	return w.cause
}

func (w *withStatus) Error() string {
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

func (w *withStatus) StackTrace() string {
	return w.stack
}

func (w *withStatus) Is(target error) bool {
	var s Status
	if errors.As(target, &s) && w.status.Code() == s.Code() {
		return true
	}
	return false
}

func (w *withStatus) As(target any) bool {
	return errors.As(w.status, target)
}
