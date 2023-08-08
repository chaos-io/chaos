package templates

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	path2 "path"
	"strings"
)

//go:embed entity_model.go.tmpl
var templateFile string

func InitModel() *Model {
	if len(os.Args) < 2 {
		log.Panicf("can't get enough argument (%v)", os.Args)
	}
	fmt.Printf("get os.Args: %v\n", os.Args)
	return &Model{
		Name:           os.Args[1],
		LowerCamelName: toLowerCamelCase(os.Args[1]),
		GoPackageName:  "model",
	}
}

type Model struct {
	Name           string
	LowerCamelName string
	GoPackageName  string
}

func (m Model) Generate() (string, error) {
	tmpl, err := template.New(m.Name).Parse(templateFile)
	if err != nil {
		return "", err
	}

	generateFile := path2.Join("generates/model", m.Name+"_model.go")
	file, err := os.Create(generateFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := tmpl.Execute(file, m); err != nil {
		return "", err
	}

	return generateFile, nil
}

func toLowerCamelCase(input string) string {
	words := strings.Fields(input)
	for i := range words {
		if i > 0 {
			words[i] = strings.Title(words[i])
		} else {
			words[i] = strings.ToLower(words[i])
		}
	}
	return strings.Join(words, "")
}
