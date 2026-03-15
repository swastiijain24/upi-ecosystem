package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/npci-switch/internals/services"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitiatePayment(c *gin.Context) {
	var paymentParams PaymentReq
	if err:= c.ShouldBindJSON(&paymentParams); err !=nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	paymentResponse, err := h.service.InitiatePayment(c.Request.Context(), paymentParams.SenderAccountId, paymentParams.ReceiverAccountId, paymentParams.Amount)
	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, paymentResponse)
}

type PaymentReq struct {
	SenderAccountId string `json:"sender_account_id"`
	ReceiverAccountId string `json:"receiver_account_id"`
	Amount int64 `json:"amount"`
}