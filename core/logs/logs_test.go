package logs

import (
	"testing"
)

func TestDebugw(t *testing.T) {
	Debugw("debugw", "map",
		map[string]interface{}{
			"1": 1,
			"2": "two",
		})
}
