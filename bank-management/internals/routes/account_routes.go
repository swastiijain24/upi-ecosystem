package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
)

func RegisterAccountRoutes(r *gin.Engine,  apiAuthMiddleware *apiAuth.APIMiddleware, accountHandler *handlers.AccountHandler){

	accountRoutes := r.Group("/accounts")
	{
		accountRoutes.POST("/",   accountHandler.CreateAccount)
		accountRoutes.GET("/:id",  accountHandler.GetAccountById)
		accountRoutes.POST("/balance/:id", apiAuthMiddleware.ApiAuthentication(), accountHandler.GetBalance)
		accountRoutes.DELETE("/:id", apiAuthMiddleware.ApiAuthentication(),  accountHandler.DeleteAccount)
		accountRoutes.POST("/discover", apiAuthMiddleware.ApiAuthentication(), accountHandler.DiscoverAccounts)
		accountRoutes.POST("/mpin/:id", apiAuthMiddleware.ApiAuthentication(), accountHandler.SetMpin)
		accountRoutes.PUT("/mpin/:id", apiAuthMiddleware.ApiAuthentication(), accountHandler.ChangeMpin)
	}
}

