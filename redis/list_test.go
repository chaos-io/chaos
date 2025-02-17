package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	ctx := context.Background()
	key := "testList"
	rPush, err := RPush(ctx, key, "1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rPush)
	fmt.Printf("rPush got %v\n", rPush)

	rPop, err := RPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "1", rPop)
	fmt.Printf("rPop got %v\n", rPop)

	lPush, err := LPush(ctx, key, "2")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lPush)
	fmt.Printf("lPush got %v\n", lPush)

	_, _ = LPush(ctx, key, "3")

	lRange, err := LRange(ctx, key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(lRange))
	fmt.Printf("lRange got %v\n", lRange)

	lPop, err := LPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, lPop, "3")

	lLen, err := LLen(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lLen)

	err = Del(ctx, key)
	assert.NoError(t, err)
}
