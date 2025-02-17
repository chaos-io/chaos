package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	ctx := context.Background()

	key := "testStringKey"
	set, err := Set(ctx, key, "value1", 0)
	assert.NoError(t, err)
	assert.NotEqual(t, "1", set)

	get, err := Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "value1", get)

	key2 := "testStringKey2"
	setNX, err := SetNX(ctx, key2, "value2", 2*time.Second)
	assert.NoError(t, err)
	assert.True(t, setNX)

	get2, _ := Get(ctx, key2)
	assert.Equal(t, "value2", get2)
	time.Sleep(2 * time.Second)
	_, err = Get(ctx, key2)
	assert.True(t, IsErrNil(err))

	key3 := "testStringKey3"
	incr, err := Incr(ctx, key3)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), incr)

	incrBy, err := IncrBy(ctx, key3, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), incrBy)

	err = Del(ctx, key, key2, key3)
	assert.NoError(t, err)
}
