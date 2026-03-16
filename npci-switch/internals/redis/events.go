package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type EventPublisher struct {
	rdb *redis.Client
}

func NewPublisher(rdb *redis.Client) *EventPublisher{
	return &EventPublisher{
		rdb: rdb,
	}
}

func (p* EventPublisher) Publish(ctx context.Context, stream string, values map[string]interface{}) error{

	return p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID: "*",
		Values: values,
	}).Err()
}