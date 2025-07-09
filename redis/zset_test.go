//go:build local
// +build local

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	zSetKey1 = "zSetKey1"
	zSetKey2 = "zSetKey2"

	zSetScore1 = 1.1
	zSetScore2 = 2.2

	zSetMember1 = "testMember1"
	zSetMember2 = "testMember2"
)

func Test_ZSet(t *testing.T) {
	zAdd, err := ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zAdd)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRange(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRange, err := ZRange(ctx, zSetKey1, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{zSetMember1, zSetMember2}, zRange)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRevRange(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRevRange, err := ZRevRange(ctx, zSetKey1, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{zSetMember2, zSetMember1}, zRevRange)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRangeWithScores(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRange, err := ZRangeWithScores(ctx, zSetKey1, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []Z{{zSetScore1, zSetMember1}, {zSetScore2, zSetMember2}}, zRange)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRevRangeWithScores(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRevRange, err := ZRevRangeWithScores(ctx, zSetKey1, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []Z{{zSetScore2, zSetMember2}, {zSetScore1, zSetMember1}}, zRevRange)

	_ = Del(ctx, zSetKey1)
}

func Test_ZCard(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zCard, err := ZCard(ctx, zSetKey1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), zCard)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRem(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRem, err := ZRem(ctx, zSetKey1, zSetMember1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zRem)

	zCard, _ := ZCard(ctx, zSetKey1)
	assert.Equal(t, int64(1), zCard)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRank(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRank, err := ZRank(ctx, zSetKey1, zSetMember1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), zRank)

	zRank2, _ := ZRank(ctx, zSetKey1, zSetMember2)
	assert.Equal(t, int64(1), zRank2)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRevRank(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRank, err := ZRevRank(ctx, zSetKey1, zSetMember1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zRank)

	zRank2, _ := ZRevRank(ctx, zSetKey1, zSetMember2)
	assert.Equal(t, int64(0), zRank2)

	_ = Del(ctx, zSetKey1)
}

func Test_ZRemRangeByRank(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zRemRangeByRank, err := ZRemRangeByRank(ctx, zSetKey1, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), zRemRangeByRank)

	zRange, _ := ZRange(ctx, zSetKey1, 0, -1)
	assert.Equal(t, []string{zSetMember1}, zRange)

	_ = Del(ctx, zSetKey1)
}

func Test_ZScore(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	zScore, err := ZScore(ctx, zSetKey1, zSetMember1)
	assert.NoError(t, err)
	assert.Equal(t, zSetScore1, zScore)

	zscore2, _ := ZScore(ctx, zSetKey1, zSetMember2)
	assert.Equal(t, zSetScore2, zscore2)

	_ = Del(ctx, zSetKey1)
}

func Test_ZIncrBy(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1})

	zIncrBy, err := ZIncrBy(ctx, zSetKey1, 2.0, zSetMember1)
	assert.NoError(t, err)
	assert.Equal(t, float64(3.1), zIncrBy)

	_ = Del(ctx, zSetKey1)
}

func Test_ZInterStore(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	_, _ = ZAdd(ctx, zSetKey2, Z{Score: 3.0, Member: zSetMember1})

	newZSetKey := "newZSetKey"
	{
		zInterStore, err := ZInterStore(ctx, newZSetKey, []string{zSetKey1, zSetKey2}, nil, "")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), zInterStore)

		zRange, _ := ZRangeWithScores(ctx, newZSetKey, 0, -1)
		assert.Equal(t, []Z{{4.1, zSetMember1}}, zRange)
	}

	newZSetKeyInMax := "newZSetKey2"
	{
		_, _ = ZInterStore(ctx, newZSetKeyInMax, []string{zSetKey1, zSetKey2}, nil, "max")
		zRange, _ := ZRangeWithScores(ctx, newZSetKeyInMax, 0, -1)
		assert.Equal(t, []Z{{3.0, zSetMember1}}, zRange)
	}

	newZSetKeyInMin := "newZSetKey3"
	{
		_, _ = ZInterStore(ctx, newZSetKeyInMin, []string{zSetKey1, zSetKey2}, nil, "min")
		zRange, _ := ZRangeWithScores(ctx, newZSetKeyInMin, 0, -1)
		assert.Equal(t, []Z{{zSetScore1, zSetMember1}}, zRange)
	}

	_ = Del(ctx, zSetKey1, zSetKey2, newZSetKey, newZSetKeyInMax, newZSetKeyInMin)
}

func Test_ZUnionStore(t *testing.T) {
	_, _ = ZAdd(ctx, zSetKey1, Z{Score: zSetScore1, Member: zSetMember1}, Z{Score: zSetScore2, Member: zSetMember2})

	_, _ = ZAdd(ctx, zSetKey2, Z{Score: 3.0, Member: zSetMember1})

	newZSetKey := "newZSetKey"
	{
		zUnionStore, err := ZUnionStore(ctx, newZSetKey, []string{zSetKey1, zSetKey2}, nil, "")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), zUnionStore)

		zRange, _ := ZRangeWithScores(ctx, newZSetKey, 0, -1)
		assert.Equal(t, []Z{{Score: zSetScore2, Member: zSetMember2}, {Score: 4.1, Member: zSetMember1}}, zRange)
	}

	newZSetKeyInMax := "newZSetKey2"
	{
		_, _ = ZUnionStore(ctx, newZSetKeyInMax, []string{zSetKey1, zSetKey2}, nil, "max")
		zRange, _ := ZRangeWithScores(ctx, newZSetKeyInMax, 0, -1)
		assert.Equal(t, []Z{{Score: zSetScore2, Member: zSetMember2}, {Score: 3.0, Member: zSetMember1}}, zRange)
	}

	newZSetKeyInMin := "newZSetKey3"
	{
		_, _ = ZUnionStore(ctx, newZSetKeyInMin, []string{zSetKey1, zSetKey2}, nil, "min")
		zRange, _ := ZRangeWithScores(ctx, newZSetKeyInMin, 0, -1)
		assert.Equal(t, []Z{{Score: zSetScore1, Member: zSetMember1}, {Score: zSetScore2, Member: zSetMember2}}, zRange)
	}

	_ = Del(ctx, zSetKey1, zSetKey2, newZSetKey, newZSetKeyInMax, newZSetKeyInMin)
}
