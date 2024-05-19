//go:generate mockgen -source=transaction_service.go -destination=mocks/transaction_service.go -package=mocks
package service

import (
	"github.com/gin-gonic/gin"
	"transact-api/model/dto/request"
	"transact-api/repository"
)

type TransactionService interface {
	CreateTransaction(ctx *gin.Context, req request.TransactionCreateRequest) error
}

type transactionService struct {
	repo repository.TransactionRepository
}

func (a transactionService) CreateTransaction(ctx *gin.Context, req request.TransactionCreateRequest) error {
	return a.repo.CreateTransaction(ctx, req.ToEntity())
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}
