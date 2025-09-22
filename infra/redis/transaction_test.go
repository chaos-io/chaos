//go:build local
// +build local

package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/chaos-io/chaos/pkg/logs"
)

func Test_Trans(t *testing.T) {
	key := "transKey"

	for i := 0; i < 3; i++ {
		go trans(ctx, key)
	}

	time.Sleep(500 * time.Millisecond)
	_ = Del(ctx, key)
}

// trans 等trans结束时，val才被赋值
func trans(ctx context.Context, key string) {
	pipeline := Pipeline()
	val, _ := pipeline.Incr(ctx, key).Result()
	logs.Debugw("trans incr value", "key", key, "value", val)
	time.Sleep(100 * time.Millisecond)
	val, _ = pipeline.Decr(ctx, key).Result()
	logs.Debugw("trans decr value", "key", key, "value", val)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		logs.Warnw("trans exec err", "error", err)
	}
}

func Test_NotTrans(t *testing.T) {
	key := "notTransKey"

	for i := 0; i < 3; i++ {
		go notTrans(ctx, key)
	}

	time.Sleep(500 * time.Millisecond)
	_ = Del(ctx, key)
}

func notTrans(ctx context.Context, key string) {
	val, _ := Incr(ctx, key)
	logs.Debugw("notTrans incr value", "key", key, "value", val)
	time.Sleep(100 * time.Millisecond)
	val, _ = Decr(ctx, key)
	logs.Debugw("notTrans decr value", "key", key, "value", val)
}

func Test_Watch(t *testing.T) {
	PurchaseItem("buyer", "item", "seller", 122)
	Do(ctx, "FLUSHDB")
}

func PurchaseItem(buyerId, itemId, sellerId string, lprice int64) bool {
	buyer := fmt.Sprintf("users:%s", buyerId)
	seller := fmt.Sprintf("users:%s", sellerId)
	item := fmt.Sprintf("item:%s", itemId)
	inventory := fmt.Sprintf("inventory:%s", buyerId)
	end := time.Now().Unix() + 10

	for time.Now().Unix() < end {
		err := Watch(ctx, func(tx *Tx) error {
			if _, err := tx.TxPipelined(ctx, func(pipeliner Pipeliner) error {
				price := int64(tx.ZScore(ctx, "market:", item).Val())
				funds, _ := tx.HGet(ctx, buyer, "funds").Int64()
				if price != lprice || price > funds {
					return errors.New("can not afford this item")
				}

				pipeliner.HIncrBy(ctx, seller, "funds", price)
				pipeliner.HIncrBy(ctx, buyer, "funds", -price)
				pipeliner.SAdd(ctx, inventory, itemId)
				pipeliner.ZRem(ctx, "market:", item)
				return nil
			}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logs.Warnw("failed to do tx", "error", err)
			return false
		}

		return true
	}

	return false
}

func benchmarkUpdateToken(b *testing.B) {
	b.Run("updateToken", func(b *testing.B) {
		count := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count++
			UpdateToken("token", "user", "item")
		}
		defer Do(ctx, "FLUSHDB")
	})

	b.Run("updateTokenPipeline", func(b *testing.B) {
		count := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count++
			UpdateTokenPipeline("token", "user", "item")
		}
		defer Do(ctx, "FLUSHDB")
	})
}

/*
go test -bench .
benchmarkUpdateToken/updateToken-8                  8964            135006 ns/op
benchmarkUpdateToken/updateTokenPipeline-8         35094             34103 ns/op
*/

func UpdateToken(token, user, item string) {
	ts := float64(time.Now().UnixNano())
	_, _ = HSet(ctx, "login:", token, user)
	_, _ = ZAdd(ctx, "recent:", Z{Score: ts, Member: token})
	if len(item) > 0 {
		_, _ = ZAdd(ctx, "viewed:"+token, Z{Score: ts, Member: item})
		_, _ = ZRemRangeByRank(ctx, "viewed:"+token, 0, -26)
		_, _ = ZIncrBy(ctx, "viewed:", -1, item)
	}
}

func UpdateTokenPipeline(token, user, item string) {
	ts := float64(time.Now().UnixNano())
	pipe := Pipeline()
	pipe.HSet(ctx, "login:", token, user)
	pipe.ZAdd(ctx, "recent:", Z{Score: ts, Member: token})
	if len(item) > 0 {
		pipe.ZAdd(ctx, "viewed:"+token, Z{Score: ts, Member: item})
		pipe.ZRemRangeByRank(ctx, "viewed:"+token, 0, -26)
		pipe.ZIncrBy(ctx, "viewed:", -1, item)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		logs.Warnw("failed to do tx", "error", err)
	}
}
