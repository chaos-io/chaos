package messaging

import "errors"

var (
	ErrNilConfig                 = errors.New("messaging: config is nil")
	ErrNilClient                 = errors.New("messaging: client is nil")
	ErrNilProvider               = errors.New("messaging: queue provider is nil")
	ErrNilQueueConstructor       = errors.New("messaging: queue constructor is nil")
	ErrInvalidProviderName       = errors.New("messaging: provider name is empty")
	ErrProviderAlreadyRegistered = errors.New("messaging: provider already registered")
	ErrProviderNotRegistered     = errors.New("messaging: provider not registered")
	ErrNilSubscription           = errors.New("messaging: subscription is nil")
	ErrNilHandler                = errors.New("messaging: message handler is nil")
	ErrEmptyTopic                = errors.New("messaging: topic is empty")
)
