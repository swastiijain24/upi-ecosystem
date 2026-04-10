package services

import (
	"context"
	"fmt"

	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/utils"
)

type LedgerService interface {
	ReconcileAccount(ctx context.Context, accountId string, storedBalance int64) (bool, error)
}

type ledgersvc struct {
	repo           repo.Querier
}

func NewLedgerService(repo repo.Querier) LedgerService {
	return &ledgersvc{
		repo:           repo,
	}
}

func (s *ledgersvc) ReconcileAccount(ctx context.Context, accountId string, storedBalance int64) (bool, error) {

	calculated, err := s.repo.BalanceFromEntries(ctx, utils.StringtoUUID(accountId))
	if err != nil {
		return false, fmt.Errorf("failed to calculate balance: %w", err)
	}

	if storedBalance != calculated {
		return false, fmt.Errorf("balance mismatch: stored %d, calculated %d",
			storedBalance, calculated)
	}
	return true, nil
}
