package log

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/chaos-io/chaos/core/log/zap"
	"go.uber.org/zap/zapcore"
)

func LevelEnabled(level Level) bool {
	return zap.DefaultLevel.Enabled(zapcore.Level(level))
}

func SWith(args ...interface{}) *zap.SugaredLogger {
	return zap.DefaultLog.With(args...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(args ...interface{}) {
	zap.ZapLogger().Debug(args...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(args ...interface{}) {
	zap.ZapLogger().Info(args...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(args ...interface{}) {
	zap.ZapLogger().Warn(args...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(args ...interface{}) {
	zap.ZapLogger().Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	zap.ZapLogger().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	zap.ZapLogger().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	zap.ZapLogger().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	zap.ZapLogger().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	zap.ZapLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	zap.ZapLogger().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	zap.ZapLogger().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zap.ZapLogger().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zap.ZapLogger().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zap.ZapLogger().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	zap.ZapLogger().Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	zap.ZapLogger().Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func NewErrorf(template string, args ...interface{}) error {
	zap.ZapLogger().Errorf(template, args...)
	return fmt.Errorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	zap.ZapLogger().Errorw(msg, keysAndValues...)

	buffer := bytes.NewBufferString(msg)
	buffer.WriteString(" ")
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf("%v: %v", keysAndValues[i], keysAndValues[i+1]))
	}

	return errors.New(buffer.String())
}
