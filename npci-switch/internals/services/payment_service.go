package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/swastiijain24/npci-switch/internals/clients"
	"github.com/swastiijain24/npci-switch/internals/models"
)

type Service interface {
	InitiatePayment(ctx context.Context, sender_accountId string , receiver_accountId string , amount int64)  (models.PaymentResponse, error)
}

type svc struct {
	bankClient clients.BankClient
}

func NewService(bankClient clients.BankClient) Service {
	return &svc{
		bankClient: bankClient,
	}
}

func (s *svc) InitiatePayment(ctx context.Context,senderAccountId string , receiverAccountId string , amount int64) (models.PaymentResponse, error) {
	txId := generateId()

	if err:= s.bankClient.Debit(senderAccountId, amount); err!=nil{
		return models.PaymentResponse{}, err
	}

	if err:= s.bankClient.Credit(receiverAccountId, amount); err!=nil{
		s.bankClient.Credit(senderAccountId, amount)
		return models.PaymentResponse{}, err
	}

	return models.PaymentResponse{
		TransactionID: txId,
		Status: "SUCCESS",
	} , nil

}

func generateId() string {
	return "tx_" + uuid.New().String()
}
