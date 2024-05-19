//go:generate mockgen -source=transaction_repository.go -destination=mocks/transaction_repository_mock.go -package=mocks TransactionRepository
package repository

import (
	"transact-api/model/entities"
	"transact-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(ctx *gin.Context, transaction entities.Transaction) error
}

type transactionRepo struct {
	db *gorm.DB
}

func (t *transactionRepo) CreateTransaction(ctx *gin.Context, transaction entities.Transaction) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionRepository").
		WithField("Method", "CreateTransaction")

	result := t.db.Create(&transaction)
	if result.Error != nil {
		logger.Errorf("error while inserting transaction - %v", result.Error)
		return result.Error
	}

	return nil
}

func NewTransactionRepo(db *gorm.DB) TransactionRepository {
	return &transactionRepo{
		db: db,
	}
}
