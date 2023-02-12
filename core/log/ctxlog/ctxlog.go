package ctxlog

import (
	"context"
	"fmt"

	log2 "github.com/chaos-io/chaos/core/log"
)

type ctxKey struct{}

// ContextFields returns log.Fields bound with ctx.
// If no fields are bound, it returns nil.
func ContextFields(ctx context.Context) []log2.Field {
	fs, _ := ctx.Value(ctxKey{}).([]log2.Field)
	return fs
}

// WithFields returns a new context that is bound with given fields and based
// on parent ctx.
func WithFields(ctx context.Context, fields ...log2.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, mergeFields(ContextFields(ctx), fields))
}

// Trace logs at Trace log level using fields both from arguments and ones that
// are bound to ctx.
func Trace(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Trace(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Debug logs at Debug log level using fields both from arguments and ones that
// are bound to ctx.
func Debug(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Debug(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Info logs at Info log level using fields both from arguments and ones that
// are bound to ctx.
func Info(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Info(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Warn logs at Warn log level using fields both from arguments and ones that
// are bound to ctx.
func Warn(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Warn(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Error logs at Error log level using fields both from arguments and ones that
// are bound to ctx.
func Error(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Error(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Fatal logs at Fatal log level using fields both from arguments and ones that
// are bound to ctx.
func Fatal(ctx context.Context, l log2.Logger, msg string, fields ...log2.Field) {
	log2.AddCallerSkip(l, 1).Fatal(msg, mergeFields(ContextFields(ctx), fields)...)
}

// Tracef logs at Trace log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Tracef(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Trace(msg, ContextFields(ctx)...)
}

// Debugf logs at Debug log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Debugf(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Debug(msg, ContextFields(ctx)...)
}

// Infof logs at Info log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Infof(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Info(msg, ContextFields(ctx)...)
}

// Warnf logs at Warn log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Warnf(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Warn(msg, ContextFields(ctx)...)
}

// Errorf logs at Error log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Errorf(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Error(msg, ContextFields(ctx)...)
}

// Fatalf logs at Fatal log level using fields that are bound to ctx.
// The message is formatted using provided arguments.
func Fatalf(ctx context.Context, l log2.Logger, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log2.AddCallerSkip(l, 1).Fatal(msg, ContextFields(ctx)...)
}

func mergeFields(a, b []log2.Field) []log2.Field {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	// NOTE: just append() here is unsafe. If a caller passed slice of fields
	// followed by ... with capacity greater than length, then simultaneous
	// logging will lead to a data race condition.
	//
	// See https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	c := make([]log2.Field, len(a)+len(b))
	n := copy(c, a)
	copy(c[n:], b)
	return c
}
