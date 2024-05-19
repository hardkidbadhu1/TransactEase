package entities

import "transact-api/model/dto/response"

type Account struct {
	ID             int    `gorm:"column:account_id;primary_key"`
	DocumentNumber string `gorm:"column:document_number"`
}

func (a Account) ToResponse() *response.AccountResponse {
	return &response.AccountResponse{
		ID:             a.ID,
		DocumentNumber: a.DocumentNumber,
	}
}

func (a Account) TableName() string {
	return "accounts"
}
