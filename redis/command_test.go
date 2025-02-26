package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Do(t *testing.T) {
	do, err := Do(ctx, "PING")
	assert.NoError(t, err)
	assert.Equal(t, "PONG", do)
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
