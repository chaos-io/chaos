package testdata

import (
	"testing"
	"time"

	"github.com/chaos-io/chaos/config"
	"github.com/chaos-io/chaos/config/reader"
)

type ProjectLogs struct {
	Level string `json:"level" default:"info"`
}

func TestScanFrom(t *testing.T) {
	cfg := &ProjectLogs{}
	if err := config.ScanFrom(cfg, "projectLogs"); err != nil {
		t.Errorf("ScanFrom() error = %v", err)
	}
	t.Logf("got level %v", cfg.Level)
}

func TestHotLoad(t *testing.T) {
	cfg := &ProjectLogs{}
	if err := config.ScanFrom(cfg, "projectLogs"); err != nil {
		t.Errorf("ScanFrom() error = %v", err)
	}

	t.Logf("got first level %v", cfg.Level)

	watcherObj, err := config.WatchFunc(func(v reader.Value) {
		if err := v.Scan(&cfg); err != nil {
			t.Logf("scan error: %v", err)
		}
	}, "projectLogs")
	if err != nil {
		t.Logf("watch error: %v", err)
	} else {
		defer func() { watcherObj.Close() }()
	}

	time.Sleep(5 * time.Second)
	t.Logf("got level %v", cfg.Level)
}
