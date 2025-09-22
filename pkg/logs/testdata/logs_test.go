package testdata

import (
	"os"
	"strings"
	"testing"

	"github.com/chaos-io/chaos/pkg/logs"
	"github.com/stretchr/testify/assert"
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

/* caller
zap@v1.27.0/logger.go:222 -3
zap@v1.27.0/sugar.go:354  -2
zap@v1.27.0/sugar.go:198  -1
logs/default.go:83         0 defaultLogger.Debugf
logs/logs.go:56            1 Debugf
testdata/logs_test.go:22   2 test
testing/testing.go:1792    3
*/

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
	printFunc := func(level string) {
		t.Logf("%s-------------------------------------------------\n", level)
		logs.SetLogLevel(logs.LevelConv(level))
		logs.Debug("debug")
		logs.Info("info")
		logs.Warn("warn")
		logs.Error("error")
		// Fatal("fatal")
		t.Log("-------------------------------------------------")
	}
	printFunc("debug")
	printFunc("info")
	printFunc("warn")
	printFunc("error")
	// printFunc("fatal")
}

func TestLogFileJSON(t *testing.T) {
	logs.Debug("begin")
	defer logs.Debug("end")

	filename := "./app.log"
	l := logs.NewSugaredLogger(&logs.Config{
		Output: "file",
		File: logs.FileConfig{
			Path: filename,
		},
	})

	value := "foo"
	values := []string{"foo", "bar", "baz"}
	mapVals := map[string]any{"foo": true, "bar": 100}
	l.Infow("info", "value", value, "values", values, "map", mapVals)

	logCont := `{"value": "foo", "values": ["foo", "bar", "baz"], "map": {"foo":true,"bar":100}}`
	readFile, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(readFile), logCont))
	_ = os.Remove(filename)
}

type MyStruct struct {
	l logs.Logger
}

func TestLogger(t *testing.T) {
	logs.Debugw("begin", "1", "2")
	my := &MyStruct{}
	my.l = logs.DefaultLogger()
	// TODO: caller has an extra layer
	my.l.Debugw("my.debug", "1", "2")
	logs.Debugw("end", "1", "2")
}
