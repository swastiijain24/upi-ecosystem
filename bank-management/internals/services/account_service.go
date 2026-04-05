package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type AccountService interface {
	GetAccountById(ctx context.Context, id string) (repo.Account, error)
	CreateAccount(ctx context.Context, Name string, Phone string) (repo.CreateAccountRow, error)
	GetBalance(ctx context.Context, accountId string) (int64, error)
	DeleteAccount(ctx context.Context, accountId string) error
	CreateSettlementAccount(ctx context.Context) error 
}

type accsvc struct {
	repo repo.Querier
	db   *pgx.Conn
}

func NewAccountService(repo repo.Querier, db *pgx.Conn) AccountService {
	return &accsvc{
		repo: repo,
		db:   db,
	}
}

func (s *accsvc) GetAccountById(ctx context.Context, id string) (repo.Account, error) {
	return s.repo.GetAccountByID(ctx, utils.StringtoUUID(id))
}

func (s *accsvc) CreateAccount(ctx context.Context, Name string, Phone string) (repo.CreateAccountRow, error) {

	if Name == "" {
		return repo.CreateAccountRow{}, fmt.Errorf("Name not given")

	}

	if Phone == "" {
		return repo.CreateAccountRow{}, fmt.Errorf("Phone not given")

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
	return s.repo.GetBalance(ctx, utils.StringtoUUID(accountId))
}

func (s* accsvc) DeleteAccount(ctx context.Context, accountId string) error{
	err := s.repo.DeleteAccount(ctx, utils.StringtoUUID(accountId)) 
	if err !=nil{
		return err
	}
	return nil
}

func (s*accsvc) CreateSettlementAccount(ctx context.Context) error {
	if err := s.repo.CreateSettlementAccount(ctx); err!=nil{
		return err
	}
	return nil
}