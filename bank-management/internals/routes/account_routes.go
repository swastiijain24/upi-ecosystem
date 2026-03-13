package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
)

func RegisterRoutes(r *gin.Engine, accountHandler* handlers.Handler){

	accountRoutes := r.Group("/accounts")
	{
		accountRoutes.POST("/", accountHandler.CreateAccount)
		accountRoutes.GET("/:id", accountHandler.GetAccountById)
		accountRoutes.POST("/debit", accountHandler.Debit)
		accountRoutes.POST("/credit", accountHandler.Credit)
		accountRoutes.GET("/:id/transactions", accountHandler.GetTransactions)
		accountRoutes.GET("/:id/balance", accountHandler.CheckBalance)

	}
	
}

