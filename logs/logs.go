package logs

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/chaos-io/chaos/config"
)

type Service struct {
	logger Logger
}

func NewService(logger Logger) *Service {
	if logger == nil {
		logger = newNopLogger()
	}
	return &Service{logger: logger}
}

func newDefaultService() *Service {
	return NewService(NewLoggerWith(NewDefaultConfig()))
}

var (
	defaultServiceMu sync.RWMutex
	defaultService   = newDefaultService()
)

// ReloadDefaultServiceFromConfig re-applies the package logger settings from the
// current default config. Call this after config.InitDefault and Load*.
func ReloadDefaultServiceFromConfig() error {
	if !config.IsDefaultInitialized() {
		return nil
	}

	cfg := NewDefaultConfig()
	if err := config.ScanFrom(cfg, "logs"); err != nil {
		return err
	}

	defaultServiceMu.Lock()
	defaultService = NewService(NewLoggerWith(cfg))
	defaultServiceMu.Unlock()

	return nil
}

func DefaultLogger() Logger {
	defaultServiceMu.RLock()
	defer defaultServiceMu.RUnlock()
	return defaultService.logger
}

func SetLogger(logger Logger) {
	defaultServiceMu.Lock()
	defaultService = NewService(logger)
	defaultServiceMu.Unlock()
}

func SetLogLevel(level Level) {
	defaultServiceMu.RLock()
	svc := defaultService
	defaultServiceMu.RUnlock()
	svc.SetLogLevel(level)
}

func (s *Service) Logger() Logger {
	if s == nil {
		return newNopLogger()
	}
	return s.logger
}

func (s *Service) SetLogger(logger Logger) {
	if s == nil {
		return
	}
	if logger == nil {
		logger = newNopLogger()
	}
	s.logger = logger
}

func (s *Service) SetLogLevel(level Level) {
	if s == nil || s.logger == nil {
		return
	}
	s.logger.SetLevel(level)
}

func (s *Service) log(level Level, msg string, fields ...Field) {
	if s == nil || s.logger == nil {
		return
	}
	s.logger.Log(Entry{Level: level, Message: msg, Fields: fields})
}

func (s *Service) Debug(args ...interface{}) {
	s.log(DebugLevel, fmt.Sprint(args...))
}

func (s *Service) Info(args ...interface{}) {
	s.log(InfoLevel, fmt.Sprint(args...))
}

func (s *Service) Warn(args ...interface{}) {
	s.log(WarnLevel, fmt.Sprint(args...))
}

func (s *Service) Error(args ...interface{}) {
	s.log(ErrorLevel, fmt.Sprint(args...))
}

func (s *Service) Fatal(args ...interface{}) {
	s.log(FatalLevel, fmt.Sprint(args...))
}

func (s *Service) Debugf(template string, args ...interface{}) {
	s.log(DebugLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Infof(template string, args ...interface{}) {
	s.log(InfoLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Warnf(template string, args ...interface{}) {
	s.log(WarnLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Errorf(template string, args ...interface{}) {
	s.log(ErrorLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Fatalf(template string, args ...interface{}) {
	s.log(FatalLevel, fmt.Sprintf(template, args...))
}

func (s *Service) Debugw(msg string, keysAndValues ...interface{}) {
	s.log(DebugLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Infow(msg string, keysAndValues ...interface{}) {
	s.log(InfoLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Warnw(msg string, keysAndValues ...interface{}) {
	s.log(WarnLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Errorw(msg string, keysAndValues ...interface{}) {
	s.log(ErrorLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) Fatalw(msg string, keysAndValues ...interface{}) {
	s.log(FatalLevel, msg, keyValuesToFields(keysAndValues...)...)
}

func (s *Service) NewError(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	s.Error(msg)
	return errors.New(msg)
}

func (s *Service) NewErrorf(template string, args ...interface{}) error {
	msg := fmt.Sprintf(template, args...)
	s.Error(msg)
	return fmt.Errorf(template, args...)
}

func (s *Service) NewErrorw(msg string, keysAndValues ...interface{}) error {
	fields := keyValuesToFields(keysAndValues...)
	s.Errorw(msg, keysAndValues...)
	return errors.New(renderErrorMessage(msg, fields))
}

func Debug(args ...interface{}) {
	defaultService.Debug(args...)
}

func Info(args ...interface{}) {
	defaultService.Info(args...)
}

func Warn(args ...interface{}) {
	defaultService.Warn(args...)
}

func Error(args ...interface{}) {
	defaultService.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultService.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	defaultService.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	defaultService.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultService.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultService.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultService.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultService.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultService.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultService.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultService.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultService.Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	return defaultService.NewError(args...)
}

func NewErrorf(template string, args ...interface{}) error {
	return defaultService.NewErrorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	return defaultService.NewErrorw(msg, keysAndValues...)
}

func keyValuesToFields(keysAndValues ...interface{}) []Field {
	if len(keysAndValues) < 2 {
		return nil
	}

	fields := make([]Field, 0, len(keysAndValues)/2)
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		fields = append(fields, Field{
			Key:   fmt.Sprint(keysAndValues[i]),
			Value: keysAndValues[i+1],
		})
	}
	return fields
}

func renderErrorMessage(msg string, fields []Field) string {
	if len(fields) == 0 {
		return msg
	}

	var b bytes.Buffer
	b.WriteString(msg)
	b.WriteByte(' ')
	for i, field := range fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%v: %v", field.Key, field.Value))
	}
	return b.String()
}

type nopLogger struct{}

func newNopLogger() Logger { return nopLogger{} }

func (nopLogger) SetLevel(Level)       {}
func (nopLogger) GetLevel() Level      { return InfoLevel }
func (nopLogger) With(...Field) Logger { return nopLogger{} }
func (nopLogger) Log(Entry)            {}
