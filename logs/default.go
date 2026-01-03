package logs

import (
	"github.com/chaos-io/core/go/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/chaos-io/chaos/config"
)

type defaultLogger struct {
	log *zap.SugaredLogger
}

func newDefaultLogger() Logger {
	cfg := &logs.Config{}
	_ = config.ScanFrom(cfg, "logs")

	log := logs.NewSugaredLogger(cfg).WithOptions(zap.AddCallerSkip(1))
	return &defaultLogger{log: log}
}

func NewLoggerWith(cfg *logs.Config) Logger {
	log := logs.NewSugaredLogger(cfg).WithOptions(zap.AddCallerSkip(1))
	return &defaultLogger{log: log}
}

func (l *defaultLogger) GetLevel() Level {
	return Level(l.log.Level())
}

func (l *defaultLogger) SetLevel(level Level) {
	lev := zapcore.Level(level).String()
	logs.SetLevel(lev)
}

func (l *defaultLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *defaultLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *defaultLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *defaultLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *defaultLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *defaultLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *defaultLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.log.Debugw(msg, keysAndValues...)
}

func (l *defaultLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.log.Infow(msg, keysAndValues...)
}

func (l *defaultLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.log.Warnw(msg, keysAndValues...)
}

func (l *defaultLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.log.Errorw(msg, keysAndValues...)
}

func (l *defaultLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.log.Fatalw(msg, keysAndValues...)
}
