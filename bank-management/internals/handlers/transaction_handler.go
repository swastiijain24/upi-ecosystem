package handlers

import (
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


func (h *TransactionHandler) Debit(c *gin.Context){
	var debitReq DebitCreditRequest
	if err:= c.ShouldBindJSON(&debitReq); err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}

	transaction, err := h.TransactionService.Debit(c.Request.Context(), debitReq.FromAccountID, debitReq.ToAccountId, debitReq.Amount, debitReq.Description)
	if err!=nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, transaction)
}


func (h *TransactionHandler) Credit(c *gin.Context){
	var creditReq DebitCreditRequest
	if err := c.ShouldBindJSON(&creditReq); err!=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return 
	}
	transaction, err:= h.TransactionService.Credit(c.Request.Context(), creditReq.FromAccountID, creditReq.ToAccountId, creditReq.Amount, creditReq.Description)
	if err!=nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, transaction)
}


func (h *TransactionHandler) GetTransactions(c *gin.Context){
	id := c.Param("id")
	
	transactions, err := h.TransactionService.GetTransactions(c.Request.Context(), id)
	if err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return 
	}
	c.JSON(200, transactions)

}


type DebitCreditRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountId string `json:"to_account_id"`
	Amount    int64  `json:"amount"`
	Description string `json:"description"`
}
