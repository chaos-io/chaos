package errorx

import (
	"errors"
	"strings"
	"testing"
)

func TestWithStackKeepsErrorAndAddsTrace(t *testing.T) {
	err := WithStack(errors.New("boom"))
	if err == nil {
		t.Fatal("WithStack returned nil")
	}

	var stackTracer StackTracer
	if !errors.As(err, &stackTracer) {
		t.Fatalf("expected stack tracer, got %T", err)
	}
	if stackTracer.StackTrace() == "" {
		t.Fatal("expected stack trace to be present")
	}
}

func TestErrorWithoutStack(t *testing.T) {
	err := New("boom")
	if got := ErrorWithoutStack(err); !strings.Contains(got, "boom") {
		t.Fatalf("unexpected error text: %q", got)
	}

	plain := errors.New("plain error")
	if got := ErrorWithoutStack(plain); got != "plain error" {
		t.Fatalf("unexpected plain error text: %q", got)
	}
}

func TestRegisterRejectsDuplicateCode(t *testing.T) {
	const code int32 = 912340001

	Register(code, "first")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected duplicate register panic")
		}
	}()
	Register(code, "second")
}

func TestWithStatusAsStatus(t *testing.T) {
	Register(912340002, "with status as")

	err := NewByCode(912340002)

	var target Status
	if !errors.As(err, &target) {
		t.Fatalf("expected status target, got %T", err)
	}
	if target.Code() != 912340002 {
		t.Fatalf("unexpected status code: %d", target.Code())
	}
}
