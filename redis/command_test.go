package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Do(t *testing.T) {
	do, err := Do(ctx, "PING")
	assert.NoError(t, err)
	assert.Equal(t, "PONG", do)
}

func Test_Del(t *testing.T) {
	key := "delKey"
	err := Del(ctx, key)
	assert.NoError(t, err)
}

func Test_Exists(t *testing.T) {
	key := "existsKey"
	result, err := Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, result)

	_, _ = Incr(ctx, key)
	result2, err2 := Exists(ctx, key)
	assert.NoError(t, err2)
	assert.True(t, result2)
}

func Test_Type(t *testing.T) {
	stringKey := "stringKey"

	_, _ = Set(ctx, stringKey, "a", 0)
	strType, err := Type(ctx, stringKey)
	assert.NoError(t, err)
	assert.Equal(t, "string", strType)

	listKey := "listKey"
	_, _ = RPush(ctx, listKey, 1)
	listType, err2 := Type(ctx, listKey)
	assert.NoError(t, err2)
	assert.Equal(t, "list", listType)

	setKey := "setKey"
	_, _ = SAdd(ctx, setKey, 1)
	setType, err3 := Type(ctx, setKey)
	assert.NoError(t, err3)
	assert.Equal(t, "set", setType)

	hashKey := "hashKey"
	_, _ = HSet(ctx, hashKey, "a", 1)
	hashType, err4 := Type(ctx, hashKey)
	assert.NoError(t, err4)
	assert.Equal(t, "hash", hashType)

	zsetKey := "zsetKey"
	_, _ = ZAdd(ctx, zsetKey, Z{Score: 10, Member: "a"})
	zsetType, err5 := Type(ctx, zsetKey)
	assert.NoError(t, err5)
	assert.Equal(t, "zset", zsetType)

	_ = Del(ctx, stringKey, listKey, setKey, hashKey, zsetKey)
}

func Test_Expire(t *testing.T) {
	key := "expireKey"

	_, _ = Incr(ctx, key)
	ttl, err := TTL(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-1), ttl)

	expire, err := Expire(ctx, key, 1*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, true, expire)

	_ = Del(ctx, key)
}

func Test_ExpireAt(t *testing.T) {
	key := "expireAtKey"

	_, _ = Incr(ctx, key)
	now := time.Now()
	expireTime := now.Add(1 * time.Second)

	expireAt, err := ExpireAt(ctx, key, expireTime)
	assert.NoError(t, err)
	assert.Equal(t, true, expireAt)

	ttl, err := PTTL(ctx, key)
	assert.NoError(t, err)
	assert.True(t, ttl > 0)

	time.Sleep(1 * time.Second)
	ttl2, err := PTTL(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-2), ttl2)

	_ = Del(ctx, key)
}

func Test_Persist(t *testing.T) {
	key := "persistKey"

	_, _ = Incr(ctx, key)
	_, _ = Expire(ctx, key, 1*time.Second)
	pttl, _ := PTTL(ctx, key)
	assert.True(t, pttl > 0)

	persist, err := Persist(ctx, key)
	assert.NoError(t, err)
	assert.True(t, persist)

	pttl2, err2 := PTTL(ctx, key)
	assert.NoError(t, err2)
	assert.Equal(t, time.Duration(-1), pttl2)

	_ = Del(ctx, key)
}

func Test_Sort_List(t *testing.T) {
	key := "sortListKey"
	_, _ = RPush(ctx, key, "3")
	_, _ = RPush(ctx, key, "2")
	_, _ = RPush(ctx, key, "8")
	sort, err := Sort(ctx, key, "", nil, "", 0, 0, true)
	assert.NoError(t, err)
	assert.Equal(t, []string{"2", "3", "8"}, sort)

	_ = Del(ctx, key)
}

func Test_Sort_Set(t *testing.T) {
	key := "sortSetKey"
	_, _ = SAdd(ctx, key, "3", "2", "8")

	sort, err := Sort(ctx, key, "", nil, "", 0, 0, true)
	assert.NoError(t, err)
	assert.Equal(t, []string{"2", "3", "8"}, sort)

	_ = Del(ctx, key)
}

func Test_Sort_ZSet(t *testing.T) {
	key := "sortZSetKey"

	_, _ = ZAdd(ctx, key, Z{Score: 10.0, Member: "zsetMemeber1"}, Z{Score: 2.0, Member: "zsetMemeber2"})

	sort, err := Sort(ctx, key, "", nil, "", 0, 0, true)
	assert.NoError(t, err)
	assert.Equal(t, []string{"zsetMemeber1", "zsetMemeber2"}, sort)

	sort2, err2 := Sort(ctx, key, "score", nil, "", 0, 0, true)
	assert.NoError(t, err2)
	assert.Equal(t, []string{"zsetMemeber2", "zsetMemeber1"}, sort2)

	_ = Del(ctx, key)
}
