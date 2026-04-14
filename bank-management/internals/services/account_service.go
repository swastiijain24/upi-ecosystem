package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type AccountService interface {
	GetAccountById(ctx context.Context, id string) (repo.Account, error)
	CreateAccount(ctx context.Context, Name string, Phone string, mpinHash string) (repo.CreateAccountRow, error)
	GetBalance(ctx context.Context, accountId string) (int64, error)
	DeleteAccount(ctx context.Context, accountId string) error
	CreateSettlementAccount(ctx context.Context) error
}

type accsvc struct {
	repo          repo.Querier
	db            *pgxpool.Pool
	ledgerService LedgerService
}

func NewAccountService(repo repo.Querier, db *pgxpool.Pool, ledgerService LedgerService) AccountService {
	return &accsvc{
		repo: repo,
		db:   db,
		ledgerService: ledgerService,
	}
}

func (s *accsvc) SetAccountMPIN(ctx context.Context, accountID string, hashedMpin string) error {
  
    return s.repo.SetMpinHash(ctx, repo.SetMpinHashParams{
        ID: utils.StringtoUUID(accountID),
        MpinHash:  hashedMpin,
    })
}

func (s *accsvc) GetAccountById(ctx context.Context, id string) (repo.Account, error) {
	return s.repo.GetAccountByID(ctx, utils.StringtoUUID(id))
}

func (s *accsvc) CreateAccount(ctx context.Context, Name string, Phone string, mpinHash string) (repo.CreateAccountRow, error) {

	if Name == "" {
		return repo.CreateAccountRow{}, fmt.Errorf("Name not given")

	}

	if Phone == "" {
		return repo.CreateAccountRow{}, fmt.Errorf("Phone not given")

	}

	if mpinHash == ""{
		return repo.CreateAccountRow{}, fmt.Errorf("mpinHash not given")
	}

	accountParams := repo.CreateAccountParams{
		Name:  Name,
		Phone: Phone,
	}

	account, err := s.repo.CreateAccount(ctx, accountParams)
	if err != nil {
		return repo.CreateAccountRow{}, fmt.Errorf("Error creating account")
	}

	return account, err
}

func (s *accsvc) GetBalance(ctx context.Context, accountId string) (int64, error) {

	stored, err := s.repo.GetBalance(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return 0, fmt.Errorf("failed to fetch balance: %w", err)
	}

	balanced, err := s.ledgerService.ReconcileAccount(ctx, accountId, stored)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch balance: %w", err)
	}
	if !balanced {
		return 0, fmt.Errorf("balance integrity check failed for account %s", accountId)
	}

	return stored, nil
}

func (s *accsvc) DeleteAccount(ctx context.Context, accountId string) error {
	err := s.repo.DeleteAccount(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return err
	}
	return nil
}

func (s *accsvc) CreateSettlementAccount(ctx context.Context) error {
	if err := s.repo.CreateSettlementAccount(ctx); err != nil {
		return err
	}
	return nil
}
