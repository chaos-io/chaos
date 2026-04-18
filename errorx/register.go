package errorx

import (
	"errors"
	"fmt"
	"sync"
)

const (
	DefaultErrorMsg         = "Service Internal Error"
	DefaultAffectsStability = true
)

var ErrRegisterConflict = errors.New("errorx: register conflict")

var (
	registeredStatuses     map[int32]*RegisteredStatus
	registeredStatusesOnce sync.Once
	registeredStatusesMu   sync.RWMutex
)

type RegisteredStatus struct {
	Code             int32
	Message          string
	AffectsStability bool
}

type RegisterOption func(s *RegisteredStatus)

func Register(code int32, msg string, opts ...RegisterOption) error {
	registeredStatusesOnce.Do(func() {
		registeredStatuses = make(map[int32]*RegisteredStatus)
	})

	s := &RegisteredStatus{
		Code:             code,
		Message:          msg,
		AffectsStability: DefaultAffectsStability,
	}

	for _, opt := range opts {
		opt(s)
	}

	registeredStatusesMu.Lock()
	defer registeredStatusesMu.Unlock()

	if current, ok := registeredStatuses[code]; ok {
		if sameRegisteredStatus(current, s) {
			return nil
		}
		return fmt.Errorf(
			"%w: code=%d current_message=%q new_message=%q current_affects_stability=%t new_affects_stability=%t",
			ErrRegisterConflict,
			code,
			current.Message,
			s.Message,
			current.AffectsStability,
			s.AffectsStability,
		)
	}
	registeredStatuses[code] = s
	return nil
}

func MustRegister(code int32, msg string, opts ...RegisterOption) {
	if err := Register(code, msg, opts...); err != nil {
		panic(err)
	}
}

func WithAffectsStability(affectsStability bool) RegisterOption {
	return func(s *RegisteredStatus) {
		s.AffectsStability = affectsStability
	}
}

func GetRegisteredStatus(code int32) *RegisteredStatus {
	registeredStatusesOnce.Do(func() {
		registeredStatuses = make(map[int32]*RegisteredStatus)
	})

	registeredStatusesMu.RLock()
	defer registeredStatusesMu.RUnlock()

	registeredStatus, ok := registeredStatuses[code]
	if ok {
		return cloneRegisteredStatus(registeredStatus)
	}

	return &RegisteredStatus{
		Code:             code,
		Message:          DefaultErrorMsg,
		AffectsStability: DefaultAffectsStability,
	}
}

func sameRegisteredStatus(left, right *RegisteredStatus) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.Code == right.Code &&
		left.Message == right.Message &&
		left.AffectsStability == right.AffectsStability
}

func cloneRegisteredStatus(s *RegisteredStatus) *RegisteredStatus {
	if s == nil {
		return nil
	}
	return &RegisteredStatus{
		Code:             s.Code,
		Message:          s.Message,
		AffectsStability: s.AffectsStability,
	}
}
