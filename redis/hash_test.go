package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHSet(t *testing.T) {
	ctx := context.Background()
	key := "testHSet"
	got, err := HSet(ctx, key, "name", "testName", "id", 1)
	assert.NoError(t, err)
	fmt.Printf("hSet got %v\n", got)

	all, err := HGetAll(ctx, key)
	assert.NoError(t, err)
	fmt.Printf("hGetAll got %v\n", all)

	name, err := HGet(ctx, key, "name")
	assert.NoError(t, err)
	assert.Equal(t, name, "testName")

	hmGet, err := HMGet(ctx, key, "id", "name")
	assert.NoError(t, err)
	assert.Equal(t, len(hmGet), 2)
	fmt.Printf("hMGet got %v\n", hmGet)

	hIncrBy, err := HIncrBy(ctx, key, "id", 2)
	assert.NoError(t, err)
	assert.Equal(t, hIncrBy, int64(3))
	fmt.Printf("hIncrBy got %v\n", hIncrBy)
}
