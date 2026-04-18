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
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task not found
    affectsStability: false
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
	if spec.ErrorCode[0].affectsStability() {
		t.Fatalf("expected affectsStability to be loaded: %+v", spec.ErrorCode[0])
	}
}

func TestGenerateGoCodeWritesFormattedFile(t *testing.T) {
	outputDir := t.TempDir()
	inputFile := inputFile{
		Path:     filepath.Join(t.TempDir(), "loop-task.yaml"),
		BizName:  "loop-task",
		FileName: "loop_task",
		Spec: ErrorCodeFile{
			AppCode: 6,
			BizCode: 12,
			ErrorCode: []ErrorDefinition{
				{Name: "TaskNotFound", Code: 1001, Message: "task not found"},
			},
		},
	}

	outputPath, err := generateGoCode(inputFile, outputDir)
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
	if !strings.Contains(got, "func registerLoopTask() error") {
		t.Fatalf("generated file missing register helper:\n%s", got)
	}
	if !strings.Contains(got, "if err := errorx.Register(") {
		t.Fatalf("generated file missing register error handling:\n%s", got)
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
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task not found
`)
	writeFile(t, filepath.Join(inputDir, "user.yaml"), `
appCode: 6
bizCode: 13
errorCode:
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
	if len(outputs) != 3 {
		t.Fatalf("unexpected output count: %v", outputs)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "loop_task.go")); err != nil {
		t.Fatalf("loop_task.go not generated: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "user.go")); err != nil {
		t.Fatalf("user.go not generated: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "register_all.go")); err != nil {
		t.Fatalf("register_all.go not generated: %v", err)
	}
}

func TestRunPrintsAllGeneratedFiles(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	writeFile(t, filepath.Join(inputDir, "loop_task.yaml"), `
appCode: 6
bizCode: 12
errorCode:
  - name: TaskNotFound
    code: 1001
    message: task not found
`)
	writeFile(t, filepath.Join(inputDir, "user.yaml"), `
appCode: 6
bizCode: 13
errorCode:
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
	if len(lines) != 3 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}

func TestValidateSpecChecksNineDigitLayout(t *testing.T) {
	err := validateSpec("biz.yaml", ErrorCodeFile{
		AppCode: 10,
		BizCode: 12,
		ErrorCode: []ErrorDefinition{
			{Name: "TaskNotFound", Code: 1001},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "appCode must be in [1,9]") {
		t.Fatalf("unexpected validation error: %v", err)
	}

	err = validateSpec("biz.yaml", ErrorCodeFile{
		AppCode: 6,
		BizCode: 10000,
		ErrorCode: []ErrorDefinition{
			{Name: "TaskNotFound", Code: 1001},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "bizCode must be in [1,9999]") {
		t.Fatalf("unexpected validation error: %v", err)
	}

	err = validateSpec("biz.yaml", ErrorCodeFile{
		AppCode: 6,
		BizCode: 12,
		ErrorCode: []ErrorDefinition{
			{Name: "TaskNotFound", Code: 10000},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "must be in [1,9999]") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateInputFilesRejectsDuplicateErrorNames(t *testing.T) {
	err := validateInputFiles([]inputFile{
		{
			Path:     "order.yaml",
			BizName:  "order",
			FileName: "order",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []ErrorDefinition{
					{Name: "NotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "payment.yaml",
			BizName:  "payment",
			FileName: "payment",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 13,
				ErrorCode: []ErrorDefinition{
					{Name: "NotFound", Code: 1002},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), `duplicate error name "NotFound"`) {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateInputFilesRejectsDuplicateFullCodes(t *testing.T) {
	err := validateInputFiles([]inputFile{
		{
			Path:     "order.yaml",
			BizName:  "order",
			FileName: "order",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []ErrorDefinition{
					{Name: "OrderNotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "payment.yaml",
			BizName:  "payment",
			FileName: "payment",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []ErrorDefinition{
					{Name: "PaymentNotFound", Code: 1001},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "duplicate full error code 600121001") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateInputFilesRejectsOutputConflicts(t *testing.T) {
	err := validateInputFiles([]inputFile{
		{
			Path:     "foo-bar.yaml",
			BizName:  "foo-bar",
			FileName: "foo_bar",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 12,
				ErrorCode: []ErrorDefinition{
					{Name: "FooBarNotFound", Code: 1001},
				},
			},
		},
		{
			Path:     "foo_bar.yaml",
			BizName:  "foo_bar",
			FileName: "foo_bar",
			Spec: ErrorCodeFile{
				AppCode: 6,
				BizCode: 13,
				ErrorCode: []ErrorDefinition{
					{Name: "FooBarAlreadyExists", Code: 1001},
				},
			},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "output file conflict") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestBuildRegisterAllCode(t *testing.T) {
	source, err := buildRegisterAllCode([]inputFile{
		{BizName: "loop_task"},
		{BizName: "user"},
	})
	if err != nil {
		t.Fatalf("buildRegisterAllCode returned error: %v", err)
	}

	got := string(source)
	if !strings.Contains(got, "func RegisterAll() error") {
		t.Fatalf("generated register_all missing RegisterAll:\n%s", got)
	}
	if !strings.Contains(got, "if err := registerLoopTask(); err != nil") {
		t.Fatalf("generated register_all missing loop_task registration:\n%s", got)
	}
	if !strings.Contains(got, "if err := registerUser(); err != nil") {
		t.Fatalf("generated register_all missing user registration:\n%s", got)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
