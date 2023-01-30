package logger

import (
	"github.com/chaos-io/chaos/core/log"
	"github.com/chaos-io/chaos/core/log/nop"
	"github.com/chaos-io/chaos/core/log/zap"
)

var Log log.Logger = new(nop.Logger)

func Setup(lvl log.Level) {
	Log = zap.Must(zap.ConsoleConfig(lvl))
}
