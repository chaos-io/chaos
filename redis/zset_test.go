package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ZAdd(t *testing.T) {
	ctx := context.Background()
	key := "testZAdd"
	score := 1.1
	member := "testZAdd"
	zAdd, err := ZAdd(ctx, key, score, member)
	assert.NoError(t, err)
	fmt.Println(zAdd)

	zRange, err := ZRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	fmt.Println(zRange)

	zCard, err := ZCard(ctx, key)
	assert.NoError(t, err)
	fmt.Println(zCard)

	zScore, err := ZScore(ctx, key, member)
	assert.NoError(t, err)
	assert.Equal(t, zScore, score)

	increment := 2.2
	zIncrBy, err := ZIncrBy(ctx, key, increment, member)
	assert.NoError(t, err)
	assert.Equal(t, zIncrBy, score+increment)

	member2 := "testZAdd2"
	_, _ = ZAdd(ctx, key, score, member2)
	zRevRange, err := ZRevRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	fmt.Println(zRevRange)

	key2 := "group:"
	_, _ = ZAdd(ctx, key2, score, member)
	newKey := "newKey"
	zRange, _ = ZRange(ctx, key, 0, -1)
	zRange2, _ := ZRange(ctx, key2, 0, -1)
	fmt.Printf("zRange: %v, zRange2: %v\n", zRange, zRange2)
	zInterStore, err := ZInterStore(ctx, newKey, []string{key, key2}, nil, "")
	assert.NoError(t, err)
	zRange3, _ := ZRange(ctx, newKey, 0, -1)
	fmt.Printf("zInterStore: %v, zRange3: %v\n", zInterStore, zRange3)

	zRem, err := ZRem(ctx, key, member, member2)
	assert.NoError(t, err)
	fmt.Println(zRem)

	_, _ = ZRem(ctx, newKey, member)
}
