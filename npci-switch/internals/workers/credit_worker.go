package workers

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type CreditWorker struct {
	rdb *redis.Client
}

func NewCreditWorker(rdb *redis.Client) *CreditWorker{
	return &CreditWorker{
		rdb: rdb,
	}
}

func (w *CreditWorker) Start(ctx context.Context){
	for{
		streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "credit_group",
			Consumer: "worker-1",
			Streams:  []string{"txn_stream", ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err!=nil{
			log.Println(err)
			continue
		}

		for _, stream := range streams{
			for _, message := range stream.Messages{
				log.Println("Credit worker received:", message.Values)

				// TODO credit receiver bank
				//bank client ka credit 

				w.rdb.XAck(ctx, "txn_stream", "credit_group", message.ID)
			}
		}
	}
}