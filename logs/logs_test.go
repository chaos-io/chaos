package logs

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chaos-io/chaos/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerLevelAndFields(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{
		Level:  "warn",
		Encode: "json",
		Output: "console",
	}, buf))

	svc.Debug("debug")
	svc.Info("info")
	svc.Warnw("warn", "user", "alice")
	svc.Errorw("error", "count", 2)

	out := buf.String()
	require.NotContains(t, out, `"message":"debug"`)
	require.NotContains(t, out, `"message":"info"`)
	require.Contains(t, out, `"message":"warn"`)
	require.Contains(t, out, `"user":"alice"`)
	require.Contains(t, out, `"message":"error"`)
	require.Contains(t, out, `"count":2`)
}

func TestServiceHelpers(t *testing.T) {
	buf := &bytes.Buffer{}
	svc := NewService(newZapLogger(&Config{
		Level:  "info",
		Encode: "json",
		Output: "console",
	}, buf))

	require.Equal(t, InfoLevel, svc.Logger().GetLevel())
	svc.SetLogLevel(DebugLevel)
	require.Equal(t, DebugLevel, svc.Logger().GetLevel())

	err := svc.NewErrorw("failed", "kind", "network")
	require.EqualError(t, err, "failed kind: network")

	svc.Debugf("value=%s", "x")
	svc.Infow("ready", "id", 1)

	out := buf.String()
	require.Contains(t, out, `"message":"failed"`)
	require.Contains(t, out, `"kind":"network"`)
	require.Contains(t, out, `"message":"value=x"`)
	require.Contains(t, out, `"message":"ready"`)
}

func TestPackageDefaultLoggerSwap(t *testing.T) {
	prev := DefaultLogger()
	t.Cleanup(func() {
		SetLogger(prev)
	})

	buf := &bytes.Buffer{}
	SetLogger(newZapLogger(&Config{
		Level:  "info",
		Encode: "json",
		Output: "console",
	}, buf))

	Infow("hello", "id", 7)
	require.Contains(t, buf.String(), `"message":"hello"`)
	require.Contains(t, buf.String(), `"id":7`)
}

func TestDefaultServiceLoadsLogsConfig(t *testing.T) {
	prev := DefaultLogger()
	t.Cleanup(func() {
		SetLogger(prev)
	})

	require.NoError(t, config.InitDefault(config.WithWatcherDisabled()))

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "logs.yaml"), []byte(`logs:
  level: debug
  encode: console
  output: console
`), 0o644))
	require.NoError(t, config.LoadPath(dir))

	require.NoError(t, ReloadDefaultServiceFromConfig())
	logger := DefaultLogger()
	require.Equal(t, DebugLevel, logger.GetLevel())
}

func TestLevelHandlerServeHTTP(t *testing.T) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	handler := levelHandler(&level, "/level")

	req := httptest.NewRequest(http.MethodGet, "/level", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"level":"info"`)

	req = httptest.NewRequest(http.MethodPut, "/level", strings.NewReader(`{"level":"debug"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, zapcore.DebugLevel, level.Level())
	require.Contains(t, rec.Body.String(), `"level":"debug"`)
}

func TestStartLevelServer(t *testing.T) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	svc := startLevelServerWithListener(&level, "/level", ln)
	t.Cleanup(func() {
		_ = svc.Close()
	})

	client := &http.Client{Timeout: 2 * time.Second}
	url := "http://" + ln.Addr().String() + "/level"

	require.Eventually(t, func() bool {
		resp, err := client.Get(url)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		return resp.StatusCode == http.StatusOK && strings.Contains(string(body), `"level":"info"`)
	}, time.Second, 20*time.Millisecond)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(`{"level":"debug"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), `"level":"debug"`)
	require.Equal(t, zapcore.DebugLevel, level.Level())
}

func TestFileOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	logger := NewLoggerWith(&Config{
		Level:  "info",
		Encode: "json",
		Output: "file",
		File: FileConfig{
			Path:   path,
			Encode: "json",
			// Encode:     "console",
			MaxSize:    1,
			MaxBackups: 1,
			MaxAge:     1,
		},
	})

	logger.Log(Entry{
		Level:   InfoLevel,
		Message: "persisted",
		Fields:  []Field{{Key: "user", Value: "bob"}},
	})

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(data), `"message":"persisted"`))
	require.True(t, strings.Contains(string(data), `"user":"bob"`))
}
