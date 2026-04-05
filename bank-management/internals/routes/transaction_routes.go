package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	"github.com/swastiijain24/bank-management/internals/middlewares"
)

func RegisterTransactionRoutes(r *gin.Engine, idempotencyMiddleware* middlewares.IdempotencyMiddleware , transactionHandler* handlers.TransactionHandler){

	transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/debit", idempotencyMiddleware.IdempotencyMiddleware(), transactionHandler.Debit)
		transactionRoutes.POST("/credit", idempotencyMiddleware.IdempotencyMiddleware(), transactionHandler.Credit)
		transactionRoutes.GET("/:id/transactions", transactionHandler.GetTransactions)

	}
	
}
