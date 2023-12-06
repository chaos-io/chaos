package logs

import (
	"testing"
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
