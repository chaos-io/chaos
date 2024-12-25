package logs

import (
	"bytes"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/chaos-io/chaos/config"
)

var defaultLog *SugaredLogger

func init() {
	logCfg := &Config{}
	_ = config.ScanFrom(logCfg, "logs")
	defaultLog = New(logCfg)
}

func Logger() *SugaredLogger {
	return defaultLog
}

func With(args ...interface{}) *SugaredLogger {
	return defaultLog.With(args...)
}

func AddCallerSkip(skip int) *SugaredLogger {
	// add a new caller, and -1
	return Logger().WithOptions(zap.AddCallerSkip(skip - 1))
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(args ...interface{}) {
	Logger().Debug(args...)
}

// Info logs a message at InfoLevel.
func Info(args ...interface{}) {
	Logger().Info(args...)
}

// Warn logs a message at WarnLevel.
func Warn(args ...interface{}) {
	Logger().Warn(args...)
}

// Error logs a message at ErrorLevel.
func Error(args ...interface{}) {
	Logger().Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	Logger().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	Logger().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	Logger().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	Logger().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	Logger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Logger().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Logger().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Logger().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Logger().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Logger().Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	Logger().Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func NewErrorf(template string, args ...interface{}) error {
	Logger().Errorf(template, args...)
	return fmt.Errorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	Logger().Errorw(msg, keysAndValues...)

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
