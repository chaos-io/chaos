package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/chaos-io/chaos/config/reader"
)

type Level = zapcore.Level

const (
	LevelPath = "logs.level"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zap.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zap.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

func ChangeLogLevel(value reader.Value) {
	level := value.String("")
	if len(level) > 0 && (level == "debug" || level == "info" || level == "error" ||
		level == "warn" || level == "panic" || level == "fatal") {
		SetLevel(level)
	}
}
