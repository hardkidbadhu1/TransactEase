package request

import "transact-api/model/entities"

type AccountCreateRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

func (req AccountCreateRequest) ToEntity() entities.Account {
	return entities.Account{
		DocumentNumber: req.DocumentNumber,
	}
}
