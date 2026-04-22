package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
)

func RegisterAccountRoutes(r *gin.Engine,  apiAuthMiddleware* apiAuth.APIMiddleware, accountHandler* handlers.AccountHandler){

	accountRoutes := r.Group("/accounts")
	{
		accountRoutes.POST("/",  accountHandler.CreateAccount)
		accountRoutes.GET("/:id",  accountHandler.GetAccountById)
		accountRoutes.GET("/balance/:id", accountHandler.GetBalance)
		accountRoutes.DELETE("/:id",  accountHandler.DeleteAccount)
		accountRoutes.GET("/discover", accountHandler.DiscoverAccounts)
		accountRoutes.POST("/mpin/:id",  accountHandler.SetMpin)
		accountRoutes.PUT("/mpin/:id",  accountHandler.ChangeMpin)
	}
}

