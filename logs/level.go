package logs

import (
	"go.uber.org/zap/zapcore"

	"github.com/chaos-io/chaos/config/reader"
)

type Level = zapcore.Level

const (
	LevelPath = "logs.level"
)

func ChangeLogLevel(value reader.Value) {
	level := value.String("")
	if len(level) > 0 && (level == "debug" || level == "info" || level == "error" ||
		level == "warn" || level == "panic" || level == "fatal") {
		SetLevel(level)
	}
}
