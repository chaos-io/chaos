package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestDefinitionNewCreatesCodedError(t *testing.T) {
	def := Define(600121001, "task {task_id} not found", CountInSLA(false), HTTPStatus(http.StatusNotFound))

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
	if def.CountInSLA {
		t.Fatalf("definition should keep explicit SLA flag: %+v", def)
	}
	if got.StatusCode() != http.StatusNotFound {
		t.Fatalf("StatusCode() = %d, want %d", got.StatusCode(), http.StatusNotFound)
	}
}

func TestDefinitionDefaultsHTTPStatusToInternalServerError(t *testing.T) {
	def := Define(600121006, "unknown")
	if got := def.StatusCode(); got != http.StatusInternalServerError {
		t.Fatalf("StatusCode() = %d, want %d", got, http.StatusInternalServerError)
	}
	if got := def.New().(*Error).StatusCode(); got != http.StatusInternalServerError {
		t.Fatalf("StatusCode() = %d, want %d", got, http.StatusInternalServerError)
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

func TestErrorStringIsSafeForClientResponse(t *testing.T) {
	def := Define(600121004, "get task failed")
	err := def.Wrap(errors.New("record not found"))

	if got := err.Error(); got != "get task failed" {
		t.Fatalf("Error should return client-safe message, got %q", got)
	}

	detail := fmt.Sprintf("%+v", err)
	for _, want := range []string{
		"code=600121004",
		"message=get task failed",
		"cause=record not found",
	} {
		if !strings.Contains(detail, want) {
			t.Fatalf("formatted detail missing %q:\n%s", want, detail)
		}
	}
	if strings.ContainsAny(detail, "{}\"") {
		t.Fatalf("formatted detail should be readable text, not JSON:\n%s", detail)
	}
	if strings.Contains(detail, "stack=") {
		t.Fatalf("formatted detail should not include stack by default:\n%s", detail)
	}
}

func TestErrorMarshalJSONIsSafeForClientResponse(t *testing.T) {
	def := Define(600121005, "get task failed")
	err := def.Wrap(errors.New("record not found"))

	body, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("MarshalJSON returned error: %v", marshalErr)
	}

	if got, want := string(body), `{"code":"600121005","message":"get task failed"}`; got != want {
		t.Fatalf("unexpected json body: got %s want %s", got, want)
	}
	if strings.Contains(string(body), "cause") {
		t.Fatalf("json body should not expose cause: %s", body)
	}
	if strings.Contains(string(body), "stack") {
		t.Fatalf("json body should not expose stack: %s", body)
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
	first := Define(912340001, "first", CountInSLA(false))
	same := Define(912340001, "first", CountInSLA(false))
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
}
