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
	errorxImportPath     = "chaos-io/chaos/errorx"
)

//go:generate go run . -out ./testdata/generated ./testdata

type Error struct {
	Code                 int    `yaml:"code"`
	Name                 string `yaml:"name"`
	Message              string `yaml:"message,omitempty"`
	Description          string `yaml:"description,omitempty"`
	NoAffectStability    bool   `yaml:"no_affect_stability,omitempty"`
	LegacyNoAffectStable bool   `yaml:"no_affects_stability,omitempty"`
}

type Spec struct {
	AppCode   int     `yaml:"app_code"`
	BizCode   int     `yaml:"biz_code"`
	ErrorCode []Error `yaml:"error_code"`
}

type specFile struct {
	Path     string
	BizName  string
	FileName string
	Spec     Spec
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
	_, _ = fmt.Fprintln(w, "yaml format: app_code, biz_code, error_code")
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
	specFiles := make([]specFile, 0, len(files))
	for _, path := range files {
		spec, err := loadSpec(path)
		if err != nil {
			return nil, err
		}
		if err := validateSpec(path, spec); err != nil {
			return nil, err
		}
		bizName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		specFiles = append(specFiles, specFile{
			Path:     path,
			BizName:  bizName,
			FileName: sanitizeFileName(bizName),
			Spec:     spec,
		})
	}
	if err := validateSpecFiles(specFiles); err != nil {
		return nil, err
	}

	outputPaths := make([]string, 0, len(specFiles))
	for _, specFile := range specFiles {
		outputPath, err := generateGoCode(specFile, outputDir)
		if err != nil {
			return nil, err
		}
		outputPaths = append(outputPaths, outputPath)
	}
	return outputPaths, nil
}

func loadSpec(path string) (Spec, error) {
	var spec Spec

	data, err := os.ReadFile(path)
	if err != nil {
		return spec, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return spec, fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return spec, nil
}

func validateSpec(path string, spec Spec) error {
	if spec.AppCode < 1 || spec.AppCode > 9 {
		return fmt.Errorf("%s: app_code must be in [1,9]", path)
	}
	if spec.BizCode < 1 || spec.BizCode > 9999 {
		return fmt.Errorf("%s: biz_code must be in [1,9999]", path)
	}
	if len(spec.ErrorCode) == 0 {
		return fmt.Errorf("%s: error_code cannot be empty", path)
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

func validateSpecFiles(specFiles []specFile) error {
	seenFiles := make(map[string]string, len(specFiles))
	seenNames := make(map[string]string)
	seenCodes := make(map[int]string)

	for _, specFile := range specFiles {
		outputName := specFile.FileName + ".go"
		if previousPath, ok := seenFiles[outputName]; ok {
			return fmt.Errorf("output file conflict: %s and %s both map to %s", previousPath, specFile.Path, outputName)
		}
		seenFiles[outputName] = specFile.Path

		for _, errDef := range specFile.Spec.ErrorCode {
			if previousPath, ok := seenNames[errDef.Name]; ok {
				return fmt.Errorf("duplicate error name %q across files: %s and %s", errDef.Name, previousPath, specFile.Path)
			}
			seenNames[errDef.Name] = specFile.Path

			fullCode := errorCode(specFile.Spec.AppCode, specFile.Spec.BizCode, errDef.Code)
			location := fmt.Sprintf("%s:%s", specFile.Path, errDef.Name)
			if previousLocation, ok := seenCodes[fullCode]; ok {
				return fmt.Errorf("duplicate full error code %d: %s and %s", fullCode, previousLocation, location)
			}
			seenCodes[fullCode] = location
		}
	}

	return nil
}

func generateGoCode(specFile specFile, outputDir string) (string, error) {
	if strings.TrimSpace(outputDir) == "" {
		outputDir = "."
	}

	source, err := buildGoCode(specFile.BizName, specFile.Spec)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(outputDir, specFile.FileName+".go")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(outputPath, source, 0o644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func buildGoCode(bizName string, spec Spec) ([]byte, error) {
	var buf bytes.Buffer

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

	_, _ = fmt.Fprintf(&buf, ")\n\nfunc init() {\n")
	for _, errDef := range spec.ErrorCode {
		writeRegistration(&buf, errDef)
	}
	_, _ = fmt.Fprintf(&buf, "}\n\n")
	for _, errDef := range spec.ErrorCode {
		writeFunctions(&buf, errDef)
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated code: %w", err)
	}
	return source, nil
}

func writeConstBlock(buf *bytes.Buffer, appCode int, bizCode int, errDef Error) {
	comment := ""
	if errDef.Description != "" {
		comment = " // " + errDef.Description
	}

	privateName := lowerFirst(errDef.Name)
	_, _ = fmt.Fprintf(
		buf,
		"\t%sCode = %d%s\n\t%sMessage = %q\n\t%sNoAffectStability = %t\n",
		errDef.Name,
		errorCode(appCode, bizCode, errDef.Code),
		comment,
		privateName,
		errDef.Message,
		privateName,
		errDef.noAffectStability(),
	)
}

func writeRegistration(buf *bytes.Buffer, errDef Error) {
	privateName := lowerFirst(errDef.Name)
	_, _ = fmt.Fprintf(
		buf,
		"\terrorx.Register(\n\t\t%sCode,\n\t\t%sMessage,\n\t\terrorx.WithAffectsStability(!%sNoAffectStability),\n\t)\n",
		errDef.Name,
		privateName,
		privateName,
	)
}

func writeFunctions(buf *bytes.Buffer, errDef Error) {
	_, _ = fmt.Fprintf(
		buf,
		"func New%s(opts ...errorx.Option) error {\n\treturn errorx.NewByCode(%sCode, opts...)\n}\n\n",
		errDef.Name,
		errDef.Name,
	)
	_, _ = fmt.Fprintf(
		buf,
		"func Is%s(err error) bool {\n\tstatus, ok := errorx.FromStatus(err)\n\treturn ok && status.Code() == %sCode\n}\n\n",
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

func (e Error) noAffectStability() bool {
	return e.NoAffectStability || e.LegacyNoAffectStable
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func isYAMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

func errorCode(appCode int, bizCode int, subCode int) int {
	return appCode*100000000 + bizCode*10000 + subCode
}
