package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type TransactionService interface {
	Debit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, mpinHash string, externalId string) (repo.Transaction, error)
	Credit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, externalId string) (repo.Transaction, error)
	Refund(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, externalId string) (repo.Transaction, error)
	GetTransactions(ctx context.Context, FromAccountId string) ([]repo.Transaction, error)
	GetStatus(ctx context.Context, externalId string, transactionType string) (string, string, error)
}

type txnsvc struct {
	repo                repo.Querier
	db                  *pgxpool.Pool
	settlementAccountId string
}

func NewTransactionService(repo repo.Querier, db *pgxpool.Pool, settlementAccountId string) TransactionService {
	return &txnsvc{
		repo:                repo,
		db:                  db,
		settlementAccountId: settlementAccountId,
	}
}

func (s *txnsvc) Debit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, mpinHash string, externalId string) (repo.Transaction, error) {

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	mpinHashstored, err := s.repo.GetMpinHash(ctx, utils.StringtoUUID(FromAccountID))
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error fetching pin")
	}
	if utils.ToPGText(mpinHash) != mpinHashstored {
		return repo.Transaction{}, fmt.Errorf("invalid mpin")
	}

	existingTransaction, err := s.repo.GetTransactionForIdempotency(ctx, repo.GetTransactionForIdempotencyParams{
		ExternalID: externalId,
		Type:       "DEBIT",
	})
	if err == nil {
		return existingTransaction, nil
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	transaction, err := s.createTransaction(ctx, qtx, FromAccountID, ToAccountId, Amount, externalId, "DEBIT")
	if err != nil {
		return repo.Transaction{}, err
	}

	newUserBalance, err := qtx.UpdateUserBalanceDebit(ctx, repo.UpdateUserBalanceDebitParams{
		Balance: Amount,
		ID:      utils.StringtoUUID(FromAccountID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Transaction{}, fmt.Errorf("insufficient balance or account not found")
		}
		return repo.Transaction{}, fmt.Errorf("Database error: %v", err)
	}

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, FromAccountID, "DEBIT", Amount, 0, newUserBalance, Description)
	if err != nil {
		return repo.Transaction{}, err
	}

	newSettlementAccountBalance, err := qtx.UpdateSettlementBalanceAtomic(ctx, Amount)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error updating settlement account balance")
	}

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, s.settlementAccountId, "CREDIT", 0, Amount, newSettlementAccountBalance, "settlement account")
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, _ := s.repo.GetTransactionById(ctx, transaction.ID)
	return finalTransaction, nil

}

func (s *txnsvc) Credit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, externalId string) (repo.Transaction, error) {

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	existingTransaction, err := s.repo.GetTransactionForIdempotency(ctx, repo.GetTransactionForIdempotencyParams{
		ExternalID: externalId,
		Type:       "CREDIT",
	})
	if err == nil {
		return existingTransaction, nil
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	transaction, err := s.createTransaction(ctx, qtx, FromAccountID, ToAccountId, Amount, externalId, "CREDIT")
	if err != nil {
		return repo.Transaction{}, err
	}

	newUserBalance, err := qtx.UpdateAccountBalanceCredit(ctx, repo.UpdateAccountBalanceCreditParams{
		Balance: Amount,
		ID:      utils.StringtoUUID(ToAccountId),
	})

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, ToAccountId, "CREDIT", 0, Amount, newUserBalance, Description)
	if err != nil {
		return repo.Transaction{}, err
	}

	//will send negative of amount
	newSettlementAccountBalance, err := qtx.UpdateSettlementBalanceAtomic(ctx, -Amount)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error updating settlement account balance")
	}

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, s.settlementAccountId, "DEBIT", Amount, 0, newSettlementAccountBalance, "settlement account")
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, _ := s.repo.GetTransactionById(ctx, transaction.ID)
	return finalTransaction, nil
}

func (s *txnsvc) Refund(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, externalId string) (repo.Transaction, error) {

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	existingTransaction, err := s.repo.GetTransactionForIdempotency(ctx, repo.GetTransactionForIdempotencyParams{
		ExternalID: externalId,
		Type:       "REFUND",
	})
	if err == nil {
		return existingTransaction, nil
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	transaction, err := s.createTransaction(ctx, qtx, FromAccountID, ToAccountId, Amount , externalId, "REFUND")
	if err != nil {
		return repo.Transaction{}, err
	}

	newUserBalance, err := qtx.UpdateAccountBalanceCredit(ctx, repo.UpdateAccountBalanceCreditParams{
		Balance: Amount,
		ID:      utils.StringtoUUID(ToAccountId),
	})

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, ToAccountId, "CREDIT", 0, Amount, newUserBalance, "")
	if err != nil {
		return repo.Transaction{}, err
	}

	//will send negative of amount
	newSettlementAccountBalance, err := qtx.UpdateSettlementBalanceAtomic(ctx, -Amount)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error updating settlement account balance")
	}

	err = s.createLedgerEntry(ctx, qtx, transaction.ID, s.settlementAccountId, "DEBIT", Amount, 0, newSettlementAccountBalance, "settlement account")
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, _ := s.repo.GetTransactionById(ctx, transaction.ID)
	return finalTransaction, nil
}

func (s *txnsvc) GetTransactions(ctx context.Context, accountID string) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, accountID)
}

func (s *txnsvc) GetStatus(ctx context.Context, externalId string, transactionType string) (string, string, error) {
	result, err := s.repo.GetTransactionStatus(ctx, repo.GetTransactionStatusParams{
		ExternalID: externalId,
		Type:       transactionType,
	})
	if err != nil {
		return "", "", err
	}
	return result.ID.String(), result.Status, err
}

func (s *txnsvc) createTransaction(ctx context.Context, qtx *repo.Queries, FromAccountID string, ToAccountId string, Amount int64, externalId string, txnType string) (repo.Transaction, error) {
	return qtx.CreateTransaction(ctx, repo.CreateTransactionParams{
		FromAccountID:       FromAccountID,
		ToAccountIdentifier: ToAccountId,
		Amount:              Amount,
		Status:              "SUCCESS",
		ExternalID:          externalId,
		Type:                txnType,
	})
}

func (s *txnsvc) createLedgerEntry(ctx context.Context, qtx *repo.Queries, transactionId pgtype.UUID, accountId string, txnType string, debit int64, credit int64, balance int64, description string) error {
	return qtx.CreateLedgerEntry(ctx, repo.CreateLedgerEntryParams{
		TransactionID: transactionId,
		AccountID:     utils.StringtoUUID(accountId),
		Type:          txnType,
		Debit:         debit,
		Credit:        credit,
		BalanceAfter:  balance,
		Description:   utils.ToPGText(description),
	})
}
