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
	if val := reflect.ValueOf(value); val.IsZero() {
		return nil, nil
	}

	return jsoniter.ConfigFastest.MarshalToString(value)
}

type JSONScanner struct{}

func (s JSONScanner) Scan(value, src any) error {
	if src == nil {
		return nil
	}
	if val := reflect.ValueOf(value); !val.IsValid() || (val.CanAddr() && val.IsNil()) {
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return jsoniter.ConfigFastest.Unmarshal(v, value)
	case string:
		return jsoniter.ConfigFastest.Unmarshal([]byte(v), value)
	default:
		return fmt.Errorf("failed to unmarshall JSONB value: %v", src)
	}
}

type JSONDbDataType struct{}

func (JSONDbDataType) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return gormDBJSONType(db)
}

func gormDBJSONType(db *gorm.DB) string {
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
