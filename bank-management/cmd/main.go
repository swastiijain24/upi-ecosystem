package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	
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
	r.Use(gin.Logger())
	Initialize(r, conn)

	port := os.Getenv("PORT")

	if port == ""{
		port = "8080"
	}
	

	log.Printf("server running")
	r.Run(":" + port)

}