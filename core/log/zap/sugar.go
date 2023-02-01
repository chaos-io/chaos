package zap

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/core/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SugaredLogger = zap.SugaredLogger

var defaultLog *SugaredLogger
var defaultLevel zap.AtomicLevel
var eLog *SugaredLogger

func init() {
	logCfg := &Config{}
	config.Get("log").Scan(logCfg)
	fmt.Printf("logCfg=%+v\n", logCfg)
	defaultLog = NewZap(logCfg)

	elogCfg := &Config{}
	config.Get("elog").Scan(elogCfg)
	eLog = NewZap(elogCfg)
}

func ZapLogger() *SugaredLogger {
	return defaultLog
}

func ELogger() *SugaredLogger {
	return eLog
}

func ELog(args ...interface{}) {
	ELogger().Info(args...)
}

func ELogf(template string, args ...interface{}) {
	ELogger().Infof(template, args)
}

func ELogw(msg string, keysAndValues ...interface{}) {
	ELogger().Infow(msg, keysAndValues)
}

// NewZap constructs zap-based logger from config
func NewZap(cfg *Config) *SugaredLogger {
	var opts []zap.Option
	opts = append(opts, zap.Development())
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.AddCallerSkip(1))
	opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))

	defaultLevel = zap.NewAtomicLevel()
	SetLevel(cfg.Level)

	if len(cfg.LevelPattern) > 0 && cfg.LevelPort > 0 {
		http.HandleFunc(cfg.LevelPattern, defaultLevel.ServeHTTP)
		go func() {
			fmt.Printf("level serve on port:%d\nusage: [GET] curl http://localhost:%d%s\nusage: [PUT] curl -XPUT --data '{\"level\":\"debug\"}' http://localhost:%d%s\n", cfg.LevelPort, cfg.LevelPort, cfg.LevelPattern, cfg.LevelPort, cfg.LevelPattern)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.LevelPort), nil); err != nil {
				panic(err)
			}
		}()
	}

	cores := make([]zapcore.Core, 0)
	output := strings.ToLower(cfg.Output)
	if strings.Contains(output, "console") {
		cores = append(cores, newConsoleCore(cfg.Encode))
	}
	if strings.Contains(output, "file") {
		cores = append(cores, newFileCore(cfg.File.Encode, cfg.File.Path, cfg.File.MaxSize, cfg.File.MaxBackups, cfg.File.MaxAge))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core)

	initFields := cfg.InitFields
	if len(initFields) > 0 {
		initFieldList := make([]zap.Field, 0)
		for k, v := range initFields {
			initFieldList = append(initFieldList, zap.Any(k, v))
		}
		logger = logger.With(initFieldList...)
	}

	logger = logger.WithOptions(opts...)
	return logger.Sugar()
}

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

// StandardSugaredConfig returns default zap config
func StandardSugaredConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newConsoleCore(encode string) zapcore.Core {
	encodeConfig := StandardSugaredConfig()

	var formatEncoder zapcore.Encoder
	enc := strings.ToLower(encode)
	if enc == "json" {
		formatEncoder = zapcore.NewJSONEncoder(encodeConfig)
	} else {
		formatEncoder = zapcore.NewConsoleEncoder(encodeConfig)
	}
	consoleDebugging := zapcore.Lock(os.Stdout)
	return zapcore.NewCore(formatEncoder, consoleDebugging, defaultLevel)
}

func newFileCore(encode, filename string, maxSize, maxBackups, maxAge int) zapcore.Core {
	encodeConfig := StandardSugaredConfig()

	var formatEncoder zapcore.Encoder
	enc := strings.ToLower(encode)
	if enc == "json" {
		formatEncoder = zapcore.NewJSONEncoder(encodeConfig)
	} else {
		formatEncoder = zapcore.NewConsoleEncoder(encodeConfig)
	}

	f := handleFileName(filename)
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   f,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // days
		Compress:   true,
	})

	return zapcore.NewCore(formatEncoder, w, defaultLevel)
}

func handleFileName(filename string) string {
	filename = path.Clean(filename)
	parts := make([]string, 0)
	var ret string
	paths := strings.Split(filename, string(os.PathSeparator))
	for _, v := range paths {
		val := handleTemplateFileName(v)
		if len(val) > 0 {
			parts = append(parts, val)
		}
	}

	if path.IsAbs(filename) {
		ret = string(os.PathSeparator) + path.Join(parts...)
	} else {
		ret = path.Join(parts...)
	}
	return ret
}

func handleTemplateFileName(template string) string {
	// foo1{hostname}foo2{port}foo3
	lefts := make([]int, 0)
	rights := make([]int, 0)

	size := len(template)
	for i := 0; i < size; i++ {
		if template[i] == '{' {
			lefts = append(lefts, i)
		} else if template[i] == '}' {
			rights = append(rights, i)
		}
	}

	leftSize := len(lefts)
	rightSize := len(rights)
	var minSize int
	if leftSize < rightSize {
		minSize = leftSize
	} else {
		minSize = rightSize
	}

	ret := template
	for i := minSize - 1; i >= 0; i-- {
		variableName := ret[lefts[i]+1 : rights[i]]
		v := os.Getenv(variableName)
		ret = ret[:lefts[i]] + v + ret[rights[i]+1:]
	}
	return ret
}

func LevelEnabled(level log.Level) bool {
	return defaultLevel.Enabled(zapcore.Level(level))
}

func With(args ...interface{}) *zap.SugaredLogger {
	return defaultLog.With(args...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(args ...interface{}) {
	ZapLogger().Debug(args...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(args ...interface{}) {
	ZapLogger().Info(args...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(args ...interface{}) {
	ZapLogger().Warn(args...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(args ...interface{}) {
	ZapLogger().Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	ZapLogger().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	ZapLogger().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	ZapLogger().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	ZapLogger().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	ZapLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	ZapLogger().Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	ZapLogger().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	ZapLogger().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	ZapLogger().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	ZapLogger().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	ZapLogger().Fatalw(msg, keysAndValues...)
}

func NewError(args ...interface{}) error {
	ZapLogger().Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func NewErrorf(template string, args ...interface{}) error {
	ZapLogger().Errorf(template, args...)
	return fmt.Errorf(template, args...)
}

func NewErrorw(msg string, keysAndValues ...interface{}) error {
	ZapLogger().Errorw(msg, keysAndValues...)

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
