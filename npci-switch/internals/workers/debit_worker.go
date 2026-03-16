package workers

//worker is responsible for processing debit operations.

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type DebitWorker struct {
	rdb *redis.Client
}

func NewDebitWorker(rdb *redis.Client) *DebitWorker{
	return &DebitWorker{
		rdb: rdb,
	}
}

func (w *DebitWorker) Start(ctx context.Context){

	//infinite loop - continuously waits for events, processes them and acknowledges them

	for{
		streams, err:= w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group: "debit_group", // a grp can have multiple workers where messages are distributed between workers helps in horizontal sclaing
			Consumer: "worker-1",
			Streams: []string{"txn_stream", ">"}, // read from txn_stream and > means read only which were never delivered before , only new messgaes 
			Count: 10,	//read upto 10 msgs 
			Block: 0, // wait until event arrives 
		}).Result()

		if err !=nil{
			log.Println(err)
			continue
		}

		for _, stream := range streams{
			for _, message := range stream.Messages{
				log.Println("debit worker received", message.Values)

				//todo call debit bank service 
				//will call bank client ka debit (with sender id and amount and bank)
				//ack message
				w.rdb.XAck(ctx, "txn_stream", "debit_group", message.ID)
			}
		}
	}
}