package redis_idgen

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRedis(t *testing.T) (*goredis.Client, func(), error) {
	m := miniredis.NewMiniRedis()
	if err := m.Start(); err != nil {
		return nil, nil, err
	}

	opts := &goredis.Options{Addr: m.Addr()}
	cli := goredis.NewClient(opts)

	cleanup := func() {
		_ = cli.Close()
		m.Close()
	}

	return cli, cleanup, nil
}

func Test_generator_GenMultiIDs(t *testing.T) {
	ctx := context.Background()

	cli, cleanup, err := newTestRedis(t)
	require.NoError(t, err)
	t.Cleanup(cleanup)

	idgen, err := NewIDGenerator(cli, []int64{0, 1, 2})
	require.NoError(t, err)

	ids, err := idgen.GenMultiIDs(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(ids))

	id, err := idgen.GenID(ctx)
	assert.Nil(t, err)
	assert.True(t, id >= time.Now().UnixNano()-(id>>32))
}

func TestNewIDGeneratorValidateArgs(t *testing.T) {
	idgen, err := NewIDGenerator(nil, []int64{1})
	require.ErrorIs(t, err, ErrNilStore)
	assert.Nil(t, idgen)

	cli, cleanup, err := newTestRedis(t)
	require.NoError(t, err)
	t.Cleanup(cleanup)

	idgen, err = NewIDGenerator(cli, nil)
	require.ErrorIs(t, err, ErrEmptyServerID)
	assert.Nil(t, idgen)
}
