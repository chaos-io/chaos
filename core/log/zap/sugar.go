package zap

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/chaos-io/chaos/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SugaredLogger = zap.SugaredLogger

var DefaultLog *SugaredLogger
var DefaultLevel zap.AtomicLevel
var eLog *SugaredLogger

func init() {
	logCfg := &Config{}
	config.Get("log").Scan(logCfg)
	DefaultLog = NewZap(logCfg)

	elogCfg := &Config{}
	config.Get("elog").Scan(elogCfg)
	eLog = NewZap(elogCfg)
}

func ZapLogger() *SugaredLogger {
	return DefaultLog
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

	DefaultLevel = zap.NewAtomicLevel()
	SetLevel(cfg.Level)

	if len(cfg.LevelPattern) > 0 && cfg.LevelPort > 0 {
		http.HandleFunc(cfg.LevelPattern, DefaultLevel.ServeHTTP)
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
		DefaultLevel.SetLevel(zap.InfoLevel)
	} else if l == "debug" {
		DefaultLevel.SetLevel(zap.DebugLevel)
	} else if l == "error" {
		DefaultLevel.SetLevel(zap.ErrorLevel)
	} else if l == "warn" {
		DefaultLevel.SetLevel(zap.WarnLevel)
	} else if l == "panic" {
		DefaultLevel.SetLevel(zap.PanicLevel)
	} else if l == "fatal" {
		DefaultLevel.SetLevel(zap.FatalLevel)
	} else {
		DefaultLevel.SetLevel(zap.InfoLevel)
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
	return zapcore.NewCore(formatEncoder, consoleDebugging, DefaultLevel)
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

	return zapcore.NewCore(formatEncoder, w, DefaultLevel)
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
