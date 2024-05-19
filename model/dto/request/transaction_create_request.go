package request

import "transact-api/model/entities"

type TransactionCreateRequest struct {
	AccountID       uint    `json:"account_id" binding:"required"`
	OperationTypeID uint    `json:"operation_type_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
}

func (req TransactionCreateRequest) ToEntity() entities.Transaction {
	return entities.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
	}
}
