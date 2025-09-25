package tracer

var _ Span = (*noopSpan)(nil)

//go:generate mockgen -destination=mocks/span.go -package=mocks . Span
type Span interface {
	SetCallType(callType string)
}

type noopSpan struct{}

func (n noopSpan) SetCallType(callType string) {}
