//go:build local
// +build local

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	key := "setKey"

	set, err := Set(ctx, key, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, "OK", set)

	_ = Del(ctx, key)
}

func Test_Get(t *testing.T) {
	key := "getKey"

	_, _ = Set(ctx, key, "a", 1*time.Second)

	get, err := Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "a", get)

	time.Sleep(1 * time.Second)
	get2, err2 := Get(ctx, key)
	assert.Equal(t, IsErrNil(err2), true)
	assert.Equal(t, "", get2)

	_ = Del(ctx, key)
}

func Test_SetNX(t *testing.T) {
	key := "setNXKey"

	setNX, err := SetNX(ctx, key, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, true, setNX)

	setNX2, err2 := SetNX(ctx, key, "a", 0)
	assert.NoError(t, err2)
	assert.Equal(t, false, setNX2)

	get, _ := Get(ctx, key)
	assert.Equal(t, "1", get)

	_ = Del(ctx, key)
}

func Test_Incr(t *testing.T) {
	key := "incrKey"

	incr, err := Incr(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), incr)

	incr2, err2 := Incr(ctx, key)
	assert.NoError(t, err2)
	assert.Equal(t, int64(2), incr2)

	_ = Del(ctx, key)
}

func Test_IncrBy(t *testing.T) {
	key := "incrByKey"

	incrBy, err := IncrBy(ctx, key, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), incrBy)

	incrBy2, err2 := IncrBy(ctx, key, 3)
	assert.NoError(t, err2)
	assert.Equal(t, int64(5), incrBy2)

	_ = Del(ctx, key)
}

func Test_incrByFloat(t *testing.T) {
	key := "incrByFloatKey"

	incrByFloat, err := IncrByFloat(ctx, key, 2.3)
	assert.NoError(t, err)
	assert.Equal(t, float64(2.3), incrByFloat)

	incrByFloat2, err2 := IncrByFloat(ctx, key, 3.0)
	assert.NoError(t, err2)
	assert.Equal(t, float64(5.3), incrByFloat2)

	_ = Del(ctx, key)
}

func Test_Decr(t *testing.T) {
	key := "decrKey"

	decr, err := Decr(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), decr)

	_ = Del(ctx, key)
}

func Test_DecrBy(t *testing.T) {
	key := "decrByKey"

	decr, err := DecrBy(ctx, key, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), decr)

	_ = Del(ctx, key)
}

func Test_Append(t *testing.T) {
	key := "appendKey"

	str1, err := Append(ctx, key, "hello ")
	assert.NoError(t, err)
	assert.Equal(t, int64(6), str1)
	str2, err2 := Append(ctx, key, "world!")
	assert.NoError(t, err2)
	assert.Equal(t, int64(12), str2)
	get, _ := Get(ctx, key)
	assert.Equal(t, "hello world!", get)

	_ = Del(ctx, key)
}

func Test_GetRange(t *testing.T) {
	key := "getRangeKey"

	_, _ = Set(ctx, key, "hello world!", 0)
	getRange, err := GetRange(ctx, key, 3, 7)
	assert.NoError(t, err)
	assert.Equal(t, "lo wo", getRange)

	_ = Del(ctx, key)
}

func Test_SetRange(t *testing.T) {
	key := "setRangeKey"

	_, _ = Set(ctx, key, "hello world!", 0)
	setRange, err := SetRange(ctx, key, 0, "H")
	assert.NoError(t, err)
	assert.Equal(t, int64(12), setRange)
	get, _ := Get(ctx, key)
	assert.Equal(t, "Hello world!", get)

	_ = Del(ctx, key)
}

func Test_Bit(t *testing.T) {
	ctx := context.Background()
	key := "testBitKey"

	bit, err := SetBit(ctx, key, 2, 1) // 0100 0000
	assert.NoError(t, err)
	assert.Equal(t, int64(0), bit)
	bit2, err := SetBit(ctx, key, 7, 1) // 0100 0001
	assert.NoError(t, err)
	assert.Equal(t, int64(0), bit2)

	getBit, err := GetBit(ctx, key, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), getBit)
	getBit2, _ := GetBit(ctx, key, 2)
	assert.Equal(t, int64(1), getBit2)

	bitCount, err := BitCount(ctx, key, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), bitCount)

	get, _ := Get(ctx, key)
	assert.Equal(t, "!", get)

	_ = Del(ctx, key)
}
