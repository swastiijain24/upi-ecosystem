package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/swastiijain24/bank-management/internals/dtos"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type TransactionService interface {
	Debit(ctx context.Context, debitDetails dtos.DebitCreditRequest) (repo.Transaction, error)
	Credit(ctx context.Context, creditDetails dtos.DebitCreditRequest) (repo.Transaction, error)
	GetTransactions(ctx context.Context, FromAccountId string) ([]repo.Transaction, error)
}

type txnsvc struct {
	repo repo.Querier
	db   *pgx.Conn
}

func NewTransactionService(repo repo.Querier, db *pgx.Conn) TransactionService {
	return &txnsvc{
		repo: repo,
		db:   db,
	}
}

func (s *txnsvc) Debit(ctx context.Context,debitDetails dtos.DebitCreditRequest) (repo.Transaction, error) {

	// _, err := s.repo.CheckIdempotencyKey(key)
	// if err == nil {
	// 	return repo.Transaction{}, fmt.Errorf("request already initiated")
	// }

	if debitDetails.Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, utils.StringtoUUID(debitDetails.FromAccountID))
	if err != nil {
		return repo.Transaction{}, err
	}

	if account.Balance < debitDetails.Amount {
		return repo.Transaction{}, fmt.Errorf("insufficient balance")
	}

	txParamsObj := repo.CreateTransactionParams{
		FromAccountID:       debitDetails.FromAccountID,
		ToAccountIdentifier: debitDetails.ToAccountId,
		Amount:              debitDetails.Amount,
		Status:              "PENDING",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}

	newBalance := account.Balance - debitDetails.Amount

	ledgerParams := repo.CreateLedgerEntryParams{
		AccountID:    utils.StringtoUUID(debitDetails.FromAccountID),
		Type:         "DEBIT",
		Amount:       debitDetails.Amount,
		BalanceAfter: newBalance,
		Description:  utils.ToPGText(debitDetails.Description),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParams); err != nil {
		return repo.Transaction{}, err
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      utils.StringtoUUID(debitDetails.FromAccountID),
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	transactionStatusParams := repo.UpdatePaymentStatusParams{
		ID: transaction.ID,
		Status: "SUCCESS",
	}

	if err := qtx.UpdatePaymentStatus(ctx, transactionStatusParams); err!=nil{
		return repo.Transaction{}, err
	}
	

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return transaction, nil

}

func (s *txnsvc) Credit(ctx context.Context, creditDetails dtos.DebitCreditRequest) (repo.Transaction, error) {

	// _, err := s.repo.CheckIdempotencyKey(key)
	// if err == nil {
	// 	return repo.Transaction{}, fmt.Errorf("request already initiated")
	// }

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, utils.StringtoUUID(creditDetails.ToAccountId))
	if err != nil {
		return repo.Transaction{}, err
	}

	txParamsObj := repo.CreateTransactionParams{
		FromAccountID:       creditDetails.FromAccountID,
		ToAccountIdentifier: creditDetails.ToAccountId,
		Amount:              creditDetails.Amount,
		Status:              "PENDING",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}


	newBalance := account.Balance + creditDetails.Amount

	ledgerParams := repo.CreateLedgerEntryParams{
		AccountID:    utils.StringtoUUID(creditDetails.ToAccountId),
		Type:         "CREDIT",
		Amount:       creditDetails.Amount,
		BalanceAfter: newBalance,
		Description:  utils.ToPGText(creditDetails.Description),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParams); err != nil {
		return repo.Transaction{}, err
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      utils.StringtoUUID(creditDetails.ToAccountId),
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}


	transactionStatusParams := repo.UpdatePaymentStatusParams{
		ID: transaction.ID,
		Status: "SUCCESS",
	}


	if err := qtx.UpdatePaymentStatus(ctx, transactionStatusParams); err!=nil{
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return transaction, nil
}

func (s *txnsvc) GetTransactions(ctx context.Context, accountID string) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, accountID)
}
