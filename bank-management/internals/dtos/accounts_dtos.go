package dtos

type CreateAccountReq struct{
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}