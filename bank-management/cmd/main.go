package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()
	dsn := os.Getenv("GOOSE_DBSTRING")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	log.Printf("connected to database")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	log.Printf("initializing")
	Initialize(r, pool)

	log.Printf("initialized all")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("server running")

	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile!="" && keyFile!= "" {
		log.Fatal(r.RunTLS(":"+port, certFile, keyFile))
	} else {
		r.Run(":" + port)
	}

}
