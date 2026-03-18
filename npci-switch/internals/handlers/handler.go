package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
	 err := h.service.InitiatePayment(c.Request.Context(), paymentParams.PayerVpa, paymentParams.PayerBank,  paymentParams.PayeeVpa, paymentParams.PayeeBank, paymentParams.Amount, paymentParams.ReferenceID)
	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, err)
}


type PaymentReq struct {
	PayerVpa    string             `json:"payer_vpa"`
	PayerBank   string             `json:"payer_bank"`
	PayeeVpa    string             `json:"payee_vpa"`
	PayeeBank   string             `json:"payee_bank"`
	Amount      int64              `json:"amount"`
	ReferenceID pgtype.Text        `json:"reference_id"`

}

