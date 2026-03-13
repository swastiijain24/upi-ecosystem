package models

import "time"

type Transaction struct {
	ID         		string
	Account_id 		string
	Type 			TransactionType
	Amount     		int64
	Status     		TransactionStatus
	Created_at 		time.Time
}

type TransactionType string

const(
	Credit TransactionType = "CREDIT"
	Debit TransactionType = "DEBIT"
)

type TransactionStatus string 

const(
	Pending TransactionStatus = "PENDING"
	Success TransactionStatus = "SUCCESS"
	Failed TransactionStatus = "FAILED"
	Reversed TransactionStatus = "REVERSED"
)