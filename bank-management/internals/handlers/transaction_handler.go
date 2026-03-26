package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/dtos"
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
	var debitReq dtos.DebitCreditRequest
	if err:= c.ShouldBindJSON(&debitReq); err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}

	transaction, err := h.TransactionService.Debit(c.Request.Context(), debitReq)
	if err!=nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, transaction)
}


func (h *TransactionHandler) Credit(c *gin.Context){
	var creditReq dtos.DebitCreditRequest
	if err := c.ShouldBindJSON(&creditReq); err!=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return 
	}
	transaction, err:= h.TransactionService.Credit(c.Request.Context(), creditReq)
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