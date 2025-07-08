package testdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

func TestGet(t *testing.T) {
	cfg := &ProjectLogs{}
	get, err := config.Get("projectLogs")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	err = get.Scan(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, "debug", cfg.Level)

	level, err := config.Get("projectLogs.level")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	get1 := level.String("s1")
	assert.Equal(t, "debug", get1)

	level2, err := config.Get("projectLogs.level1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	get2 := level2.String("s1")
	assert.Equal(t, "s1", get2)
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

type Host struct {
	IP string
}

func Test1(t *testing.T) {
	var host Host
	err := config.ScanFrom(&host, "host")
	t.Logf("host ip: %v, err: %v", host.IP, err)
}
