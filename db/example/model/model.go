package model

import (
	"github.com/pkg/errors"
)

// TemplatePath is the path to the entity gotemplate file.
const TemplatePath = "generates/ENTITY_model.go.tmpl"

type Model struct {
	Name           string
	LowerCamelName string
}

func (m Model) Generate(path string, service *data.Service) (string, error) {
	if path != TemplatePath {
		return "", errors.Errorf("cannot render unknown file: %q", path)
	}

	return "", nil
}
