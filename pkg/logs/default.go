package logs

import (
	"github.com/chaos-io/chaos/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type defaultLogger struct {
	log *SugaredLogger
}

func newDefaultLogger() Logger {
	cfg := &Config{}
	_ = config.ScanFrom(cfg, "logs")

	log := New(cfg).WithOptions(zap.AddCallerSkip(1))
	l := &defaultLogger{log: log}
	l.SetLevel(LevelConv(cfg.Level))
	return l
}

var defaultLevel zap.AtomicLevel

func (l *defaultLogger) GetLevel() LogLevel {
	switch l.log.Level() {
	case zap.DebugLevel:
		return DebugLevel
	case zap.InfoLevel:
		return InfoLevel
	case zap.WarnLevel:
		return WarnLevel
	case zap.ErrorLevel:
		return ErrorLevel
	case zap.FatalLevel:
		return FatalLevel
	default:
		return InfoLevel
	}
}

func (l *defaultLogger) SetLevel(level LogLevel) {
	switch level {
	case DebugLevel:
		defaultLevel.SetLevel(zapcore.DebugLevel)
	case InfoLevel:
		defaultLevel.SetLevel(zapcore.InfoLevel)
	case WarnLevel:
		defaultLevel.SetLevel(zapcore.WarnLevel)
	case ErrorLevel:
		defaultLevel.SetLevel(zapcore.ErrorLevel)
	case FatalLevel:
		defaultLevel.SetLevel(zapcore.FatalLevel)
	default:
		defaultLevel.SetLevel(zapcore.InfoLevel)
	}
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
