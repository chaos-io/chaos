package testdata

import (
	"os"
	"strings"
	"testing"

	"github.com/chaos-io/chaos/pkg/logs"
	logs2 "github.com/chaos-io/core/go/logs"
	"github.com/stretchr/testify/assert"
)

func TestDebugw(t *testing.T) {
	logs.Debugw("debugw", "map", map[string]interface{}{"1": 1, "2": "two"})
	logs.Debugw("debugw", "slice", []string{"1", "2"})
	i := new(int)
	*i = 8
	logs.Debugw("debugw", "ptr", i)
	logs.Debugw("debugw", "addr", &i)

	logs.Debugf("the debugf, string=%s", "aaa")

	err := logs.NewErrorf("the newErrorf, string=%s", "aaa")
	logs.Debugw("the debugw", "debugw error", err)

	err = logs.NewErrorw("the newErrorw", "err", "this is a error")
	logs.Debugw("the debugw", "debugw error", err)
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
// priority: debug < info < warn < error < DPanic < panic < fatal
func TestLevelLogs(t *testing.T) {
	printFunc := func(level logs.Level) {
		t.Logf("%v-------------------------------------------------\n", level)
		logs.SetLogLevel(level)
		logs.Debug("debug")
		logs.Info("info")
		logs.Warn("warn")
		logs.Error("error")
		t.Log("-------------------------------------------------")
	}
	printFunc(logs.DebugLevel)
	printFunc(logs.InfoLevel)
	printFunc(logs.WarnLevel)
	printFunc(logs.ErrorLevel)
	// printFunc(logs.FatalLevel)
}

func TestLogFileJSON(t *testing.T) {
	logs.Debug("begin")
	defer logs.Debug("end")

	filename := "./app.log"
	l := logs.NewLoggerWith(&logs2.Config{
		Output: "file",
		File: logs2.FileConfig{
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
