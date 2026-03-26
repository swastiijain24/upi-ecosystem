package dtos

type DebitCreditRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountId string `json:"to_account_id"`
	Amount    int64  `json:"amount"`
	Description string `json:"description"`
}
