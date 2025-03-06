package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	"github.com/chaos-io/chaos/logs"
)

func XAdd(ctx context.Context, streamName string, kvs ...any) (string, error) {
	res, err := GetRedis().XAdd(ctx, &goredis.XAddArgs{
		Stream: streamName,
		Values: kvs,
	}).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func XRead(ctx context.Context, streamName string) ([]goredis.XStream, error) {
	xStreams, err := GetRedis().XRead(ctx, &goredis.XReadArgs{
		Streams: []string{streamName, "0"},
		Count:   10,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	return xStreams, nil
}

func XGroupCreate(ctx context.Context, streamName, groupName string) (string, error) {
	groups, err := GetRedis().XInfoGroups(ctx, streamName).Result()
	if err != nil {
		return "", err
	}

	for _, group := range groups {
		if group.Name == groupName {
			logs.Debugw("stream's group already exists", "stream", streamName, "group", groupName)
			return "", nil
		}
	}

	res, err := GetRedis().XGroupCreateMkStream(ctx, streamName, groupName, "0").Result()
	if err != nil {
		return "", err
	}

	logs.Debugw("stream's group created", "stream", streamName, "group", groupName)
	return res, nil
}

func XReadGroup(ctx context.Context, streamName, groupName, consumer string) ([]goredis.XStream, error) {
	xStreams, err := GetRedis().XReadGroup(ctx, &goredis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumer,
		Streams:  []string{streamName, ">"},
		Count:    1,
		Block:    0,
		NoAck:    false,
	}).Result()
	if err != nil {
		return nil, err
	}

	return xStreams, nil
}

func XPending(ctx context.Context, streamName, groupName string) (*goredis.XPending, error) {
	xPending, err := GetRedis().XPending(ctx, streamName, groupName).Result()
	if err != nil {
		return nil, err
	}

	return xPending, nil
}

func XAck(ctx context.Context, streamName, groupName string, ids ...string) (int64, error) {
	num, err := GetRedis().XAck(ctx, streamName, groupName, ids...).Result()
	if err != nil {
		return 0, err
	}

	return num, nil
}
