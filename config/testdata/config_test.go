package testdata

import (
	"testing"
	"time"

	"github.com/chaos-io/chaos/config"
)

type ProjectLog struct {
	Config *ProjectLogConfig
}

type ProjectLogConfig struct {
	Level string `json:"level" default:"info"`
}

func TestScanFrom(t *testing.T) {
	s := &ProjectLog{}
	cfg := &ProjectLogConfig{}
	if err := config.ScanFrom(cfg, "projectLogs"); err != nil {
		t.Errorf("ScanFrom() error = %v", err)
	}
	s.Config = cfg
	t.Logf("got level %v", s.Config.Level)
}

func TestHotLoad(t *testing.T) {
	cfg := &ProjectLogConfig{}
	if err := config.ScanFrom(cfg, "projectLogs"); err != nil {
		t.Errorf("ScanFrom() error = %v", err)
	}

	t.Logf("got first level %v", cfg.Level)

	ticker := time.NewTicker(time.Second)
	timer := time.NewTimer(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			t.Logf("got level %v", cfg.Level)
		case <-timer.C:
			return
		}
	}
}
