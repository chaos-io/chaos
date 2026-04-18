package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var goIdentPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

const (
	generatedPackageName = "errcode"
	errorxImportPath     = "github.com/chaos-io/chaos/errorx"
)

//go:generate go run . -out ./testdata/generated ./testdata

type ErrorDefinition struct {
	Code             int    `yaml:"code"`
	Name             string `yaml:"name"`
	Message          string `yaml:"message,omitempty"`
	Description      string `yaml:"description,omitempty"`
	AffectsStability *bool  `yaml:"affectsStability,omitempty"`
}

type ErrorCodeFile struct {
	AppCode   int               `yaml:"appCode"`
	BizCode   int               `yaml:"bizCode"`
	ErrorCode []ErrorDefinition `yaml:"errorCode"`
}

type inputFile struct {
	Path     string
	BizName  string
	FileName string
	Spec     ErrorCodeFile
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	files, outputDir, err := parseArgs(args, stderr)
	if err != nil {
		return err
	}

	outputPaths, err := generateFromFiles(files, outputDir)
	if err != nil {
		return err
	}

	for _, outputPath := range outputPaths {
		_, _ = fmt.Fprintln(stdout, outputPath)
	}
	return nil
}

func parseArgs(args []string, stderr io.Writer) ([]string, string, error) {
	flagSet := flag.NewFlagSet("gen_error_code", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	outputDir := "."
	flagSet.StringVar(&outputDir, "out", ".", "output directory")

	if err := flagSet.Parse(args); err != nil {
		printUsage(stderr)
		return nil, "", err
	}

	files, err := collectYAMLFiles(flagSet.Args())
	if err != nil {
		printUsage(stderr)
		return nil, "", err
	}

	return files, outputDir, nil
}

func printUsage(w io.Writer) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintln(w, "usage: gen_error_code [-out output-dir] <yaml-file-or-dir> [more-yaml-files-or-dirs...]")
	_, _ = fmt.Fprintln(w, "yaml format: appCode, bizCode, errorCode")
	_, _ = fmt.Fprintln(w, "example: go run ./gen_error_code -out ./generated ./biz_errors")
}

func collectYAMLFiles(inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("at least one yaml file or directory is required")
	}

	var files []string
	for _, input := range inputs {
		info, err := os.Stat(input)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			if !isYAMLFile(input) {
				return nil, fmt.Errorf("unsupported file type: %s", input)
			}
			files = append(files, input)
			continue
		}

		count := 0
		err = filepath.WalkDir(input, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || !isYAMLFile(path) {
				return nil
			}
			files = append(files, path)
			count++
			return nil
		})
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("no yaml files found in %s", input)
		}
	}

	sort.Strings(files)
	return files, nil
}

func generateFromFiles(files []string, outputDir string) ([]string, error) {
	inputFiles := make([]inputFile, 0, len(files))
	for _, path := range files {
		spec, err := loadSpec(path)
		if err != nil {
			return nil, err
		}
		if err := validateSpec(path, spec); err != nil {
			return nil, err
		}
		bizName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		inputFiles = append(inputFiles, inputFile{
			Path:     path,
			BizName:  bizName,
			FileName: sanitizeFileName(bizName),
			Spec:     spec,
		})
	}
	if err := validateInputFiles(inputFiles); err != nil {
		return nil, err
	}

	outputPaths := make([]string, 0, len(inputFiles))
	for _, inputFile := range inputFiles {
		outputPath, err := generateGoCode(inputFile, outputDir)
		if err != nil {
			return nil, err
		}
		outputPaths = append(outputPaths, outputPath)
	}
	registerAllPath, err := generateRegisterAll(inputFiles, outputDir)
	if err != nil {
		return nil, err
	}
	outputPaths = append(outputPaths, registerAllPath)
	return outputPaths, nil
}

