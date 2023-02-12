package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	log2 "github.com/chaos-io/chaos/core/log"
)

// ZapifyLevel turns interface log level to zap log level
func ZapifyLevel(level log2.Level) zapcore.Level {
	switch level {
	case log2.TraceLevel:
		return zapcore.DebugLevel
	case log2.DebugLevel:
		return zapcore.DebugLevel
	case log2.InfoLevel:
		return zapcore.InfoLevel
	case log2.WarnLevel:
		return zapcore.WarnLevel
	case log2.ErrorLevel:
		return zapcore.ErrorLevel
	case log2.FatalLevel:
		return zapcore.FatalLevel
	default:
		// For when new log level is not added to this func (most likely never).
		panic(fmt.Sprintf("unknown log level: %d", level))
	}
}

// UnzapifyLevel turns zap log level to interface log level.
func UnzapifyLevel(level zapcore.Level) log2.Level {
	switch level {
	case zapcore.DebugLevel:
		return log2.DebugLevel
	case zapcore.InfoLevel:
		return log2.InfoLevel
	case zapcore.WarnLevel:
		return log2.WarnLevel
	case zapcore.ErrorLevel:
		return log2.ErrorLevel
	case zapcore.FatalLevel, zapcore.DPanicLevel, zapcore.PanicLevel:
		return log2.FatalLevel
	default:
		// For when new log level is not added to this func (most likely never).
		panic(fmt.Sprintf("unknown log level: %d", level))
	}
}

// nolint: gocyclo
func zapifyField(field log2.Field) zap.Field {
	switch field.Type() {
	case log2.FieldTypeNil:
		return zap.Reflect(field.Key(), nil)
	case log2.FieldTypeString:
		return zap.String(field.Key(), field.String())
	case log2.FieldTypeBinary:
		return zap.Binary(field.Key(), field.Binary())
	case log2.FieldTypeBoolean:
		return zap.Bool(field.Key(), field.Bool())
	case log2.FieldTypeSigned:
		return zap.Int64(field.Key(), field.Signed())
	case log2.FieldTypeUnsigned:
		return zap.Uint64(field.Key(), field.Unsigned())
	case log2.FieldTypeFloat:
		return zap.Float64(field.Key(), field.Float())
	case log2.FieldTypeTime:
		return zap.Time(field.Key(), field.Time())
	case log2.FieldTypeDuration:
		return zap.Duration(field.Key(), field.Duration())
	case log2.FieldTypeError:
		return zap.NamedError(field.Key(), field.Error())
	case log2.FieldTypeArray:
		return zap.Any(field.Key(), field.Interface())
	case log2.FieldTypeAny:
		return zap.Any(field.Key(), field.Interface())
	case log2.FieldTypeReflect:
		return zap.Reflect(field.Key(), field.Interface())
	case log2.FieldTypeByteString:
		return zap.ByteString(field.Key(), field.Binary())
	default:
		// For when new field type is not added to this func
		panic(fmt.Sprintf("unknown field type: %d", field.Type()))
	}
}

func zapifyFields(fields ...log2.Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zapifyField(field))
	}

	return zapFields
}
