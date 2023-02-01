package zap

import (
	"fmt"
	"testing"
)

func TestDebugw(t *testing.T) {
	Info("info", "111")
	fmt.Println("123")
	Debugw("debugw", "this is a map",
		map[string]interface{}{"1": 1, "2":"two"})
	Infow("infow", "this is a map",
		map[string]interface{}{"1": 1, "2":"two"})
}
