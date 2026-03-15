package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/npci-switch/internals/handlers"
)

func SetupRoutes(r *gin.Engine, paymentHandler *handlers.Handler) {
	r.POST("/payments", paymentHandler.InitiatePayment)
}