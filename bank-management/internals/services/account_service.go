package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	repo "github.com/swastiijain24/bank-management/internals/adapters/sqlc"
)

type Service interface {
	GetAccountById(ctx context.Context, id string) (repo.Account, error)
	CreateAccount(ctx context.Context, accountName string, accountPhone string) (repo.Account, error)
	Debit(ctx context.Context, accountId string, amount int64) (repo.Transaction, error)
	Credit(ctx context.Context, accountId string, amount int64) (repo.Transaction, error)
	GetTransactions(ctx context.Context, accountId string) ([]repo.Transaction, error)
	CheckBalance(ctx context.Context, accountId string) (int64, error)
}

type svc struct {
	repo repo.Querier
	db   *pgx.Conn
}

func NewService(repo repo.Querier, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) GetAccountById(ctx context.Context, id string) (repo.Account, error) {
	return s.repo.GetAccountByID(ctx, id)

}

func (s *svc) CreateAccount(ctx context.Context, accountName string, accountPhone string) (repo.Account, error) {

	if accountName == "" {
		return repo.Account{}, fmt.Errorf("Name not given")

	}

	if accountPhone == "" {
		return repo.Account{}, fmt.Errorf("Phone not given")

	}

	accountDetails := repo.CreateAccountParams{
		ID:    generateId("acc"),
		Name:  accountName,
		Phone: accountPhone,
	}

	account, err := s.repo.CreateAccount(ctx, accountDetails)
	if err != nil {
		return repo.Account{}, fmt.Errorf("Phone not given")
	}

	return account, err
}

func generateId(prefix string) string {
	return prefix + "-" + uuid.New().String()
}

func (s *svc) Debit(ctx context.Context, accountId string, amount int64) (repo.Transaction, error) {
	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, accountId)
	if err != nil {
		return repo.Transaction{}, err
	}

	if account.Balance < amount {
		return repo.Transaction{}, fmt.Errorf("insufficient balance")
	}

	if amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: account.Balance - amount,
		ID:      accountId,
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	txParamsObj := repo.CreateTransactionParams{
		ID:        generateId("tx"),
		AccountID: accountId,
		Amount:    amount,
		Type:      "DEBIT",
		Status:    "SUCCESS",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return transaction, err

}

func (s *svc) Credit(ctx context.Context, accountId string, amount int64) (repo.Transaction, error) {
	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, accountId)
	if err != nil {
		return repo.Transaction{}, err
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		ID:      accountId,
		Balance: account.Balance + amount,
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	txParamsObj := repo.CreateTransactionParams{
		ID:        generateId("tx"),
		AccountID: accountId,
		Amount:    amount,
		Type:      "CREDIT",
		Status:    "SUCCESS",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return transaction, err
}

func (s *svc) GetTransactions(ctx context.Context, accountId string) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, accountId)
}

func (s *svc) CheckBalance(ctx context.Context, accountId string) (int64, error) {
	return s.repo.CheckBalance(ctx, accountId)
}
