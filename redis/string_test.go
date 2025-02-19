package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
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

func TestStrings(t *testing.T) {
	ctx := context.Background()
	key := "testStringsKey"

	str1, err := Append(ctx, key, "hello ")
	assert.NoError(t, err)
	assert.Equal(t, int64(6), str1)
	str2, err := Append(ctx, key, "world!")
	assert.NoError(t, err)
	assert.Equal(t, int64(12), str2)
	get1, _ := Get(ctx, key)
	assert.Equal(t, "hello world!", get1)

	getRange, err := GetRange(ctx, key, 3, 7)
	assert.NoError(t, err)
	assert.Equal(t, "lo wo", getRange)

	setRange, err := SetRange(ctx, key, 0, "H")
	assert.NoError(t, err)
	assert.Equal(t, int64(12), setRange)
	get2, _ := Get(ctx, key)
	assert.Equal(t, "Hello world!", get2)

	_ = Del(ctx, key)
}

func TestBit(t *testing.T) {
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
