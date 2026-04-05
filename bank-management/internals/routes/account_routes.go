package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/handlers"
	apiAuth "github.com/swastiijain24/bank-management/internals/middlewares/api_auth"
)

func RegisterAccountRoutes(r *gin.Engine,  apiAuthMiddleware* apiAuth.APIMiddleware, accountHandler* handlers.AccountHandler){

	accountRoutes := r.Group("/accounts")
	{
		accountRoutes.POST("/",  apiAuthMiddleware.ApiAuthentication(), accountHandler.CreateAccount)
		accountRoutes.GET("/:id",  apiAuthMiddleware.ApiAuthentication(), accountHandler.GetAccountById)
		accountRoutes.GET("/:id/balance", apiAuthMiddleware.ApiAuthentication(), accountHandler.GetBalance)
		accountRoutes.DELETE("/:id",  apiAuthMiddleware.ApiAuthentication(), accountHandler.DeleteAccount)
	}
}

