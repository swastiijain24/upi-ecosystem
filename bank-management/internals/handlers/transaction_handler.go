package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/services"
)

type TransactionHandler struct {
	TransactionService services.TransactionService
}

func NewTransactionHandler(s services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		TransactionService: s,
	}
}

func (h *TransactionHandler) Debit(c *gin.Context) {
	var debitReq DebitRequest
	if err := c.ShouldBindJSON(&debitReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.TransactionService.Debit(c.Request.Context(), debitReq.FromAccountID, debitReq.ToAccountId, debitReq.Amount, debitReq.Description, debitReq.MpinHash, debitReq.ExternalId)
	response := Response{
		bankReferenceId: transaction.ID.String(),
		status: transaction.Status,
		created_at: transaction.CreatedAt.Time,
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, response)
}

func (h *TransactionHandler) Credit(c *gin.Context) {
	var creditReq CreditRequest
	if err := c.ShouldBindJSON(&creditReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	transaction, err := h.TransactionService.Credit(c.Request.Context(), creditReq.FromAccountID, creditReq.ToAccountId, creditReq.Amount, creditReq.Description, creditReq.ExternalId)
	response := Response{
		bankReferenceId: transaction.ID.String(),
		status: transaction.Status,
		created_at: transaction.CreatedAt.Time,
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, response)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	id := c.Param("id")

	transactions, err := h.TransactionService.GetTransactions(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, transactions)

}

func (h *TransactionHandler) GetStatusByExternalId(c *gin.Context){
	var statusReq StatusReq
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	txnId, status, err := h.TransactionService.GetStatusByExternalId(c.Request.Context(), statusReq.ExternalId)
	response := Response{
		bankReferenceId: txnId,
		status: status,
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, response)

}

type DebitRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountId   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Description   string `json:"description"`
	MpinHash      string `json:"mpin_hash" binding:"required,e164"`
	ExternalId    string `json:"external_id" binding:"required"`
}

type CreditRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountId   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Description   string `json:"description"`
	ExternalId    string `json:"external_id" binding:"required"`
}

type StatusReq struct {
	ExternalId string `json:"external_id" binding:"required"`
}

type Response struct{
	bankReferenceId string 
	status string 
	created_at  time.Time
}