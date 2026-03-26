package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
)

func RegisterTransactionRoutes(r *gin.Engine, transactionHandler* handlers.TransactionHandler){

	transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/debit", transactionHandler.Debit)
		transactionRoutes.POST("/credit", transactionHandler.Credit)
		transactionRoutes.GET("/:id/transactions", transactionHandler.GetTransactions)

	}
	
}

