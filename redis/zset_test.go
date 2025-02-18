package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ZSet(t *testing.T) {
	ctx := context.Background()
	key := "testZSetKey"
	score := 1.1
	member := "testZAdd"
	zAdd, err := ZAdd(ctx, key, score, member)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zAdd)

	zRange, err := ZRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, member, zRange[0])

	zRangeWithScores, err := ZRangeWithScores(ctx, key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, score, zRangeWithScores[0].Score)
	assert.Equal(t, member, zRangeWithScores[0].Member.(string))

	zCard, err := ZCard(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zCard)

	zScore, err := ZScore(ctx, key, member)
	assert.NoError(t, err)
	assert.Equal(t, zScore, score)

	increment := 2.0
	zIncrBy, err := ZIncrBy(ctx, key, increment, member)
	assert.NoError(t, err)
	assert.Equal(t, zIncrBy, score+increment)

	member2 := "testZAdd2"
	score2 := 10.1
	_, _ = ZAdd(ctx, key, score2, member2)

	zRank, err := ZRank(ctx, key, member2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zRank)

	// 从高到低
	zRevRange, err := ZRevRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(zRevRange))
	assert.Equal(t, member2, zRevRange[0])
	assert.Equal(t, member, zRevRange[1])

	member3 := "testZAdd3"
	score3 := 20.2
	_, _ = ZAdd(ctx, key, score3, member3)
	zRemRangeByRank, err := ZRemRangeByRank(ctx, key, -1, -1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zRemRangeByRank)

	key2 := "testZSetKey2"
	_, _ = ZAdd(ctx, key2, score2, member2)
	newKey := "testZSetNewKey"
	// zRangeScores, _ := ZRangeWithScores(ctx, key, 0, -1)
	// zRangeScores2, _ := ZRangeWithScores(ctx, key2, 0, -1)
	// fmt.Printf("zRangeScores: %v, zRangeScores2: %v\n", zRangeScores, zRangeScores2)
	// 对交集元素的分数进行求和
	zInterStore, err := ZInterStore(ctx, newKey, []string{key, key2}, nil, "")
	assert.NoError(t, err)
	zRangeScores3, _ := ZRangeWithScores(ctx, newKey, 0, -1)
	assert.Equal(t, score2+score2, zRangeScores3[0].Score)
	assert.Equal(t, member2, zRangeScores3[0].Member.(string))
	assert.Equal(t, int64(1), zInterStore)
	// fmt.Printf("zInterStore: %v, zRangeScores3: %v\n", zInterStore, zRangeScores3)

	zRem, err := ZRem(ctx, key, member, member2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), zRem)

	_, _ = ZRem(ctx, newKey, member)
	_ = Del(ctx, key, key2, newKey)
}
