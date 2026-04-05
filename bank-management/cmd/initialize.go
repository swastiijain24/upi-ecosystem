package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
	"github.com/swastiijain24/bank-management/internals/middlewares/idempotency"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/routes"
	"github.com/swastiijain24/bank-management/internals/services"
)

func Initialize(r *gin.Engine, pool *pgxpool.Pool, ctx context.Context) {

	repository := repo.New(pool)

	accountService := services.NewAccountService(repository, pool)
	transactionService := services.NewTransactionService(repository, pool)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	redisClient := idempotency.NewRedis()
	redisStore := idempotency.NewRedisStore(redisClient, 24*time.Hour)
	idempotencyMiddleware := idempotency.NewIdempotencyMiddleware(*redisStore)

	keyAuth := apiAuth.NewKeyAuth()
	apiKeyService := services.NewApiKeyService(repository)
	apiAuthMiddleware := apiAuth.NewApiAuthMiddleware(keyAuth, apiKeyService)

	routes.RegisterAccountRoutes(r, apiAuthMiddleware, accountHandler)
	routes.RegisterTransactionRoutes(r, apiAuthMiddleware, idempotencyMiddleware, transactionHandler)

	if err := accountService.CreateSettlementAccount(ctx); err != nil {
		log.Fatal("Failed to create settlement account:", err)
	}
}
