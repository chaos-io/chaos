package logs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewSugaredLogger level string, encode string, port int, pattern string, initFields map[string]interface{}
func NewSugaredLogger(cfg *Config) *zap.SugaredLogger {
	sl := &zap.SugaredLogger{}
	defaultLevel = zap.NewAtomicLevel()

	if len(cfg.LevelPattern) > 0 && cfg.LevelPort > 0 {
		http.HandleFunc(cfg.LevelPattern, defaultLevel.ServeHTTP)
		go func() {
			fmt.Printf("level serve on port:%d\nusage: [GET] curl http://localhost:%d%s\nusage: [PUT] curl -XPUT --data '{\"level\":\"debug\"}' http://localhost:%d%s\n", cfg.LevelPort, cfg.LevelPort, cfg.LevelPattern, cfg.LevelPort, cfg.LevelPattern)
			svc := http.Server{
				Addr:         fmt.Sprintf(":%d", cfg.LevelPort),
				Handler:      nil,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
			}
			if err := svc.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}

	cores := make([]zapcore.Core, 0)
	if strings.EqualFold(cfg.Output, "console") {
		cores = append(cores, initWithConsole(cfg.Encode))
	} else if strings.EqualFold(cfg.Output, "file") {
		cores = append(cores, initWithFile(cfg.File))
	}

	core := zapcore.NewTee(cores...)
	coreLogger := zap.New(core)

	initFields := cfg.InitFields
	if len(initFields) > 0 {
		initFieldList := make([]zap.Field, 0, len(initFields))
		for k, v := range initFields {
			var field zap.Field
			if _, ok := v.(proto.Message); ok {
				field = zap.Field{Key: k, Type: zapcore.ReflectType, Interface: v}
			} else {
				field = zap.Any(k, v)
			}

			initFieldList = append(initFieldList, field)
		}

		coreLogger = coreLogger.With(initFieldList...)
	}

	var opts []zap.Option
	opts = append(opts, zap.Development())
	opts = append(opts, zap.AddCaller())
	// opts = append(opts, zap.AddCallerSkip(1))
	opts = append(opts, zap.AddStacktrace(zap.FatalLevel))
	sl = coreLogger.WithOptions(opts...).Sugar()

	return sl
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
		TimeKey:             "time",
		LevelKey:            "level",
		NameKey:             "logger",
		CallerKey:           "caller",
		MessageKey:          "message",
		StacktraceKey:       "stacktrace",
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.LowercaseLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: jsoniterReflectedEncoder,
	}

	var formatEncoder zapcore.Encoder
	if strings.EqualFold(encode, "json") {
		formatEncoder = zapcore.NewJSONEncoder(encodeConfig)
	} else {
		formatEncoder = zapcore.NewConsoleEncoder(encodeConfig)
	}

	return formatEncoder
}

func jsoniterReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := jsoniter.NewEncoder(w)
	// For consistency with custom JSON encoder.
	enc.SetEscapeHTML(false)
	return enc
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
