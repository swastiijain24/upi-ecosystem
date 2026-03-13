package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	repo "github.com/swastiijain24/bank-management/internals/adapters/sqlc"
	"github.com/swastiijain24/bank-management/internals/handlers"
	"github.com/swastiijain24/bank-management/internals/routes"
	"github.com/swastiijain24/bank-management/internals/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()
	dsn := os.Getenv("GOOSE_DBSTRING")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		panic(err) 
	}
	defer conn.Close(ctx)

	log.Printf("connected to database")

	r := gin.New()

	accountService := services.NewService(repo.New(conn), conn)
	accountHandler := handlers.NewHandler(accountService)

	r.Use(gin.Logger())
	routes.RegisterRoutes(r, accountHandler)

	port := os.Getenv("PORT")

	if port == ""{
		port = "8080"
	}
	

	log.Printf("server running")
	r.Run(":" + port)

}