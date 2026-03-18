package workers

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	repo "github.com/swastiijain24/npci-switch/internals/adapters/sqlc"
	"github.com/swastiijain24/npci-switch/internals/helpers"
	"github.com/swastiijain24/npci-switch/internals/redis_stream"
)

type CreditWorker struct {
	rdb *redis.Client
	q *repo.Queries
	r *redis_stream.EventPublisher
}

func NewCreditWorker(rdb *redis.Client, q*repo.Queries, r *redis_stream.EventPublisher) *CreditWorker{
	return &CreditWorker{
		rdb: rdb,
		q: q,
		r: r,
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

				event := message.Values["event"].(string)
				txnID := message.Values["txn_id"].(string)

				switch event {

				case "txn.initiated":
					err:= w.handleCredit(ctx, txnID)

					if err == nil{
						w.rdb.XAck(ctx, "txn_stream", "credit_group", message.ID)
					}
					
				default:
					continue

				}

				
			}
		}
	}
}

func (w *CreditWorker) handleCredit(ctx context.Context, txnID string) error {
	txn, err := w.q.GetTransaction(ctx, helpers.StringToPgUUID(txnID))
	if err !=nil{
		return err
	}

	if txn.State != "DEBIT_SUCCESS"{
		return 
	}

	w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
		ID: txn.ID,
		State: "CREDIT_PENDING",
	})

	success := mockCredit()

	if success {
		w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
			ID: txn.ID,
			State: "CREDIT_SUCCESS",
		})

		w.r.Publish(ctx, "txn_stream", map[string]interface{}{
			"txn_id": txn.ID.String(),
			"event":  "txn.credit.success",
		})

	} else {

		w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
			ID: txn.ID,
			State: "CREDIT_FAILED",
		})

		//triggering reversal 
		w.r.Publish(ctx, "txn_stream", map[string]interface{}{
			"txn_id": txn.ID.String(),
			"event": "txn.reversal",
		})
	}

	return nil
}


func mockCredit() bool {
	panic("unimplemented")
}