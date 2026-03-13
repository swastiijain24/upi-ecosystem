package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/services"
)

type Handler struct {
	service services.Service
}

func NewHandler(s services.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetAccountById(c *gin.Context){
	id := c.Param("id")
	
	account, err := h.service.GetAccountById(c.Request.Context(), id)
	if err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}

	c.JSON(200, gin.H{"account":account})

}

func (h *Handler) CreateAccount(c *gin.Context){
	var accountDetails CreateAccountReq

	if err:= c.ShouldBindJSON(&accountDetails); err !=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	account, err:= h.service.CreateAccount(c.Request.Context(), accountDetails.Name, accountDetails.Phone)
	if err !=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, account)
}

func (h *Handler) Debit(c *gin.Context){
	var debitReq DebitCreditRequest
	if err:= c.ShouldBindJSON(&debitReq); err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}

	transaction, err := h.service.Debit(c.Request.Context(), debitReq.AccountID, debitReq.Amount)
	if err!=nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, transaction)
}


func (h *Handler) Credit(c *gin.Context){
	var creditReq DebitCreditRequest
	if err := c.ShouldBindJSON(&creditReq); err!=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return 
	}
	transaction, err:= h.service.Credit(c.Request.Context(), creditReq.AccountID, creditReq.Amount)
	if err!=nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200, transaction)
}


func (h *Handler) GetTransactions(c *gin.Context){
	id := c.Param("id")
	
	transactions, err := h.service.GetTransactions(c.Request.Context(), id)
	if err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return 
	}
	c.JSON(200, transactions)

}


func (h *Handler) CheckBalance(c *gin.Context){
	id := c.Param("id")
	
	balance, err := h.service.CheckBalance(c.Request.Context(), id)
	if err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return 
	}
	c.JSON(200, balance)
}

type DebitCreditRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type CreateAccountReq struct{
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}