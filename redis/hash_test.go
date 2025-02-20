package redis

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hashField1 = "id"
	hashField2 = "name"
	hashField3 = "empty_field"

	hashValue1 = 1
	hashValue2 = "testName"
)

func Test_HSet(t *testing.T) {
	key := "testHSetKey"

	hSet, err := HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), hSet)

	_ = Del(ctx, key)
}

func Test_HGet(t *testing.T) {
	key := "testHGetKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hGet, err := HGet(ctx, key, hashField1)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(hashValue1), hGet)

	hGet2, err := HGet(ctx, key, hashField2)
	assert.NoError(t, err)
	assert.Equal(t, hashValue2, hGet2)

	hGet3, err := HGet(ctx, key, hashField3)
	assert.True(t, IsErrNil(err))
	assert.Equal(t, "", hGet3)

	_ = Del(ctx, key)
}

func Test_HMGet(t *testing.T) {
	key := "testHMGetKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hMGet, err := HMGet(ctx, key, hashField1, hashField2)
	assert.NoError(t, err)
	assert.Equal(t, []any{strconv.Itoa(hashValue1), hashValue2}, hMGet)

	_ = Del(ctx, key)
}

func Test_HGetAll(t *testing.T) {
	key := "testHGetAllKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hGetAll, err := HGetAll(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		hashField1: strconv.Itoa(hashValue1),
		hashField2: hashValue2,
	}, hGetAll)

	_ = Del(ctx, key)
}

func Test_HIncrBy(t *testing.T) {
	key := "testHIncrByKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hIncrBy, err := HIncrBy(ctx, key, hashField1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), hIncrBy)

	_ = Del(ctx, key)
}

func Test_HLen(t *testing.T) {
	key := "testHLenKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hLen, err := HLen(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), hLen)

	_ = Del(ctx, key)
}

func Test_HDel(t *testing.T) {
	key := "testHDelKey"

	_, _ = HSet(ctx, key, hashField1, hashValue1, hashField2, hashValue2)
	hDel, err := HDel(ctx, key, hashField1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), hDel)

	hLen, _ := HLen(ctx, key)
	assert.Equal(t, int64(1), hLen)

	_ = Del(ctx, key)
}
