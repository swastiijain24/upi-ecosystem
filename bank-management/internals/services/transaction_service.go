package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type TransactionService interface {
	Debit(ctx context.Context, FromAccountID string, ToAccountId string, Amount string, Description string) (repo.Transaction, error)
	Credit(ctx context.Context, FromAccountID string, ToAccountId string, Amount string, Description string) (repo.Transaction, error)
	GetTransactions(ctx context.Context, FromAccountId string) ([]repo.Transaction, error)
}

type txnsvc struct {
	repo repo.Querier
	db   *pgxpool.Pool
}

func NewTransactionService(repo repo.Querier, db *pgxpool.Pool) TransactionService {
	return &txnsvc{
		repo: repo,
		db:   db,
	}
}

func (s *txnsvc) Debit(ctx context.Context, FromAccountID string, ToAccountId string, AmountinString string, Description string) (repo.Transaction, error) {

	Amount, err := utils.RupeesToPaise(AmountinString)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error converting amount")
	}

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, utils.StringtoUUID(FromAccountID))
	if err != nil {
		return repo.Transaction{}, err
	}

	if account.Balance < Amount {
		return repo.Transaction{}, fmt.Errorf("insufficient balance")
	}

	settlementAccount, err := qtx.GetSettlementAccountForUpdate(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}

	txParamsObj := repo.CreateTransactionParams{
		FromAccountID:       FromAccountID,
		ToAccountIdentifier: ToAccountId,
		Amount:              Amount,
		Status:              "PENDING",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}

	newBalance := account.Balance - Amount
	newSettlementAccountBalance := settlementAccount.Balance + Amount

	ledgerParams := repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(FromAccountID),
		Type:          "DEBIT",
		Debit:         Amount,
		Credit:        0,
		BalanceAfter:  newBalance,
		Description:   utils.ToPGText(Description),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParams); err != nil {
		return repo.Transaction{}, err
	}

	ledgerParamsSettlementAccount := repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     settlementAccount.ID,
		Type:          "CREDIT",
		Credit:        Amount,
		Debit:         0,
		BalanceAfter:  newSettlementAccountBalance,
		Description:   utils.ToPGText("settlement account"),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParamsSettlementAccount); err != nil {
		return repo.Transaction{}, err
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      utils.StringtoUUID(FromAccountID),
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	settlementUpdatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newSettlementAccountBalance,
		ID:      settlementAccount.ID,
	}

	if err := qtx.UpdateAccountBalance(ctx, settlementUpdatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	transactionStatusParams := repo.UpdatePaymentStatusParams{
		ID:     transaction.ID,
		Status: "SUCCESS",
	}

	if err := qtx.UpdatePaymentStatus(ctx, transactionStatusParams); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, err := qtx.GetTransactionById(ctx, transaction.ID)
	if err != nil {
		return repo.Transaction{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return repo.Transaction{}, err
	}
	return finalTransaction, nil

}

func (s *txnsvc) Credit(ctx context.Context, FromAccountID string, ToAccountId string, AmountinString string, Description string) (repo.Transaction, error) {

	Amount, err := utils.RupeesToPaise(AmountinString)
	if err != nil {
		return repo.Transaction{}, fmt.Errorf("error converting amount")
	}

	if Amount <= 0 {
		return repo.Transaction{}, fmt.Errorf("invalid amount")
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}
	defer dbTx.Rollback(ctx)

	qtx := repo.New(dbTx)

	account, err := qtx.GetAccountForUpdate(ctx, utils.StringtoUUID(ToAccountId))
	if err != nil {
		return repo.Transaction{}, err
	}

	settlementAccount, err := qtx.GetSettlementAccountForUpdate(ctx)
	if err != nil {
		return repo.Transaction{}, err
	}

	txParamsObj := repo.CreateTransactionParams{
		FromAccountID:       FromAccountID,
		ToAccountIdentifier: ToAccountId,
		Amount:              Amount,
		Status:              "PENDING",
	}

	transaction, err := qtx.CreateTransaction(ctx, txParamsObj)
	if err != nil {
		return repo.Transaction{}, err
	}

	newBalance := account.Balance + Amount
	newSettlementAccountBalance := settlementAccount.Balance - Amount

	ledgerParams := repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     utils.StringtoUUID(ToAccountId),
		Type:          "CREDIT",
		Credit:        Amount,
		Debit:         0,
		BalanceAfter:  newBalance,
		Description:   utils.ToPGText(Description),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParams); err != nil {
		return repo.Transaction{}, err
	}

	ledgerParamsSettlementAccount := repo.CreateLedgerEntryParams{
		TransactionID: transaction.ID,
		AccountID:     settlementAccount.ID,
		Type:          "DEBIT",
		Debit:         Amount,
		Credit:        0,
		BalanceAfter:  newSettlementAccountBalance,
		Description:   utils.ToPGText("settlement account"),
	}

	if err := qtx.CreateLedgerEntry(ctx, ledgerParamsSettlementAccount); err != nil {
		return repo.Transaction{}, err
	}

	updatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newBalance,
		ID:      utils.StringtoUUID(ToAccountId),
	}

	if err := qtx.UpdateAccountBalance(ctx, updatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	settlementUpdatedParamsObj := repo.UpdateAccountBalanceParams{
		Balance: newSettlementAccountBalance,
		ID:      settlementAccount.ID,
	}

	if err := qtx.UpdateAccountBalance(ctx, settlementUpdatedParamsObj); err != nil {
		return repo.Transaction{}, err
	}

	transactionStatusParams := repo.UpdatePaymentStatusParams{
		ID:     transaction.ID,
		Status: "SUCCESS",
	}

	if err := qtx.UpdatePaymentStatus(ctx, transactionStatusParams); err != nil {
		return repo.Transaction{}, err
	}

	finalTransaction, err := qtx.GetTransactionById(ctx, transaction.ID)
	if err != nil {
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
