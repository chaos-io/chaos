package db

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSONValuer struct{}

func (v JSONValuer) Value(value any) (driver.Value, error) {
	if isNilValue(value) {
		return nil, nil
	}
	if val := reflect.ValueOf(value); val.IsZero() {
		return nil, nil
	}

	return jsoniter.ConfigFastest.MarshalToString(value)
}

type JSONScanner struct{}

func (s JSONScanner) Scan(value, src any) error {
	if src == nil || isNilValue(value) {
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return jsoniter.ConfigFastest.Unmarshal(v, value)
	case string:
		return jsoniter.ConfigFastest.Unmarshal([]byte(v), value)
	default:
		return fmt.Errorf("failed to unmarshal JSON value from %T", src)
	}
}

type JSONDbDataType struct{}

func (JSONDbDataType) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	return gormDBJSONType(db)
}

func gormDBJSONType(db *gorm.DB) string {
	if db == nil || db.Dialector == nil {
		return ""
	}

	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func isNilValue(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
