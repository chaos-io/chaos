package golog

import (
	"github.com/chaos-io/chaos/core/log"
	canal_log "github.com/siddontang/go-log/log"
)

func SetLevel(level log.Level) {
	switch level {
	case log.DebugLevel:
		canal_log.SetLevel(canal_log.LevelDebug)
	case log.ErrorLevel:
		canal_log.SetLevel(canal_log.LevelError)
	case log.FatalLevel:
		canal_log.SetLevel(canal_log.LevelFatal)
	case log.InfoLevel:
		canal_log.SetLevel(canal_log.LevelInfo)
	case log.TraceLevel:
		canal_log.SetLevel(canal_log.LevelTrace)
	}
}
