package workers

//this worker is responsible for processing debit operations.

import (
	"context"
	"github.com/redis/go-redis/v9"
	repo "github.com/swastiijain24/npci-switch/internals/adapters/sqlc"
	"github.com/swastiijain24/npci-switch/internals/helpers"
	"github.com/swastiijain24/npci-switch/internals/redis_stream"
	"log"
)

type DebitWorker struct {
	rdb *redis.Client
	q   *repo.Queries
	r   *redis_stream.EventPublisher
}

func NewDebitWorker(rdb *redis.Client, q *repo.Queries, r *redis_stream.EventPublisher) *DebitWorker {
	return &DebitWorker{
		rdb: rdb,
		q:   q,
		r:   r,
	}
}

func (w *DebitWorker) Start(ctx context.Context) {

	//infinite loop - continuously waits for events, processes them and acknowledges them

	for {
		streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "debit_group", // a grp can have multiple workers where messages are distributed between workers helps in horizontal sclaing
			Consumer: "worker-1",
			Streams:  []string{"txn_stream", ">"}, // read from txn_stream and > means read only which were never delivered before , only new messgaes
			Count:    10,                          //read upto 10 msgs
			Block:    0,                           // wait until event arrives
		}).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				log.Println("debit worker received", message.Values)

				event, ok := message.Values["event"].(string)
				if !ok {
					continue
				}

				txnID, ok := message.Values["txn_id"].(string)
				if !ok {
					continue
				}

				switch event {

				case "txn.initiated":
					err:= w.handleDebit(ctx, txnID)

					if err == nil{
						w.rdb.XAck(ctx, "txn_stream", "debit_group", message.ID)
					}

				default:
					continue

				}
			}
		}
	}
}

func (w *DebitWorker) handleDebit(ctx context.Context, txnID string) error {

	txn, err := w.q.GetTransaction(ctx, helpers.StringToPgUUID(txnID))
	if err != nil {
		return err
	}

	if txn.State != "INITIATED" {
		return 
	}

	w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
		ID:    txn.ID,
		State: "DEBIT_PENDING",
	})

	success := mockDebit() //call to bankservice /debit for the payers bank

	if success {
		w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
			ID:    txn.ID,
			State: "DEBIT_SUCCESS",
		})

		w.r.Publish(ctx, "txn_stream", map[string]interface{}{
			"txn_id": txn.ID.String(),
			"event":  "txn.debit.success",
		})

	} else {
		w.q.UpdateTransactionState(ctx, repo.UpdateTransactionStateParams{
			ID:    txn.ID,
			State: "FAILED",
		})
	}

	return nil
}

func mockDebit() bool {
	panic("unimplemented")
}
