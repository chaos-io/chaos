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
