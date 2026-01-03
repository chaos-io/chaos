package logs

import (
	"bytes"
	"errors"
	"fmt"
)

var logger Logger = newDefaultLogger()

func DefaultLogger() Logger {
	return logger
}

// SetLogger sets the logger.
// Note that this method is not concurrent-safe.
func SetLogger(l Logger) {
	logger = l
}

// SetLogLevel sets the level of logs below which logs will not be output.
// The default log level is LevelInfo.
// Note that this method is not concurrent-safe.
func SetLogLevel(level Level) {
	logger.SetLevel(level)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info logs a message at InfoLevel.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn logs a message at WarnLevel.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error logs a message at ErrorLevel.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	logger.Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func NewErrorf(template string, args ...interface{}) error {
	logger.Errorf(template, args...)
	return fmt.Errorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	logger.Errorw(msg, keysAndValues...)

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
