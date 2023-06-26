package logs

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestDebugw(t *testing.T) {
	Debugw("debugw", "map", map[string]interface{}{"1": 1, "2": "two"})
	Debugw("debugw", "slice", []string{"1", "2"})
	i := new(int)
	*i = 8
	Debugw("debugw", "ptr", i)
	Debugw("debugw", "addr", &i)

}

var url = "http://192.168.1.1:8800"

func Test_SugaredLogger_Example(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

// Logger only supports structured logging.
func Test_Logger_Example(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func TestNewErrorw(t *testing.T) {
	type args struct {
		msg           string
		keysAndValues []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{
			msg:           "arg1",
			keysAndValues: []interface{}{"123", 1234, "2123", []string{"1", "2"}},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewErrorw(tt.args.msg, tt.args.keysAndValues...); (err != nil) != tt.wantErr {
				t.Errorf("NewErrorw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
