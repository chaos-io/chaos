package generator

type Generator interface {
	Generate(templateFile, generateFile string) error
}
