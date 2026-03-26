package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
)

func RegisterAccountRoutes(r *gin.Engine, accountHandler* handlers.AccountHandler){

	accountRoutes := r.Group("/accounts")
	{
		accountRoutes.POST("/", accountHandler.CreateAccount)
		accountRoutes.GET("/:id", accountHandler.GetAccountById)
		accountRoutes.GET("/:id/balance", accountHandler.GetBalance)
		accountRoutes.DELETE("/:id", accountHandler.DeleteAccount)
	}
}

