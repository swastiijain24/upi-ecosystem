package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
	"github.com/swastiijain24/bank-management/internals/middlewares/idempotency"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/routes"
	"github.com/swastiijain24/bank-management/internals/services"
)

func Initialize(r *gin.Engine, conn *pgx.Conn) {

	repository := repo.New(conn)

	accountService := services.NewAccountService(repository, conn)
	transactionService := services.NewTransactionService(repository, conn)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	redisClient := idempotency.NewRedis()
	redisStore := idempotency.NewRedisStore(redisClient, 24*time.Hour)
	idempotencyMiddleware := idempotency.NewIdempotencyMiddleware(*redisStore) 

	keyAuth := apiAuth.NewKeyAuth()
	apiKeyService := services.NewApiKeyService(repository)
	apiAuthMiddleware := apiAuth.NewApiAuthMiddleware(keyAuth, apiKeyService)

	routes.RegisterAccountRoutes(r, accountHandler)
	routes.RegisterTransactionRoutes(r, apiAuthMiddleware, idempotencyMiddleware, transactionHandler)

}
