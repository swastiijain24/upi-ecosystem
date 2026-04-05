package services

import (
	"context"

	repo "github.com/swastiijain24/bank-management/internals/repositories"
)

type ApiKeyService interface {
	GetAPIKeyByKeyID(ctx context.Context, keyID string) (repo.ApiKey, error)
	IsValid(ctx context.Context, keyID string) (bool , error)
	UpdateAPIKeyLastUsed(ctx context.Context, keyID string) (error)
}

type ApiKeySvc struct {
	repo repo.Querier
}

func NewApiKeyService(repo repo.Querier) ApiKeyService {
	return &ApiKeySvc{
		repo: repo,
	}
}

func (s* ApiKeySvc) GetAPIKeyByKeyID(ctx context.Context, keyID string) (repo.ApiKey, error){
	return s.repo.GetAPIKeyByKeyID(ctx, keyID)
	
}

func (s* ApiKeySvc) IsValid(ctx context.Context, keyID string) (bool , error){
	return s.repo.IsValid(ctx, keyID)
	
} 


func (s* ApiKeySvc) UpdateAPIKeyLastUsed(ctx context.Context, keyID string) (error){
	return s.repo.UpdateAPIKeyLastUsed(ctx, keyID)
	
} 

