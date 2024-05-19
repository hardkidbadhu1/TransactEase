//go:generate mockgen -source=account_repository.go -destination=mocks/account_repository_mock.go -package=mocks AccountRepository
package repository

import (
	"transact-api/model/entities"
	"transact-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountRepository interface
type AccountRepository interface {
	InsertAccount(ctx *gin.Context, account entities.Account) error
	FindAccountByDocumentNumber(ctx *gin.Context, documentNumber string) (*entities.Account, error)
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{
		db: db,
	}
}

func (repository *accountRepo) InsertAccount(ctx *gin.Context, account entities.Account) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "AccountRepository").
		WithField("Method", "InsertAccount")
	result := repository.db.Create(&account)
	if result.Error != nil {
		logger.Errorf("error while inserting account - %v", result.Error)
		return result.Error
	}

	return nil
}

func (repository *accountRepo) FindAccountByDocumentNumber(ctx *gin.Context, documentNumber string) (*entities.Account, error) {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "AccountRepository").
		WithField("Method", "FindAccountByDocumentNumber")
	var account = new(entities.Account)
	result := repository.db.Where("document_number = ?", documentNumber).Find(account)
	if result.Error != nil {
		logger.Errorf("error while fetching account - %v", result.Error)
		return nil, result.Error
	}

	return account, nil
}
