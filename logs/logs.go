package logs

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/chaos-io/chaos/config"
)

var (
	defaultLog   *ZapLogger
	defaultLevel zap.AtomicLevel
)

type (
	ZapLogger = zap.SugaredLogger
	Level     = zapcore.Level
)

func init() {
	logCfg := &Config{}
	_ = config.ScanFrom(logCfg, "logs")
	defaultLog = New(logCfg)
}

func Logger() *ZapLogger {
	return defaultLog
}

func New(cfg *Config) *ZapLogger {
	return newZap(cfg)
}

func LevelEnabled(level Level) bool {
	return defaultLevel.Enabled(level)
}

func With(args ...interface{}) *ZapLogger {
	return defaultLog.With(args...)
}

// level string, encode string, port int, pattern string, initFields map[string]interface{}
func newZap(cfg *Config) *ZapLogger {
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
		cores = append(cores, initWithConsole(cfg.Encode))
	} else if strings.Contains(output, "file") {
		cores = append(cores, initWithFile(cfg.File))
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

func initWithConsole(encode string) zapcore.Core {
	formatEncoder := standardEncode(encode)
	consoleDebugging := zapcore.Lock(os.Stdout)
	return zapcore.NewCore(formatEncoder, consoleDebugging, defaultLevel)
}

func initWithFile(fileCfg FileConfig) zapcore.Core {
	formatEncoder := standardEncode(fileCfg.Encode)

	f := handleFileName(fileCfg.Path)
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   f,
		MaxSize:    fileCfg.MaxSize, // megabytes
		MaxAge:     fileCfg.MaxAge,  // days
		MaxBackups: fileCfg.MaxBackups,
		Compress:   fileCfg.Compress,
	})

	return zapcore.NewCore(formatEncoder, w, defaultLevel)
}

func standardEncode(encode string) zapcore.Encoder {
	encodeConfig := zapcore.EncoderConfig{
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

	var formatEncoder zapcore.Encoder
	enc := strings.ToLower(encode)
	if enc == "json" {
		formatEncoder = zapcore.NewJSONEncoder(encodeConfig)
	} else {
		formatEncoder = zapcore.NewConsoleEncoder(encodeConfig)
	}

	return formatEncoder
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
