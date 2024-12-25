package logs

import (
	"strings"

	"go.uber.org/zap"
)

var defaultLevel zap.AtomicLevel

func SetLevel(level string) {
	l := strings.ToLower(level)
	if l == "info" {
		defaultLevel.SetLevel(zap.InfoLevel)
	} else if l == "debug" {
		defaultLevel.SetLevel(zap.DebugLevel)
	} else if l == "error" {
		defaultLevel.SetLevel(zap.ErrorLevel)
	} else if l == "warn" {
		defaultLevel.SetLevel(zap.WarnLevel)
	} else if l == "panic" {
		defaultLevel.SetLevel(zap.PanicLevel)
	} else if l == "fatal" {
		defaultLevel.SetLevel(zap.FatalLevel)
	} else {
		defaultLevel.SetLevel(zap.InfoLevel)
	}
}
