package testdata

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/chaos-io/chaos/logs"
)

func TestDebugw(t *testing.T) {
	logs.Debugw("debugw", "map", map[string]interface{}{"1": 1, "2": "two"})
	logs.Debugw("debugw", "slice", []string{"1", "2"})
	i := new(int)
	*i = 8
	logs.Debugw("debugw", "ptr", i)
	logs.Debugw("debugw", "addr", &i)
}

func TestDebugf(t *testing.T) {
	logs.Debugf("the debugf, string=%s", "aaa")
}

func TestNewErrorf(t *testing.T) {
	err := logs.NewErrorf("the newErrorf, string=%s", "aaa")
	logs.Debugw("the debugw", "debugw error", err)
}

func TestNewErrorw(t *testing.T) {
	err := logs.NewErrorw("the newErrorw", "err", "this is a error")
	logs.Debugw("the debugw", "debugw error", err)
}

// priority: debug < info < warn < error < DPanic < panic < fatal
func TestLevelLogs(t *testing.T) {
	print := func(level string) {
		t.Logf("%s-------------------------------------------------\n", level)
		logs.SetLevel(level)
		logs.Debug("debug")
		logs.Info("info")
		logs.Warn("warn")
		logs.Error("error")
		// Fatal("fatal")
		t.Log("-------------------------------------------------")
	}
	print("debug")
	print("info")
	print("warn")
	print("error")
	// print("fatal")
}

func TestLogFileJSON(t *testing.T) {
	const logFileName = "./test.log"
	l := logs.New(&logs.Config{
		Output: "file",
		File: logs.FileConfig{
			Path: logFileName,
		},
	})

	value := "foo"
	values := []string{"foo", "bar", "baz"}
	mapVals := map[string]any{"foo": true, "bar": 100}
	l.Infow("info", "value", value, "values", values, "map", mapVals)

	content, err := os.ReadFile(logFileName)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.True(t, strings.Contains(string(content), "{\"value\": \"foo\", \"values\": [\"foo\",\"bar\",\"baz\"], \"map\": {\"foo\":true,\"bar\":100}}"))

	_ = os.Remove(logFileName)
}

func TestConsoleJson(t *testing.T) {
	stat := &CpuStat{
		Number: 0,
		State:  "123",
	}
	logs.Infow("log infow", "stat", stat)
	logs.Logger().AddCallerSkip(-1).Infow("log infow-1", "stat", stat)
	logs.Logger().AddCallerSkip(0).Infow("log infow0", "stat", stat)
	logs.Logger().AddCallerSkip(1).Infow("log infow1", "stat", stat)

	logger := logs.Logger()
	logger.WithOptions(zap.AddCallerSkip(-2)).Infow("skip -2")
	logger.WithOptions(zap.AddCallerSkip(-1)).Infow("skip -1")
	logger.WithOptions(zap.AddCallerSkip(0)).Infow("skip 0")
	logger.WithOptions(zap.AddCallerSkip(1)).Infow("skip 1")
	logger.WithOptions(zap.AddCallerSkip(2)).Infow("skip 2")
}
