package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // Redis密码，没有密码则留空
		DB:       0,                // 使用默认的数据库
	})

	ctx := context.Background()
	pubSub := rdb.Subscribe(ctx, "channel")

	defer func() {
		err := pubSub.Close()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for msg := range pubSub.Channel() {
			fmt.Println("received message:", msg.Payload)
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 1000)
	for {
		select {
		case <-ticker.C:
			// err := rdb.Publish(ctx, "channel", "hello world").Err()
			// if err != nil {
			// 	panic(err)
			// }
		}
	}

}
