package logs

import (
	"strings"

	"go.uber.org/zap"
)

var defaultLevel zap.AtomicLevel

func SetLevel(level string) {
	l := strings.ToLower(level)
	switch l {
	case "info":
		defaultLevel.SetLevel(zap.InfoLevel)
	case "debug":
		defaultLevel.SetLevel(zap.DebugLevel)
	case "warn":
		defaultLevel.SetLevel(zap.WarnLevel)
	case "error":
		defaultLevel.SetLevel(zap.ErrorLevel)
	case "panic":
		defaultLevel.SetLevel(zap.PanicLevel)
	case "fatal":
		defaultLevel.SetLevel(zap.FatalLevel)
	default:
		defaultLevel.SetLevel(zap.InfoLevel)
	}
}
