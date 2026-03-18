package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/swastiijain24/npci-switch/internals/adapters/sqlc"
	"github.com/swastiijain24/npci-switch/internals/helpers"
	"github.com/swastiijain24/npci-switch/internals/redis_stream"
)

type Service interface {
	InitiatePayment(ctx context.Context, payerVpa string, payerBank string, payeeVpa string, payeeBank string, amount int64, referenceID pgtype.Text) (error)
}

type svc struct {
	q *repo.Queries
	r *redis_stream.EventPublisher
}

func NewService(q *repo.Queries, r *redis_stream.EventPublisher) Service {
	return &svc{
		q: q,
		r: r,
	}
}

func (s *svc) InitiatePayment(ctx context.Context, payerVpa string, payerBank string, payeeVpa string, payeeBank string, amount int64, referenceID pgtype.Text) ( error) {

	txn, err := s.q.GetTransactionByReference(ctx, referenceID)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) { //we only have to cehck if the row with this referid trnx exists or not if not then only create new , not for all errors
			txn, err = s.q.CreateTransaction(ctx, repo.CreateTransactionParams{
				ID:          helpers.ToPgUUID(uuid.New()),
				PayerVpa:    payerVpa,
				PayerBank:   payerBank,
				PayeeVpa:    payeeVpa,
				PayeeBank:   payeeBank,
				Amount:      amount,
				State:       "INITIATED",
				ReferenceID: referenceID,
			})
			if err != nil {
				log.Println("failed to create txn:", err)
				return err
			}

		} else {

			log.Println("db error:", err)
			return err
		}
	}

	err = s.r.Publish(ctx, "txn_stream", map[string]interface{}{
		"txn_id": txn.ID.String(),
		"event":  "txn.initiated",
	})
	if err != nil {
		log.Println("event publish failed, retry later")
	}

	return nil 

}
