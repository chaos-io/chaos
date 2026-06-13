package errorx

import (
	"errors"
	"strings"
	"testing"
)

func TestDefinitionNewCreatesCodedError(t *testing.T) {
	def := Define(600121001, "task {task_id} not found", AffectsStability(false))

	err := def.New(
		WithMessageParam("task_id", "t-1"),
		WithExtra(map[string]string{"task_id": "t-1"}),
	)

	got, ok := From(err)
	if !ok {
		t.Fatalf("From returned false for %T", err)
	}
	if got.Code() != 600121001 {
		t.Fatalf("unexpected code: %d", got.Code())
	}
	if got.Message() != "task t-1 not found" {
		t.Fatalf("unexpected message: %q", got.Message())
	}
	if got.Extra()["task_id"] != "t-1" {
		t.Fatalf("unexpected extra: %+v", got.Extra())
	}
	if got.StackTrace() == "" {
		t.Fatal("expected stack trace")
	}
	if def.AffectsStability {
		t.Fatalf("definition should keep explicit stability flag: %+v", def)
	}
}

func TestDefinitionWrapKeepsCauseAndSupportsMatching(t *testing.T) {
	def := Define(600121002, "db failed")
	cause := errors.New("db down")

	err := def.Wrap(cause, WithExtraMessage("retry later"))

	if !errors.Is(err, cause) {
		t.Fatalf("wrapped error should keep cause: %v", err)
	}
	if !def.Is(err) {
		t.Fatalf("definition should match wrapped error: %v", err)
	}
	if !Is(err, 600121002) {
		t.Fatalf("Is should match wrapped error: %v", err)
	}
	if def.Wrap(nil) != nil {
		t.Fatal("Wrap(nil) should return nil")
	}
	if got := ErrorWithoutStack(err); strings.Contains(got, "stack=") {
		t.Fatalf("ErrorWithoutStack should remove stack: %q", got)
	}
}

func TestCodeOfAndWithoutStack(t *testing.T) {
	def := Define(600121003, "plain")
	err := def.New(WithoutStack())

	code, ok := CodeOf(err)
	if !ok {
		t.Fatal("CodeOf returned false")
	}
	if code != 600121003 {
		t.Fatalf("unexpected code: %d", code)
	}

	got, ok := From(err)
	if !ok {
		t.Fatal("From returned false")
	}
	if got.StackTrace() != "" {
		t.Fatalf("WithoutStack should skip stack, got %q", got.StackTrace())
	}
}

func TestRegisterDefinitions(t *testing.T) {
	first := Define(912340001, "first", AffectsStability(false))
	same := Define(912340001, "first", AffectsStability(false))
	conflict := Define(912340001, "second")

	if err := Register(first); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if err := Register(same); err != nil {
		t.Fatalf("Register should be idempotent: %v", err)
	}

	err := Register(conflict)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !errors.Is(err, ErrRegisterConflict) {
		t.Fatalf("expected ErrRegisterConflict, got %v", err)
	}

	registered, ok := Lookup(912340001)
	if !ok {
		t.Fatal("expected registered definition")
	}
	registered.Message = "changed"

	again, ok := Lookup(912340001)
	if !ok || again.Message != "first" {
		t.Fatalf("Lookup should return a copy, got %+v", again)
	}
}
