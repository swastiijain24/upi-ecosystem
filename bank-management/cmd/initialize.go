package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swastiijain24/bank-management/internals/audit"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
	"github.com/swastiijain24/bank-management/internals/middlewares/idempotency"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/routes"
	"github.com/swastiijain24/bank-management/internals/services"
)

func Initialize(r *gin.Engine, pool *pgxpool.Pool) {

	repository := repo.New(pool)
	ctx := context.Background()

	ledgerService := services.NewLedgerService(repository)
	accountService := services.NewAccountService(repository, pool, ledgerService)
	settlementAccountId, err := accountService.CreateSettlementAccount(ctx)
	if err != nil {
		log.Print("Failed to create settlement account:", err)
	}
	transactionService := services.NewTransactionService(repository, pool, settlementAccountId)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	redisClient := idempotency.NewRedis()
	redisStore := idempotency.NewRedisStore(redisClient, 24*time.Hour)
	idempotencyMiddleware := idempotency.NewIdempotencyMiddleware(*redisStore)

	APIKeyHasher := apiAuth.NewAPIKeyHasher()
	APIKeyGenerator := apiAuth.NewAPIKeyGenerator()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	auditLogger := audit.NewLogger(logger)

	apiKeyService := services.NewApiKeyService(repository)

	apiAuthMiddleware := apiAuth.NewApiAuthMiddleware(APIKeyGenerator, APIKeyHasher, auditLogger, apiKeyService)

	routes.RegisterAccountRoutes(r, apiAuthMiddleware, accountHandler)
	routes.RegisterTransactionRoutes(r, apiAuthMiddleware, idempotencyMiddleware, transactionHandler)

}
