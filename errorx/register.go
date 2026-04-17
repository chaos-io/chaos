package errorx

import "sync"

const (
	DefaultErrorMsg         = "Service Internal Error"
	DefaultAffectsStability = true
)

var (
	initializers     map[int32]*RegisterStatus
	initializersOnce sync.Once
	initializersMu   sync.RWMutex
)

type RegisterStatus struct {
	Code             int32
	Message          string
	AffectsStability bool
}
type RegisterOption func(s *RegisterStatus)

func Register(code int32, msg string, opts ...RegisterOption) {
	initializersOnce.Do(func() {
		initializers = make(map[int32]*RegisterStatus)
	})

	s := &RegisterStatus{
		Code:             code,
		Message:          msg,
		AffectsStability: DefaultAffectsStability,
	}

	for _, opt := range opts {
		opt(s)
	}

	initializersMu.Lock()
	initializers[code] = s
	initializersMu.Unlock()
}

func WithAffectsStability(affectsStability bool) RegisterOption {
	return func(s *RegisterStatus) {
		s.AffectsStability = affectsStability
	}
}

func GetRegisterStatus(code int32) *RegisterStatus {
	initializersMu.RLock()
	defer initializersMu.RUnlock()

	registerStatus, ok := initializers[code]
	if ok {
		return registerStatus
	}

	return &RegisterStatus{
		Code:             code,
		Message:          DefaultErrorMsg,
		AffectsStability: DefaultAffectsStability,
	}
}
