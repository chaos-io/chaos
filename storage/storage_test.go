package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	chaosconfig "github.com/chaos-io/chaos/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	Register("stub", func(cfg *Config) (Storage, error) {
		return &stubStorage{bucket: cfg.BucketName}, nil
	})
	Register("stub-config", func(cfg *Config) (Storage, error) {
		return &stubStorage{bucket: cfg.BucketName}, nil
	})

	t.Run("new with config", func(t *testing.T) {
		st, err := NewWithConfig(&Config{
			Vendor:     "stub",
			BucketName: "chaos",
		})
		require.NoError(t, err)
		require.IsType(t, &stubStorage{}, st)
	})

	t.Run("rejects unknown vendor", func(t *testing.T) {
		st, err := NewWithConfig(&Config{
			Vendor:     "missing",
			BucketName: "chaos",
		})
		require.ErrorIs(t, err, ErrUnsupportedVendor)
		require.Nil(t, st)
	})

	t.Run("loads config", func(t *testing.T) {
		loadTestConfig(t, "storage.yaml", "storage:\n  vendor: stub-config\n  bucketName: chaos\n")

		st, err := New()
		require.NoError(t, err)
		require.IsType(t, &stubStorage{}, st)
	})
}

type stubStorage struct {
	bucket string
}

func (s *stubStorage) Read(context.Context, string, ...Option) (*Object, error) {
	return nil, nil
}

func (s *stubStorage) Write(context.Context, *Object, ...Option) error {
	return nil
}

func (s *stubStorage) Download(context.Context, string, string, ...Option) error {
	return nil
}

func (s *stubStorage) Upload(context.Context, string, string, ...Option) error {
	return nil
}

func (s *stubStorage) PresignedDownloadURL(context.Context, string, ...Option) (string, error) {
	return "", nil
}

func (s *stubStorage) PresignedUploadURL(context.Context, string, ...Option) (string, error) {
	return "", nil
}

func loadTestConfig(t *testing.T, filename, body string) {
	t.Helper()

	if err := chaosconfig.InitDefault(chaosconfig.WithWatcherDisabled()); err != nil {
		t.Fatalf("InitDefault() failed: %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "config", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if err := chaosconfig.LoadPath(filepath.Join(dir, "config")); err != nil {
		t.Fatalf("LoadPath() failed: %v", err)
	}
}
