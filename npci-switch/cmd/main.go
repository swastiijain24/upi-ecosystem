package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swastiijain24/npci-switch/internals/clients"
	"github.com/swastiijain24/npci-switch/internals/handlers"
	"github.com/swastiijain24/npci-switch/internals/routes"
	"github.com/swastiijain24/npci-switch/internals/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}


	r := gin.New()

	bankClient := clients.NewBankClient()
	paymentService := services.NewService(bankClient)
	paymentHandler := handlers.NewHandler(paymentService)
	routes.SetupRoutes(r, paymentHandler)

	port:= os.Getenv("PORT")

	if port ==""{
		port  = "8081"
	}

	r.Run(":"+port)
}