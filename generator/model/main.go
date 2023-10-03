package main

import (
	_ "embed"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/chaos-io/chaos/core/strcase"
)

const GoPackageName = "internal"

//go:embed templates/entity.go.tmpl
var templateStructFile string

//go:embed templates/entity_model.go.tmpl
var templateModelFile string

func main() {
	Generator()
}

func Generator() {
	m := InitModel()
	log.Printf("get the model: %+v\n", m)

	generateStructFile := path.Join(GoPackageName, strcase.ToSnake(m.Name)+".go")
	if err := m.Generate(templateStructFile, generateStructFile); err != nil {
		log.Fatalf("generate struct file error: %v", err)
	}
	log.Printf("generate struct file into %s", generateStructFile)

	generateModelFile := path.Join(GoPackageName+"/model", strcase.ToSnake(m.Name)+"_model.go")
	if err := m.Generate(templateModelFile, generateModelFile); err != nil {
		log.Fatalf("generate model file error: %v", err)
		return
	}
	log.Printf("generate model file into %s", generateModelFile)
}

func InitModel() *Model {
	if len(os.Args) < 2 {
		log.Fatalf("can't get enough arguments, length(%v), arguments(%v)\n", len(os.Args), os.Args)
	}
	log.Printf("get os.Args: %v\n", os.Args)
	name := strcase.ToCamel(os.Args[1])
	return &Model{
		Name:           name,
		LowerCamelName: strcase.ToLowerCamel(name),
		GoPackageName:  GoPackageName,
	}
}

type Model struct {
	Name           string
	LowerCamelName string
	GoPackageName  string
}

func (m *Model) Generate(templateFile, generateFile string) error {
	tmpl, err := template.New(m.Name).Parse(templateFile)
	if err != nil {
		return err
	}

	file, err := os.Create(generateFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, m); err != nil {
		return err
	}

	return nil
}
