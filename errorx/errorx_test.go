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

func TestRegisterIsIdempotentForSameDefinition(t *testing.T) {
	const code int32 = 912340001

	if err := Register(code, "first", WithAffectsStability(false)); err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}
	if err := Register(code, "first", WithAffectsStability(false)); err != nil {
		t.Fatalf("second Register returned error: %v", err)
	}
}

func TestRegisterReturnsConflictForDifferentDefinition(t *testing.T) {
	const code int32 = 912340003

	if err := Register(code, "first", WithAffectsStability(false)); err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}

	err := Register(code, "second", WithAffectsStability(true))
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !errors.Is(err, ErrRegisterConflict) {
		t.Fatalf("expected ErrRegisterConflict, got %v", err)
	}
}

func TestWithStatusAsStatus(t *testing.T) {
	if err := Register(912340002, "with status as"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	err := NewByCode(912340002)

	var target Status
	if !errors.As(err, &target) {
		t.Fatalf("expected status target, got %T", err)
	}
	if target.Code() != 912340002 {
		t.Fatalf("unexpected status code: %d", target.Code())
	}
}

func TestGetRegisteredStatusReturnsCopy(t *testing.T) {
	const code int32 = 912340004

	if err := Register(code, "immutable", WithAffectsStability(false)); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	got := GetRegisteredStatus(code)
	got.Message = "changed"
	got.AffectsStability = true

	again := GetRegisteredStatus(code)
	if again.Message != "immutable" {
		t.Fatalf("register status should not be mutated externally: %+v", again)
	}
	if again.AffectsStability {
		t.Fatalf("register stability should not be mutated externally: %+v", again)
	}
}
