package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swastiijain24/npci-switch/internals/clients"
	"github.com/swastiijain24/npci-switch/internals/db"
	"github.com/swastiijain24/npci-switch/internals/handlers"
	"github.com/swastiijain24/npci-switch/internals/redis"
	"github.com/swastiijain24/npci-switch/internals/routes"
	"github.com/swastiijain24/npci-switch/internals/services"
	"github.com/swastiijain24/npci-switch/internals/workers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()

	//we connect to the databse before starting the server obv 

	connString := os.Getenv("DBSTRING")

	pool := db.NewPostgres(connString)

	log.Println(pool)

	redisClient := redis.NewRedis()
	log.Println(redisClient)


	r := gin.New()

	bankClient := clients.NewBankClient()
	paymentService := services.NewService(bankClient)
	paymentHandler := handlers.NewHandler(paymentService)
	routes.SetupRoutes(r, paymentHandler)

	port:= os.Getenv("PORT")

	if port ==""{
		port  = "8081"
	}


	//start the workers
	debitWorker := workers.NewDebitWorker(redisClient)
	go debitWorker.Start((ctx))
	creditWorker := workers.NewCreditWorker(redisClient)
	go creditWorker.Start(ctx)


	
	if err:= r.Run(":"+port); err!=nil{
		log.Fatal("failed to start server")
	}
	
}