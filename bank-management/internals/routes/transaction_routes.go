package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
	"github.com/swastiijain24/bank-management/internals/middlewares/idempotency"
)

func RegisterTransactionRoutes(r *gin.Engine, apiAuthMiddleware* apiAuth.APIMiddleware, idempotencyMiddleware* idempotency.IdempotencyMiddleware , transactionHandler* handlers.TransactionHandler){

	transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/debit", apiAuthMiddleware.ApiAuthentication(), idempotencyMiddleware.IdempotencyCheck(), transactionHandler.Debit)
		transactionRoutes.POST("/credit", apiAuthMiddleware.ApiAuthentication(), idempotencyMiddleware.IdempotencyCheck(), transactionHandler.Credit)
		transactionRoutes.GET("/:id/transactions", transactionHandler.GetTransactions)

	}
	
}
