package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/services"
)

type AccountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(s services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: s,
	}
}

func (h *AccountHandler) GetAccountById(c *gin.Context){
	id := c.Param("id")
	
	account, err := h.accountService.GetAccountById(c.Request.Context(), id)
	if err!= nil{
		c.JSON(404, gin.H{"error":err.Error()})
		return
	}

	c.JSON(200, gin.H{"account":account})
}

func (h *AccountHandler) CreateAccount(c *gin.Context){
	var accountDetails CreateAccountReq

	if err:= c.ShouldBindJSON(&accountDetails); err !=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	account, err:= h.accountService.CreateAccount(c.Request.Context(), accountDetails.Name, accountDetails.Phone)
	if err !=nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, account)
}


func (h *AccountHandler) GetBalance(c *gin.Context){
	id := c.Param("id")
	
	balance, err := h.accountService.GetBalance(c.Request.Context(), id)
	if err!= nil{
		c.JSON(400, gin.H{"error":err.Error()})
		return 
	}
	c.JSON(200, balance)
}


func (h *AccountHandler) DeleteAccount(c *gin.Context){
	id := c.Param("id")

	err := h.accountService.DeleteAccount(c.Request.Context(), id)
	if err!=nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return 
	}
	c.JSON(204, nil)
}

type CreateAccountReq struct {
    Name  string `json:"name" binding:"required,min=1,max=255"`
    Phone string `json:"phone" binding:"required,e164"`  
}