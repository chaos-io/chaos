package logs

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/chaos-io/chaos/logs/testdata"
)

func TestDebugw(t *testing.T) {
	Debugw("debugw", "map", map[string]interface{}{"1": 1, "2": "two"})
	Debugw("debugw", "slice", []string{"1", "2"})
	i := new(int)
	*i = 8
	Debugw("debugw", "ptr", i)
	Debugw("debugw", "addr", &i)
}

func TestDebugf(t *testing.T) {
	Debugf("the debugf, string=%s", "aaa")
}

func TestNewErrorf(t *testing.T) {
	err := NewErrorf("the newErrorf, string=%s", "aaa")
	Debugw("the debugw", "debugw error", err)
}

func TestNewErrorw(t *testing.T) {
	err := NewErrorw("the newErrorw", "err", "this is a error")
	Debugw("the debugw", "debugw error", err)
}

// priority: debug < info < warn < error < DPanic < panic < fatal
func TestLevelLogs(t *testing.T) {
	print := func(level string) {
		t.Logf("%s-------------------------------------------------\n", level)
		SetLevel(level)
		Debug("debug")
		Info("info")
		Warn("warn")
		Error("error")
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
	l := New(&Config{
		Output: "file",
		File: FileConfig{
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
	stat := &testdata.CpuStat{
		Number: 0,
		State:  "123",
	}
	Infow("log infow", "stat", stat)

	// already skip 2 layer caller
	logger := Logger().With()
	logger.WithOptions(zap.AddCallerSkip(-2)).Infow("skip -2")
	logger.WithOptions(zap.AddCallerSkip(-1)).Infow("skip -1")
	logger.WithOptions(zap.AddCallerSkip(0)).Infow("skip 0")
	logger.WithOptions(zap.AddCallerSkip(1)).Infow("skip 1")
	logger.WithOptions(zap.AddCallerSkip(2)).Infow("skip 2")
}
