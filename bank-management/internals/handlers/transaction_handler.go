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
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response := Response{
		BankReferenceId: transaction.ID.String(),
		Status:          transaction.Status,
		Created_at:      transaction.CreatedAt.Time,
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
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response := Response{
		BankReferenceId: transaction.ID.String(),
		Status:          transaction.Status,
		Created_at:      transaction.CreatedAt.Time,
	}

	c.JSON(201, response)
}

func (h *TransactionHandler) Refund(c *gin.Context) {
	var refundReq RefundRequest
	if err := c.ShouldBindJSON(&refundReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	transaction, err := h.TransactionService.Refund(c.Request.Context(), refundReq.FromAccountID, refundReq.ToAccountId, refundReq.Amount, refundReq.ExternalId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response := Response{
		BankReferenceId: transaction.ID.String(),
		Status:          transaction.Status,
		Created_at:      transaction.CreatedAt.Time,
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

func (h *TransactionHandler) GetStatusOfTransaction(c *gin.Context) {
	var statusReq StatusReq
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	txnId, status, err := h.TransactionService.GetStatus(c.Request.Context(), statusReq.ExternalId, statusReq.TransactionType)
	response := Response{
		BankReferenceId: txnId,
		Status:          status,
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
	MpinHash      string `json:"mpin_hash" binding:"required"`
	ExternalId    string `json:"external_id" binding:"required"`
}

type CreditRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountId   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Description   string `json:"description"`
	ExternalId    string `json:"external_id" binding:"required"`
}

type RefundRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountId   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	ExternalId    string `json:"external_id" binding:"required"`
}

type StatusReq struct {
	ExternalId string `json:"external_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"`
}

type Response struct {
	BankReferenceId string `json:"bank_reference_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
	Created_at      time.Time `json:"created_at"`
}