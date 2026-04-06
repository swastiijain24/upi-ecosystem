package services

import (
	"context"
	"fmt"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
)

type LedgerService interface {
	ReconcileAccount(ctx context.Context, accountID string) (bool, error)
}

type ledgersvc struct {
	repo repo.Querier
	accountService AccountService
}

func NewLedgerService(repo repo.Querier, accountService AccountService) LedgerService {
	return &ledgersvc{
		repo: repo,
		accountService: accountService,
	}
}

func (s *ledgersvc) ReconcileAccount(ctx context.Context, accountID string) (bool, error) {
	account, err := s.accountService.GetAccountById(ctx, accountID)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}

	calculated, err := s.repo.BalanceFromEntries(ctx, account.ID)
	if err != nil {
		return false, fmt.Errorf("failed to calculate balance: %w", err)
	}

	stored, err := s.accountService.GetBalance(ctx, accountID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch balance: %w", err)
	}

	if stored != calculated {
		return false, fmt.Errorf("balance mismatch: stored %d, calculated %d",
			stored, calculated)
	}
	return true, nil
}
