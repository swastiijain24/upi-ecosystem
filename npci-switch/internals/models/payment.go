package models

type PaymentRequest struct{
	SenderAccountID   string `json:"sender_account_id"`
	ReceiverAccountID string `json:"receiver_account_id"`
	Amount            int64  `json:"amount"`
}

type PaymentResponse struct{
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

