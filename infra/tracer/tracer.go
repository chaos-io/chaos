package tracer

import (
	"context"

	"github.com/chaos-io/core/go/chaos/core"
)

var tracer Tracer = &noopTracer{}

//go:generate mockgen -destination=mocks/tracer.go -package=mocks . Tracer
type Tracer interface {
	// StartSpan Generate a span that automatically links to the previous span in the context.
	// The start time of the span starts counting from the call of StartSpan.
	// The generated span will be automatically written into the context.
	// Subsequent spans that need to be chained should call StartSpan based on the new context.
	StartSpan(ctx context.Context, name, spanType string, opts ...core.Options) (context.Context, Span)
	// GetSpanFromContext Get the span from the context.
	GetSpanFromContext(ctx context.Context) Span
	// Flush Force the reporting of spans in the queue.
	Flush(ctx context.Context)
	// Inject Inject the tracer into the context.
	Inject(ctx context.Context) context.Context
}

// GetTracer Get the tracer. Must call InitTracer first.
func GetTracer() Tracer {
	return tracer
}

// InitTracer Init the tracer. Must call before GetTracer.
func InitTracer(t Tracer) {
	tracer = t
}

type noopTracer struct{}

func (d *noopTracer) StartSpan(ctx context.Context, name, spanType string, opts ...core.Options) (context.Context, Span) {
	return ctx, &noopSpan{}
}

func (d *noopTracer) GetSpanFromContext(ctx context.Context) Span {
	return &noopSpan{}
}

func (d *noopTracer) Flush(ctx context.Context) {}

func (d *noopTracer) Inject(ctx context.Context) context.Context {
	return ctx
}

func (d *noopTracer) SetCallType(callType string) {}
