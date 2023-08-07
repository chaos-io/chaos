package templates

const EntityModel = `
package model

import (
	"context"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/chaos-io/chaos/db"
	"github.com/chaos-io/chaos/logs"
)

var {{.LowerCamelName}}Model *{{.Name}}Model
var {{.LowerCamelName}}ModelOnce sync.Once

type {{.Name}}Model struct {
	DB *gorm.DB
}

func Get{{.Name}}Model() *{{.Name}}Model {
	{{.LowerCamelName}}ModelOnce.Do(func() {
		{{.Name}}Model = New{{.Name}}Model()
	})

	return surveyModel
}

func New{{.Name}}Model() *{{.Name}}Model {
	m := &{{.Name}}Model{DB: InitDB()}
	if !m.DB.Config.DisableAutoMigrate || !d.Migrator().HasTable(&{{.GoPackageName}}.{{.Name}}{}) {
		if err := d.AutoMigrate(&{{.GoPackageName}}.{{.Name}}{}); err != nil {
			logs.Error("Init {{.Name}}Model model err: ", err)
			panic(err)
		}
	}

	return m
}

func (m *{{.Name}}Model) Create({{.LowerCamelName}} *{{.GoPackageName}}.{{.Name}}) (int64, error) {
	result := m.DB.Create({{.LowerCamelName}})
	return result.RowsAffected, result.Error
}

func (m *{{.Name}}Model) Get(id string) (*{{.GoPackageName}}.{{.Name}}, error) {
	{{.LowerCamelName}} := &{{.GoPackageName}}.{{.Name}}{}
	if err := m.DB.First(survey, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return {{.LowerCamelName}}, nil
}

func (m *{{.Name}}Model) Delete(id string) (int64, error) {
	result := m.DB.Where("id = ?", uid).Delete(&{{.GoPackageName}}.{{.Name}}{})
	return result.RowsAffected, result.Error
}

func (m *{{.Name}}Model) Update({{.LowerCamelName}} *{{.GoPackageName}}.{{.Name}}) (int64, error) {
	result := m.DB.Updates({{.LowerCamelName}})
	return result.RowsAffected, result.Error
}

func (m *{{.Name}}Model) List(filter string, fieldMask FieldMask) ([]*{{.GoPackageName}}.{{.Name}}, error) {
	var {{Plural .LowerCamelName}} []*{{.GoPackageName}}.{{.Name}}
	if err := m.DB.Find(&{{Plural .LowerCamelName}}).Error; err != nil {
		return nil, err
	}
	return {{Plural .LowerCamelName}}, nil
}

func (m *{{.Name}}Model) BatchCreate({{Plural .LowerCamelName}} ...*{{.GoPackageName}}.{{.Name}}) (int64, error) {
	result := m.DB.CreateInBatches({{Plural .LowerCamelName}}, len({{Plural .LowerCamelName}}))
	return result.RowsAffected, result.Error
}

func (m *{{.Name}}Model) BatchGet(ids ...string) (*{{.GoPackageName}}.{{.Name}}, error) {
	var {{Plural .LowerCamelName}} []{{.GoPackageName}}.{{.Name}}
	if err := m.DB.Find(&{{Plural .LowerCamelName}}, "id = ?", ids).Error; err != nil {
		return nil, err
	}
	return {{Plural .LowerCamelName}}, nil
}

func (m *{{.Name}}Model) BatchDelete(ids ...string) (int64, error) {
	result := m.DB.Where("id = ?", ids).Delete(&{{.GoPackageName}}.{{.Name}}{})
	return result.RowsAffected, result.Error
}

func (m *{{.Name}}Model) BatchUpdate({{Plural .LowerCamelName}} []*{{.GoPackageName}}.{{.Name}}) (int64, error) {
	result := m.DB.Updates({{Plural .LowerCamelName}}...)
	return result.RowsAffected, result.Error
}
`
