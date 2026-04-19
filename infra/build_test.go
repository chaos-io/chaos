package infra

import (
	"context"
	"database/sql"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/chaos-io/chaos/db"
	idredis "github.com/chaos-io/chaos/idgen/redis"
	"github.com/chaos-io/chaos/messaging"
	chaosredis "github.com/chaos-io/chaos/redis"
	"github.com/chaos-io/chaos/storage"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBuild(t *testing.T) {
	srv := miniredis.NewMiniRedis()
	require.NoError(t, srv.Start())
	t.Cleanup(srv.Close)

	storage.Register("stub", func(cfg *storage.Config) (storage.Storage, error) {
		return &stubStorage{bucket: cfg.BucketName}, nil
	})
	messaging.Register("stub", func(cfg *messaging.Config) (messaging.Queue, error) {
		_ = cfg
		return &stubQueue{}, nil
	})

	components, err := Build(Components{}, BuildOptions{
		DB: &db.Config{
			Driver: db.SqliteDriver,
			DSN:    ":memory:",
		},
		Redis: &chaosredis.Config{
			Addresses: []string{srv.Addr()},
		},
		IDGen: &idredis.Config{
			ServerIDs: []int64{1},
		},
		Storage: &storage.Config{
			Vendor:     "stub",
			BucketName: "chaos",
		},
		Messaging: &messaging.Config{
			Driver: "stub",
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, components.Close(context.Background()))
	})

	require.NotNil(t, components.DB)
	require.NotNil(t, components.Redis)
	require.NotNil(t, components.IDGen)
	require.NotNil(t, components.OSS)
	require.NotNil(t, components.Messaging)

	id, err := components.IDGen.GenID(context.Background())
	require.NoError(t, err)
	require.NotZero(t, id)
}

func TestComponentsCloseOnlyClosesOwned(t *testing.T) {
	dbProvider := &spyDBProvider{}
	redisProvider := &spyRedisProvider{}
	messagingProvider := &spyMessagingProvider{}

	components, err := Build(Components{
		OSS:       &stubStorage{},
		IDGen:     &spyIDGen{},
		DB:        dbProvider,
		Redis:     redisProvider,
		Messaging: messagingProvider,
	}, BuildOptions{})
	require.NoError(t, err)

	require.NoError(t, components.Close(context.Background()))
	require.False(t, dbProvider.closed)
	require.False(t, redisProvider.closed)
	require.False(t, messagingProvider.closed)
}

type stubStorage struct {
	bucket string
}

func (s *stubStorage) Read(context.Context, string, ...storage.Option) (*storage.Object, error) {
	return nil, nil
}

func (s *stubStorage) Write(context.Context, *storage.Object, ...storage.Option) error {
	return nil
}

func (s *stubStorage) Download(context.Context, string, string, ...storage.Option) error {
	return nil
}

func (s *stubStorage) Upload(context.Context, string, string, ...storage.Option) error {
	return nil
}

func (s *stubStorage) PresignedDownloadURL(context.Context, string, ...storage.Option) (string, error) {
	return "", nil
}

func (s *stubStorage) PresignedUploadURL(context.Context, string, ...storage.Option) (string, error) {
	return "", nil
}

type stubQueue struct{}

func (s *stubQueue) Publish(context.Context, string, ...*messaging.Message) error {
	return nil
}

func (s *stubQueue) Subscribe(*messaging.Subscription, messaging.Handler) error {
	return nil
}

func (s *stubQueue) Shutdown() {}

type spyIDGen struct{}

func (s *spyIDGen) GenID(context.Context) (int64, error) {
	return 0, nil
}

func (s *spyIDGen) GenMultiIDs(context.Context, int) ([]int64, error) {
	return nil, nil
}

type spyDBProvider struct {
	closed bool
}

func (s *spyDBProvider) Gorm() *gorm.DB {
	return nil
}

func (s *spyDBProvider) SQLDB() (*sql.DB, error) {
	return nil, nil
}

func (s *spyDBProvider) Close(context.Context) error {
	s.closed = true
	return nil
}

type spyRedisProvider struct {
	closed bool
}

func (s *spyRedisProvider) Raw() goredis.UniversalClient {
	return nil
}

func (s *spyRedisProvider) Ping(context.Context) error {
	return nil
}

func (s *spyRedisProvider) Close(context.Context) error {
	s.closed = true
	return nil
}

type spyMessagingProvider struct {
	closed bool
}

func (s *spyMessagingProvider) Publish(context.Context, string, ...*messaging.Message) error {
	return nil
}

func (s *spyMessagingProvider) Subscribe(*messaging.Subscription, messaging.Handler) error {
	return nil
}

func (s *spyMessagingProvider) Shutdown() {
	s.closed = true
}
