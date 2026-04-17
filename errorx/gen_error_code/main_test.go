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
}

func TestGenerateGoCodeWritesFormattedFile(t *testing.T) {
	outputDir := t.TempDir()
	path := filepath.Join(t.TempDir(), "loop-task.yaml")
	spec := Spec{
		AppCode: 6,
		BizCode: 12,
		ErrorCode: []Error{
			{Name: "TaskNotFound", Code: 1001, Message: "task not found"},
		},
	}

	outputPath, err := generateGoCode(path, spec, outputDir)
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
	if !strings.Contains(got, "package loop_task") {
		t.Fatalf("generated file missing package name:\n%s", got)
	}
	if !strings.Contains(got, "TaskNotFoundCode") || !strings.Contains(got, "600121001") {
		t.Fatalf("generated file missing error code:\n%s", got)
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

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
