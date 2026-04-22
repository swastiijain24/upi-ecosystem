package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
	"github.com/swastiijain24/bank-management/internals/middlewares/idempotency"
)

func RegisterTransactionRoutes(r *gin.Engine, apiAuthMiddleware* apiAuth.APIMiddleware, idempotencyMiddleware* idempotency.IdempotencyMiddleware , transactionHandler* handlers.TransactionHandler){

	transactionRoutes := r.Group("/transactions")
	transactionRoutes.Use(apiAuthMiddleware.ApiAuthentication())
	{
		transactionRoutes.POST("/debit", idempotencyMiddleware.IdempotencyCheck(), transactionHandler.Debit)
		transactionRoutes.POST("/credit", idempotencyMiddleware.IdempotencyCheck(), transactionHandler.Credit)
		transactionRoutes.POST("/refund", idempotencyMiddleware.IdempotencyCheck(), transactionHandler.Refund)
		transactionRoutes.GET("/account/:id", transactionHandler.GetTransactions)
		transactionRoutes.GET("/status/:external_id", transactionHandler.GetStatusOfTransaction)
	}
}