func loadSpec(path string) (ErrorCodeFile, error) {
	var spec ErrorCodeFile

	data, err := os.ReadFile(path)
	if err != nil {
		return spec, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return spec, fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return spec, nil
}

func validateSpec(path string, spec ErrorCodeFile) error {
	if spec.AppCode < 1 || spec.AppCode > 9 {
		return fmt.Errorf("%s: appCode must be in [1,9]", path)
	}
	if spec.BizCode < 1 || spec.BizCode > 9999 {
		return fmt.Errorf("%s: bizCode must be in [1,9999]", path)
	}
	if len(spec.ErrorCode) == 0 {
		return fmt.Errorf("%s: errorCode cannot be empty", path)
	}

	seenNames := make(map[string]struct{}, len(spec.ErrorCode))
	seenCodes := make(map[int]struct{}, len(spec.ErrorCode))
	for _, errDef := range spec.ErrorCode {
		if !goIdentPattern.MatchString(errDef.Name) {
			return fmt.Errorf("%s: invalid error name %q", path, errDef.Name)
		}
		if errDef.Code < 1 || errDef.Code > 9999 {
			return fmt.Errorf("%s: invalid error code for %s: must be in [1,9999]", path, errDef.Name)
		}
		if _, ok := seenNames[errDef.Name]; ok {
			return fmt.Errorf("%s: duplicate error name %q", path, errDef.Name)
		}
		if _, ok := seenCodes[errDef.Code]; ok {
			return fmt.Errorf("%s: duplicate error sub-code %d", path, errDef.Code)
		}
		seenNames[errDef.Name] = struct{}{}
		seenCodes[errDef.Code] = struct{}{}
	}

	return nil
}

func validateInputFiles(inputFiles []inputFile) error {
	seenFiles := make(map[string]string, len(inputFiles))
	seenNames := make(map[string]string)
	seenCodes := make(map[int]string)

	for _, inputFile := range inputFiles {
		outputName := inputFile.FileName + ".go"
		if previousPath, ok := seenFiles[outputName]; ok {
			return fmt.Errorf("output file conflict: %s and %s both map to %s", previousPath, inputFile.Path, outputName)
		}
		seenFiles[outputName] = inputFile.Path

		for _, errDef := range inputFile.Spec.ErrorCode {
			if previousPath, ok := seenNames[errDef.Name]; ok {
				return fmt.Errorf("duplicate error name %q across files: %s and %s", errDef.Name, previousPath, inputFile.Path)
			}
			seenNames[errDef.Name] = inputFile.Path

			fullCode := composeErrorCode(inputFile.Spec.AppCode, inputFile.Spec.BizCode, errDef.Code)
			location := fmt.Sprintf("%s:%s", inputFile.Path, errDef.Name)
			if previousLocation, ok := seenCodes[fullCode]; ok {
				return fmt.Errorf("duplicate full error code %d: %s and %s", fullCode, previousLocation, location)
			}
			seenCodes[fullCode] = location
		}
	}

	return nil
}

func generateGoCode(inputFile inputFile, outputDir string) (string, error) {
	if strings.TrimSpace(outputDir) == "" {
		outputDir = "."
	}

	source, err := buildGoCode(inputFile.BizName, inputFile.Spec)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(outputDir, inputFile.FileName+".go")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(outputPath, source, 0o644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func generateRegisterAll(inputFiles []inputFile, outputDir string) (string, error) {
	if strings.TrimSpace(outputDir) == "" {
		outputDir = "."
	}

	source, err := buildRegisterAllCode(inputFiles)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(outputDir, "register_all.go")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(outputPath, source, 0o644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func buildGoCode(bizName string, spec ErrorCodeFile) ([]byte, error) {
	var buf bytes.Buffer
	registerFuncName := registerFuncName(bizName)

	_, _ = fmt.Fprintf(&buf, `// Code generated by gen_error_code. DO NOT EDIT.
// biz: %s

package %s

import (
	"%s"
)

const (
`, bizName, generatedPackageName, errorxImportPath)

	for _, errDef := range spec.ErrorCode {
		writeConstBlock(&buf, spec.AppCode, spec.BizCode, errDef)
	}

	_, _ = fmt.Fprintf(&buf, ")\n\nfunc %s() error {\n", registerFuncName)
	for _, errDef := range spec.ErrorCode {
		writeRegistration(&buf, errDef)
	}
	_, _ = fmt.Fprintf(&buf, "\treturn nil\n}\n\n")
	for _, errDef := range spec.ErrorCode {
		writeFunctions(&buf, errDef)
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated code: %w", err)
	}
	return source, nil
}

func buildRegisterAllCode(inputFiles []inputFile) ([]byte, error) {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, `// Code generated by gen_error_code. DO NOT EDIT.

package %s

func RegisterAll() error {
`, generatedPackageName)
	for _, inputFile := range inputFiles {
		_, _ = fmt.Fprintf(
			&buf,
			"\tif err := %s(); err != nil {\n\t\treturn err\n\t}\n",
			registerFuncName(inputFile.BizName),
		)
	}
	_, _ = fmt.Fprintf(&buf, "\treturn nil\n}\n")

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format register_all code: %w", err)
	}
	return source, nil
}

func writeConstBlock(buf *bytes.Buffer, appCode int, bizCode int, errDef ErrorDefinition) {
	comment := ""
	if errDef.Description != "" {
		comment = " // " + errDef.Description
	}

	privateName := lowerFirst(errDef.Name)
	_, _ = fmt.Fprintf(
		buf,
		"\t%sCode = %d%s\n\t%sMessage = %q\n\t%sAffectsStability = %t\n",
		errDef.Name,
		composeErrorCode(appCode, bizCode, errDef.Code),
		comment,
		privateName,
		errDef.Message,
		privateName,
		errDef.affectsStability(),
	)
}

func writeRegistration(buf *bytes.Buffer, errDef ErrorDefinition) {
	privateName := lowerFirst(errDef.Name)
	_, _ = fmt.Fprintf(
		buf,
		"\tif err := errorx.Register(\n\t\t%sCode,\n\t\t%sMessage,\n\t\terrorx.WithAffectsStability(%sAffectsStability),\n\t); err != nil {\n\t\treturn err\n\t}\n",
		errDef.Name,
		privateName,
		privateName,
	)
}

func writeFunctions(buf *bytes.Buffer, errDef ErrorDefinition) {
	_, _ = fmt.Fprintf(
		buf,
		"func New%s(opts ...errorx.Option) error {\n\treturn errorx.NewByCode(%sCode, opts...)\n}\n\n",
		errDef.Name,
		errDef.Name,
	)
	_, _ = fmt.Fprintf(
		buf,
		"func Is%s(err error) bool {\n\tstatusErr, ok := errorx.FromStatus(err)\n\treturn ok && statusErr.Code() == %sCode\n}\n\n",
		errDef.Name,
		errDef.Name,
	)
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.NewReplacer("-", "_", ".", "_", "/", "_", " ", "_").Replace(name)
	if name == "" {
		return "error_code"
	}
	if name[0] >= '0' && name[0] <= '9' {
		return "biz_" + name
	}
	return name
}

func (e ErrorDefinition) affectsStability() bool {
	if e.AffectsStability != nil {
		return *e.AffectsStability
	}
	return true
}

func registerFuncName(bizName string) string {
	return "register" + upperCamelName(sanitizeFileName(bizName))
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func upperCamelName(s string) string {
	if s == "" {
		return "ErrorCode"
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == '/' || r == ' '
	})
	if len(parts) == 0 {
		return "ErrorCode"
	}

	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(part[1:])
		}
	}
	if b.Len() == 0 {
		return "ErrorCode"
	}
	return b.String()
}

func isYAMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

func composeErrorCode(appCode int, bizCode int, subCode int) int {
	return appCode*100000000 + bizCode*10000 + subCode
}
