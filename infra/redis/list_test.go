//go:build local
// +build local

package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/chaos-io/chaos/pkg/logs"
)

func Test_RPush(t *testing.T) {
	key := "rPushKey"
	rPush, err := RPush(ctx, key, 1, "2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), rPush)

	_ = Del(ctx, key)
}

func Test_RPop(t *testing.T) {
	key := "rPopKey"
	_, _ = RPush(ctx, key, 1)

	rPop, err := RPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "1", rPop)

	_ = Del(ctx, key)
}

func Test_LPush(t *testing.T) {
	key := "lPushKey"
	lPush, err := LPush(ctx, key, 1, "2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), lPush)

	_ = Del(ctx, key)
}

func Test_LPop(t *testing.T) {
	key := "lPopKey"
	_, _ = LPush(ctx, key, 1)

	rPop, err := LPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "1", rPop)

	_ = Del(ctx, key)
}

func Test_LIndex(t *testing.T) {
	key := "lIndexKey"
	_, _ = LPush(ctx, key, 1, "2")

	lIndex1, err1 := LIndex(ctx, key, 0)
	assert.NoError(t, err1)
	assert.Equal(t, "2", lIndex1)

	lIndex2, err2 := LIndex(ctx, key, 1)
	assert.NoError(t, err2)
	assert.Equal(t, "1", lIndex2)

	_ = Del(ctx, key)
}

func Test_LRange(t *testing.T) {
	key := "lRangeKey"
	_, _ = RPush(ctx, key, 1, "2")

	lRange, err := LRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2"}, lRange)

	_ = Del(ctx, key)
}

func Test_LLen(t *testing.T) {
	key := "lLenKey"
	_, _ = RPush(ctx, key, 1, "2")

	lLen, err := LLen(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), lLen)

	_ = Del(ctx, key)
}

func Test_LTrim(t *testing.T) {
	key := "testLTrim"
	_, _ = RPush(ctx, key, 1, "2", 3)

	lTrim, err := LTrim(ctx, key, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "OK", lTrim)
	val1, _ := LRange(ctx, key, 0, -1)
	assert.Equal(t, []string{"1", "2"}, val1)

	_, _ = LTrim(ctx, key, 0, 0)
	val2, _ := LRange(ctx, key, 0, -1)
	assert.Equal(t, []string{"1"}, val2)

	lTrim3, err := LTrim(ctx, key, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, "OK", lTrim3)
	val3, _ := LRange(ctx, key, 0, -1)
	assert.Equal(t, []string{}, val3)

	_ = Del(ctx, key)
}

func Test_BRPop(t *testing.T) {
	key := "testBRPop"
	_, _ = RPush(ctx, key, 1, "2")

	key2 := "testBRPop2"
	_, _ = RPush(ctx, key2, 3)

	brPop, err := BRPop(ctx, 1*time.Second, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, []string{key, "2"}, brPop)

	brPop2, _ := BRPop(ctx, 1*time.Second, key, key2)
	assert.Equal(t, []string{key, "1"}, brPop2)

	brPop3, _ := BRPop(ctx, 1*time.Second, key, key2)
	assert.Equal(t, []string{key2, "3"}, brPop3)

	go func() {
		time.Sleep(2 * time.Second)
		logs.Debugw("rPush key2 new value in goroutine")
		_, _ = RPush(ctx, key2, "4")
	}()

	{
		logs.Debugw("empty list, wait to block, and can't get value")
		brpop4, err := BRPop(ctx, 1*time.Second, key, key2)
		logs.Debugw("brpop ok")
		assert.Equal(t, true, IsErrNil(err))
		assert.Equal(t, 0, len(brpop4))
	}

	{
		logs.Debugw("empty list, wait to block, can get value")
		brpop5, err := BRPop(ctx, 3*time.Second, key, key2)
		logs.Debugw("brpop ok2")
		assert.NoError(t, err)
		assert.Equal(t, []string{key2, "4"}, brpop5)
	}

	_ = Del(ctx, key, key2)
}

func Test_RPopLPush(t *testing.T) {
	key := "testRPopLPush"
	_, _ = RPush(ctx, key, 1, "2", 3)
	key2 := "testRPopLPush2"
	_, _ = RPush(ctx, key2, "a")

	rPopLPush, err := RPopLPush(ctx, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, "3", rPopLPush)

	lRange, _ := LRange(ctx, key, 0, -1)
	assert.Equal(t, []string{"1", "2"}, lRange)

	lRange2, _ := LRange(ctx, key2, 0, -1)
	assert.Equal(t, []string{"3", "a"}, lRange2)

	_ = Del(ctx, key, key2)
}
