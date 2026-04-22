package services

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
	"golang.org/x/crypto/bcrypt"
)

type AccountService interface {
	GetAccountById(ctx context.Context, id string) (repo.Account, error)
	CreateAccount(ctx context.Context, Name string, Phone string, mpinHash string) (repo.CreateAccountRow, error)
	GetBalance(ctx context.Context, accountId string, mpinEn string) (int64, error)
	DeleteAccount(ctx context.Context, accountId string) error
	CreateSettlementAccount(ctx context.Context) (string, error)
	DiscoverAccounts(ctx context.Context, phone string) ([]string, error)
	SetMpin(ctx context.Context, accountId string, mpinEn string) error
	ChangeMpin(ctx context.Context, accountId string, oldMpinEn string, newMpinEn string) error
}

type accsvc struct {
	repo          repo.Querier
	db            *pgxpool.Pool
	ledgerService LedgerService
}

func NewAccountService(repo repo.Querier, db *pgxpool.Pool, ledgerService LedgerService) AccountService {
	return &accsvc{
		repo:          repo,
		db:            db,
		ledgerService: ledgerService,
	}
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
	if mpinHash == "" {
		return repo.CreateAccountRow{}, fmt.Errorf("mpinHash not given")
	}

	decryptedPin, err := utils.DecryptAES(mpinHash, []byte(os.Getenv("MPIN_ENCRYPTION_KEY")))
	if err != nil {
		return repo.CreateAccountRow{}, fmt.Errorf("decryption error: %w", err)
	}

	hashedMpin, err := bcrypt.GenerateFromPassword([]byte(decryptedPin), bcrypt.DefaultCost)
	if err != nil {
		return repo.CreateAccountRow{}, fmt.Errorf("failed to hash mpin: %w", err)
	}

	accountParams := repo.CreateAccountParams{
		Name:     Name,
		Phone:    Phone,
		MpinHash: utils.ToPGText(string(hashedMpin)),
	}

	account, err := s.repo.CreateAccount(ctx, accountParams)
	if err != nil {
		return repo.CreateAccountRow{}, fmt.Errorf("Error creating account")
	}

	return account, err
}

func (s *accsvc) GetBalance(ctx context.Context, accountId string, mpinEn string) (int64, error) {

	decryptedPin, err := utils.DecryptAES(mpinEn, []byte(os.Getenv("MPIN_ENCRYPTION_KEY")))
	if err != nil {
		return 0, fmt.Errorf("decryption error: %w", err)
	}

	mpinHashstored, err := s.repo.GetMpinHash(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return 0, fmt.Errorf("error fetching pin")
	}

	err = bcrypt.CompareHashAndPassword([]byte(mpinHashstored.String), []byte(decryptedPin))
	if err != nil {
		return 0, fmt.Errorf("invalid mpin")
	}

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

func (s *accsvc) DiscoverAccounts(ctx context.Context, phone string) ([]string, error) {
	accountIds, err := s.repo.DiscoverAccountsByPhone(ctx, phone)
	if err != nil {
		return []string{}, err
	}
	var result []string
	for _, accountId := range accountIds {
		result = append(result, accountId.String())
	}
	return result, err
}

func (s *accsvc) SetMpin(ctx context.Context, accountId string, mpinEn string) error {
	decryptedPin, err := utils.DecryptAES(mpinEn, []byte(os.Getenv("MPIN_ENCRYPTION_KEY")))
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	hashedMpin, err := bcrypt.GenerateFromPassword([]byte(decryptedPin), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash mpin: %w", err)
	}

	err = s.repo.UpdateMpin(ctx, repo.UpdateMpinParams{
		ID:       utils.StringtoUUID(accountId),
		MpinHash: utils.ToPGText(string(hashedMpin)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *accsvc) ChangeMpin(ctx context.Context, accountId string, oldMpinEn string, newMpinEn string) error {
	decryptedPin, err := utils.DecryptAES(oldMpinEn, []byte(os.Getenv("MPIN_ENCRYPTION_KEY")))
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	mpinHashstored, err := s.repo.GetMpinHash(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return fmt.Errorf("error fetching pin")
	}

	err = bcrypt.CompareHashAndPassword([]byte(mpinHashstored.String), []byte(decryptedPin))
	if err != nil {
		return fmt.Errorf("invalid mpin")
	}

	return s.SetMpin(ctx, accountId, newMpinEn)

}

func (s *accsvc) DeleteAccount(ctx context.Context, accountId string) error {
	err := s.repo.DeleteAccount(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return err
	}
	return nil
}

func (s *accsvc) CreateSettlementAccount(ctx context.Context) (string, error) {
	accountId, err := s.repo.CreateSettlementAccount(ctx)
	if err != nil {
		return "", err
	}
	return accountId.String(), nil
}
