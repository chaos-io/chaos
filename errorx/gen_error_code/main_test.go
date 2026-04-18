package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "loop_task.yaml")
	writeFile(t, path, `
app_code: 6
biz_code: 12
error_code:
  - name: TaskNotFound
    code: 1001
    message: task not found
    no_affect_stability: true
`)

	spec, err := loadSpec(path)
	if err != nil {
		t.Fatalf("loadSpec returned error: %v", err)
	}
	if spec.AppCode != 6 || spec.BizCode != 12 {
		t.Fatalf("unexpected spec codes: %+v", spec)
	}
	if len(spec.ErrorCode) != 1 || spec.ErrorCode[0].Name != "TaskNotFound" {
		t.Fatalf("unexpected errors: %+v", spec)
	}
	if !spec.ErrorCode[0].noAffectStability() {
		t.Fatalf("expected no_affect_stability to be loaded: %+v", spec.ErrorCode[0])
	}
}

func TestGenerateGoCodeWritesFormattedFile(t *testing.T) {
	outputDir := t.TempDir()
	specFile := specFile{
		Path:     filepath.Join(t.TempDir(), "loop-task.yaml"),
		BizName:  "loop-task",
		FileName: "loop_task",
		Spec: Spec{
			AppCode: 6,
			BizCode: 12,
			ErrorCode: []Error{
				{Name: "TaskNotFound", Code: 1001, Message: "task not found"},
			},
		},
	}

	outputPath, err := generateGoCode(specFile, outputDir)
	if err != nil {
		t.Fatalf("generateGoCode returned error: %v", err)
	}

	if want := filepath.Join(outputDir, "loop_task.go"); outputPath != want {
		t.Fatalf("unexpected output path: got %s want %s", outputPath, want)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	got := string(content)
	if !strings.Contains(got, "package errcode") {
		t.Fatalf("generated file missing package name:\n%s", got)
	}
	if !strings.Contains(got, "TaskNotFoundCode") || !strings.Contains(got, "600121001") {
		t.Fatalf("generated file missing error code:\n%s", got)
	}
	if !strings.Contains(got, "\"chaos-io/chaos/errorx\"") {
		t.Fatalf("generated file missing errorx import path:\n%s", got)
	}
	if !strings.Contains(got, "errorx.Register(") {
		t.Fatalf("generated file missing register call:\n%s", got)
	}
	if !strings.Contains(got, "func NewTaskNotFound(opts ...errorx.Option) error") {
		t.Fatalf("generated file missing constructor:\n%s", got)
	}
	if !strings.Contains(got, "func IsTaskNotFound(err error) bool") {
		t.Fatalf("generated file missing matcher:\n%s", got)
	}
}

func TestGenerateFromFilesCreatesOneGoFilePerYAML(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	writeFile(t, filepath.Join(inputDir, "loop_task.yaml"), `
app_code: 6
biz_code: 12
error_code:
  - name: TaskNotFound
    code: 1001
    message: task not found
`)
	writeFile(t, filepath.Join(inputDir, "user.yaml"), `
app_code: 6
biz_code: 13
error_code:
  - name: UserNotFound
    code: 1001
    message: user not found
`)

	files, err := collectYAMLFiles([]string{inputDir})
	if err != nil {
		t.Fatalf("collectYAMLFiles returned error: %v", err)
	}

	outputs, err := generateFromFiles(files, outputDir)
	if err != nil {
		t.Fatalf("generateFromFiles returned error: %v", err)
	}
	if len(outputs) != 2 {
		t.Fatalf("unexpected output count: %v", outputs)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "loop_task.go")); err != nil {
		t.Fatalf("loop_task.go not generated: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "user.go")); err != nil {
		t.Fatalf("user.go not generated: %v", err)
	}
}

func TestRunPrintsAllGeneratedFiles(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	writeFile(t, filepath.Join(inputDir, "loop_task.yaml"), `
app_code: 6
biz_code: 12
error_code:
  - name: TaskNotFound
    code: 1001
    message: task not found
`)
	writeFile(t, filepath.Join(inputDir, "user.yaml"), `
app_code: 6
biz_code: 13
error_code:
  - name: UserNotFound
    code: 1001
    message: user not found
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run([]string{"-out", outputDir, inputDir}, &stdout, &stderr); err != nil {
		t.Fatalf("run returned error: %v; stderr=%s", err, stderr.String())
	}

	lines := strings.Fields(strings.TrimSpace(stdout.String()))
	if len(lines) != 2 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}

func TestValidateSpecChecksNineDigitLayout(t *testing.T) {
	err := validateSpec("biz.yaml", Spec{
		AppCode: 10,
		BizCode: 12,
		ErrorCode: []Error{
			{Name: "TaskNotFound", Code: 1001},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "app_code must be in [1,9]") {
		t.Fatalf("unexpected validation error: %v", err)
	}

	err = validateSpec("biz.yaml", Spec{
		AppCode: 6,
		BizCode: 10000,
		ErrorCode: []Error{
			{Name: "TaskNotFound", Code: 1001},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "biz_code must be in [1,9999]") {
		t.Fatalf("unexpected validation error: %v", err)
	}

	err = validateSpec("biz.yaml", Spec{
		AppCode: 6,
		BizCode: 12,
		ErrorCode: []Error{
			{Name: "TaskNotFound", Code: 10000},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "must be in [1,9999]") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateSpecFilesRejectsDuplicateErrorNames(t *testing.T) {
	err := validateSpecFiles([]specFile{
		{
			Path:     "order.yaml",
			BizName:  "order",
			FileName: "order",
			Spec: Spec{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []Error{
					{Name: "NotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "payment.yaml",
			BizName:  "payment",
			FileName: "payment",
			Spec: Spec{
				AppCode: 6,
				BizCode: 13,
				ErrorCode: []Error{
					{Name: "NotFound", Code: 1002},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), `duplicate error name "NotFound"`) {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateSpecFilesRejectsDuplicateFullCodes(t *testing.T) {
	err := validateSpecFiles([]specFile{
		{
			Path:     "order.yaml",
			BizName:  "order",
			FileName: "order",
			Spec: Spec{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []Error{
					{Name: "OrderNotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "payment.yaml",
			BizName:  "payment",
			FileName: "payment",
			Spec: Spec{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []Error{
					{Name: "PaymentNotFound", Code: 1001},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "duplicate full error code 600121001") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateSpecFilesRejectsOutputConflicts(t *testing.T) {
	err := validateSpecFiles([]specFile{
		{
			Path:     "foo-bar.yaml",
			BizName:  "foo-bar",
			FileName: "foo_bar",
			Spec: Spec{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []Error{
					{Name: "FooBarNotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "foo_bar.yaml",
			BizName:  "foo_bar",
			FileName: "foo_bar",
			Spec: Spec{
				AppCode: 6,
				BizCode: 13,
				ErrorCode: []Error{
					{Name: "FooBarAlreadyExists", Code: 1001},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "output file conflict") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
