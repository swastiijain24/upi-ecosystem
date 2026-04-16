package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type TransactionService interface {
	Debit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, mpinHash string, externalId string) (repo.Transaction, error)
	Credit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, externalId string) (repo.Transaction, error)
	GetTransactions(ctx context.Context, FromAccountId string) ([]repo.Transaction, error)
	GetStatusByExternalId(ctx context.Context, externalId string) (string, string, error)
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
	if mpinHash != mpinHashstored {
		return repo.Transaction{}, fmt.Errorf("invalid mpin")
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	transaction, err := qtx.CreateTransaction(ctx, repo.CreateTransactionParams{
		FromAccountID:       FromAccountID,
		ToAccountIdentifier: ToAccountId,
		Amount:              Amount,
		Status:              "PENDING",
		ExternalID:          externalId,
	})
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := qtx.UpdatePaymentStatus(ctx, repo.UpdatePaymentStatusParams{
		ID:     transaction.ID,
		Status: "SUCCESS",
	}); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, err := qtx.GetTransactionById(ctx, transaction.ID)
	if err != nil {
		return repo.Transaction{}, err
	}

	newUserBalance, err := qtx.UpdateUserBalanceDebit(ctx, repo.UpdateUserBalanceDebitParams{
		Balance: Amount,
		ID:      utils.StringtoUUID(FromAccountID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return repo.Transaction{}, fmt.Errorf("insufficient balance or account not found")
		}
		return repo.Transaction{}, fmt.Errorf("Database error: %v", err)
	}

	//ledger entry for user account
	err = qtx.CreateLedgerEntry(ctx, repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(FromAccountID),
		Type:          "DEBIT",
		Debit:         Amount,
		Credit:        0,
		BalanceAfter:  newUserBalance,
		Description:   utils.ToPGText(Description),
	})
	if err != nil {
		return repo.Transaction{}, err
	}

	newSettlementAccountBalance, err := qtx.UpdateSettlementBalanceAtomic(ctx, Amount)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error updating settlement account balance")
	}

	//ledger entry for settlement account
	err = qtx.CreateLedgerEntry(ctx, repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(s.settlementAccountId),
		Type:          "CREDIT",
		Credit:        Amount,
		Debit:         0,
		BalanceAfter:  newSettlementAccountBalance,
		Description:   utils.ToPGText("settlement account"),
	})
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return finalTransaction, nil

}

func (s *txnsvc) Credit(ctx context.Context, FromAccountID string, ToAccountId string, Amount int64, Description string, externalId string) (repo.Transaction, error) {

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	transaction, err := qtx.CreateTransaction(ctx, repo.CreateTransactionParams{
		FromAccountID:       FromAccountID,
		ToAccountIdentifier: ToAccountId,
		Amount:              Amount,
		Status:              "PENDING",
		ExternalID:          externalId,
	})
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := qtx.UpdatePaymentStatus(ctx, repo.UpdatePaymentStatusParams{
		ID:     transaction.ID,
		Status: "SUCCESS",
	}); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, err := qtx.GetTransactionById(ctx, transaction.ID)
	if err != nil {
		return repo.Transaction{}, err
	}

	newUserBalance, err := qtx.UpdateAccountBalanceCredit(ctx, repo.UpdateAccountBalanceCreditParams{
		Balance: Amount,
		ID:      utils.StringtoUUID(ToAccountId),
	})

	//ledger entry for user account
	if err := qtx.CreateLedgerEntry(ctx, repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(ToAccountId),
		Type:          "CREDIT",
		Credit:        Amount,
		Debit:         0,
		BalanceAfter:  newUserBalance,
		Description:   utils.ToPGText(Description),
	}); err != nil {
		return repo.Transaction{}, err
	}

	//will send negative of amount
	newSettlementAccountBalance, err := qtx.UpdateSettlementBalanceAtomic(ctx, -Amount)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error updating settlement account balance")
	}

	if err := qtx.CreateLedgerEntry(ctx, repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(s.settlementAccountId),
		Type:          "DEBIT",
		Debit:         Amount,
		Credit:        0,
		BalanceAfter:  newSettlementAccountBalance,
		Description:   utils.ToPGText("settlement account"),
	}); err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return finalTransaction, nil
}

func (s *txnsvc) GetTransactions(ctx context.Context, accountID string) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, accountID)
}

func (s *txnsvc) GetStatusByExternalId(ctx context.Context, externalId string) (string, string, error) {
	result, err := s.repo.GetTransactionStatusByExternalId(ctx, externalId)
	if err != nil {
		return "", "", err
	}
	return result.ID.String(), result.Status, err
}
