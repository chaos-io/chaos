//go:build local
// +build local

package redis

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	member1 = "member1"
	member2 = 2
)

func Test_SAdd(t *testing.T) {
	key := "testSAddKey"

	sAdd, err := SAdd(ctx, key, member1, member2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), sAdd)

	_ = Del(ctx, key)
}

func Test_SMembers(t *testing.T) {
	key := "testSMembersKey"

	_, _ = SAdd(ctx, key, member1, member2)
	sMembers, err := SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, []string{member1, strconv.Itoa(member2)}, sMembers)

	_ = Del(ctx, key)
}

func Test_SCard(t *testing.T) {
	key := "testSCardKey"

	_, _ = SAdd(ctx, key, member1, member2)
	sCard, err := SCard(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), sCard)

	_ = Del(ctx, key)
}

func Test_SIsMember(t *testing.T) {
	key := "testSIsMemberKey"

	_, _ = SAdd(ctx, key, member1)
	sIsMember, err := SIsMember(ctx, key, member1)
	assert.NoError(t, err)
	assert.Equal(t, true, sIsMember)

	sIsMember2, err := SIsMember(ctx, key, true)
	assert.NoError(t, err)
	assert.Equal(t, false, sIsMember2)

	_ = Del(ctx, key)
}

func Test_SPop(t *testing.T) {
	key := "testSPopKey"

	_, _ = SAdd(ctx, key, member1)
	sPop, err := SPop(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, member1, sPop)

	sPop2, err := SPop(ctx, key)
	assert.True(t, IsErrNil(err))
	assert.Equal(t, "", sPop2)

	_ = Del(ctx, key)
}

func Test_SRandMember(t *testing.T) {
	key := "testSRandMemberKey"

	_, _ = SAdd(ctx, key, member1)
	sRandMember, err := SRandMember(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, member1, sRandMember)

	_ = Del(ctx, key)
}

func Test_SRem(t *testing.T) {
	key := "testSRemKey"

	_, _ = SAdd(ctx, key, member1)
	sRem, err := SRem(ctx, key, member1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), sRem)

	sIsMember, _ := SIsMember(ctx, key, member1)
	assert.Equal(t, false, sIsMember)

	_ = Del(ctx, key)
}

func Test_SDiff(t *testing.T) {
	key := "testSDiffKey"
	key2 := "testSDiffKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)

	sDiff, err := SDiff(ctx, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, []string{strconv.Itoa(member2)}, sDiff)

	_ = Del(ctx, key, key2)
}

func Test_SDiffStore(t *testing.T) {
	key := "testSDiffStoreKey"
	key2 := "testSDiffStoreKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)

	destKey := "testSDiffStoreKey3"
	sDiffStore, err := SDiffStore(ctx, destKey, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), sDiffStore)

	sMembers, _ := SMembers(ctx, destKey)
	assert.Equal(t, []string{strconv.Itoa(member2)}, sMembers)

	_ = Del(ctx, key, key2)
}

func Test_SInter(t *testing.T) {
	key := "testSInterKey"
	key2 := "testSInterKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)

	sInter, err := SInter(ctx, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, []string{member1}, sInter)

	_ = Del(ctx, key, key2)
}

func Test_SInterStore(t *testing.T) {
	key := "testSInterStoreKey"
	key2 := "testSInterStoreKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)

	destKey := "testSInterStoreKey3"
	sInterStore, err := SInterStore(ctx, destKey, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), sInterStore)

	sMembers, _ := SMembers(ctx, destKey)
	assert.Equal(t, []string{member1}, sMembers)

	_ = Del(ctx, key, key2)
}

func Test_SUnion(t *testing.T) {
	key := "testSUnionKey"
	key2 := "testSUnionKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)
	sUnion, err := SUnion(ctx, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, []string{member1, strconv.Itoa(member2)}, sUnion)

	_ = Del(ctx, key, key2)
}

func Test_SUnionStore(t *testing.T) {
	key := "testSUnionStoreKey"
	key2 := "testSUnionStoreKey2"

	_, _ = SAdd(ctx, key, member1, member2)
	_, _ = SAdd(ctx, key2, member1)
	destKey := "testSUnionStoreKey3"
	sUnionStore, err := SUnionStore(ctx, destKey, key, key2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), sUnionStore)

	sMembers, _ := SMembers(ctx, destKey)
	assert.Equal(t, []string{member1, strconv.Itoa(member2)}, sMembers)

	_ = Del(ctx, key, key2)
}
