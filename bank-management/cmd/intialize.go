package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/swastiijain24/bank-management/internals/handlers"
	repo "github.com/swastiijain24/bank-management/internals/repositories"
	"github.com/swastiijain24/bank-management/internals/routes"
	"github.com/swastiijain24/bank-management/internals/services"
)


func Initialize(r *gin.Engine,conn *pgx.Conn) {
	
	repository := repo.New(conn)

	accountService := services.NewAccountService(repository, conn)
	transactionService := services.NewTransactionService(repository, conn)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	routes.RegisterAccountRoutes(r, accountHandler)
	routes.RegisterTransactionRoutes(r, transactionHandler)

}