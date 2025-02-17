package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	ctx := context.Background()
	key := "setKey"
	member1 := "member1"
	member2 := 2

	sAdd, err := SAdd(ctx, key, member1, member2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), sAdd)
	fmt.Printf("sAdd: %v\n", sAdd)

	sMembers, err := SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sMembers))
	fmt.Printf("sMembers: %v\n", sMembers)

	sCard, err := SCard(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), sCard)
	fmt.Printf("sCard: %v\n", sCard)

	sIsMember, err := SIsMember(ctx, key, member1)
	assert.NoError(t, err)
	assert.Equal(t, true, sIsMember)
	fmt.Printf("sIsMember: %v\n", sIsMember)

	sRem, err := SRem(ctx, key, member1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), sRem)
	fmt.Printf("sRem: %v\n", sRem)

	sMembers, _ = SMembers(ctx, key)
	fmt.Printf("sMembers2: %v\n", sMembers)

	sRandMember, err := SRandMember(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(member2), sRandMember)
	fmt.Printf("sRandMember: %v\n", sRandMember)

	sPop, err := SPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(member2), sPop)
	fmt.Printf("sPop: %v\n", sPop)

	sMembers, _ = SMembers(ctx, key)
	fmt.Printf("sMembers3: %v\n", sMembers)

	err = Del(ctx, key)
	assert.NoError(t, err)
}
