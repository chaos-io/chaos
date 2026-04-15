package redis

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMiniRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	srv := miniredis.NewMiniRedis()
	require.NoError(t, srv.Start())
	t.Cleanup(srv.Close)
	return srv
}

func TestNewRequiresAddresses(t *testing.T) {
	cli, err := New(Config{})
	require.ErrorIs(t, err, ErrEmptyAddresses)
	assert.Nil(t, cli)
}

func TestNewTrimsAndDedupsAddresses(t *testing.T) {
	cli, err := New(Config{Addresses: []string{" 127.0.0.1:6379 ", "127.0.0.1:6379", "127.0.0.1:6380"}})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cli.Close(context.Background()))
	})

	cluster, ok := cli.Raw().(*goredis.ClusterClient)
	require.True(t, ok)
	assert.Equal(t, []string{"127.0.0.1:6379", "127.0.0.1:6380"}, cluster.Options().Addrs)
}

func TestNewRejectsClusterDB(t *testing.T) {
	cli, err := New(Config{
		Addresses: []string{"127.0.0.1:6379", "127.0.0.1:6380"},
		DB:        1,
	})
	require.ErrorIs(t, err, ErrClusterDBUnsupported)
	assert.Nil(t, cli)
}

func TestNewRejectsInvalidBackoffRange(t *testing.T) {
	cli, err := New(Config{
		Addresses:       []string{"127.0.0.1:6379"},
		MinRetryBackoff: time.Second,
		MaxRetryBackoff: 500 * time.Millisecond,
	})
	require.ErrorIs(t, err, ErrInvalidBackoff)
	assert.Nil(t, cli)
}

func TestNewBuildsWorkingService(t *testing.T) {
	srv := newMiniRedis(t)

	cli, err := New(Config{Addresses: []string{srv.Addr()}})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cli.Close(context.Background()))
	})

	ctx := context.Background()
	require.NoError(t, cli.Raw().Set(ctx, "user:1:name", "alice", time.Minute).Err())

	value, err := cli.Raw().Get(ctx, "user:1:name").Result()
	require.NoError(t, err)
	assert.Equal(t, "alice", value)
}

func TestNewAppliesDefaults(t *testing.T) {
	cli, err := New(Config{Addresses: []string{"127.0.0.1:6379"}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, cli.Close(context.Background())) })

	raw, ok := cli.Raw().(*goredis.Client)
	require.True(t, ok)
	opts := raw.Options()
	assert.Equal(t, 3, opts.MaxRetries)
	assert.Equal(t, 8*time.Millisecond, opts.MinRetryBackoff)
	assert.Equal(t, 512*time.Millisecond, opts.MaxRetryBackoff)
	assert.Equal(t, 10*runtime.GOMAXPROCS(0), opts.PoolSize)
	assert.Equal(t, time.Second*3, opts.ReadTimeout)
	assert.Equal(t, opts.ReadTimeout, opts.WriteTimeout)
}

func TestWrapRejectsNil(t *testing.T) {
	cli, err := Wrap(nil)
	require.ErrorIs(t, err, ErrNilRawClient)
	assert.Nil(t, cli)
}

func TestPingWorks(t *testing.T) {
	srv := newMiniRedis(t)

	cli, err := New(Config{Addresses: []string{srv.Addr()}})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cli.Close(context.Background()))
	})

	require.NoError(t, cli.Ping(context.Background()))
}

func TestPingNilService(t *testing.T) {
	var svc *Service
	require.ErrorIs(t, svc.Ping(context.Background()), ErrNilService)
}

func TestCloseNilSafe(t *testing.T) {
	var svc *Service
	require.NoError(t, svc.Close(context.Background()))
}
