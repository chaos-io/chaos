package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHSet(t *testing.T) {
	ctx := context.Background()
	key := "testHSet"
	got, err := HSet(ctx, key, "name", "testName", "id", 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), got)

	hLen, err := HLen(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), hLen)

	all, err := HGetAll(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"name": "testName", "id": "1"}, all)

	name, err := HGet(ctx, key, "name")
	assert.NoError(t, err)
	assert.Equal(t, "testName", name)

	hmGet, err := HMGet(ctx, key, "id", "name")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(hmGet))
	// fmt.Printf("hMGet got %v\n", hmGet)

	hIncrBy, err := HIncrBy(ctx, key, "id", 2)
	assert.NoError(t, err)
	assert.Equal(t, hIncrBy, int64(3))
	// fmt.Printf("hIncrBy got %v\n", hIncrBy)

	_ = Del(ctx, key)
}
