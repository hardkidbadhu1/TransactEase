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
	AdjustBalance(ctx *gin.Context, transaction *entities.Transaction, dbTxn *gorm.DB) error
}

type transactionRepo struct {
	db *gorm.DB
}

func (t *transactionRepo) AdjustBalance(ctx *gin.Context, transaction *entities.Transaction, dbTxn *gorm.DB) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionRepository").
		WithField("Method", "AdjustBalance")

	// fetch all the transactions with the account number and balance is negative
	var transactions []entities.Transaction
	result := dbTxn.Where("account_id = ? AND balance < 0", transaction.AccountID).Find(&transactions)
	if result.Error != nil {
		logger.Errorf("error while fetching transactions - %v", result.Error)
		return result.Error
	}

	balanceAdjustment := transaction.Balance
	logger.Info("balance adjustment - ", balanceAdjustment, transaction.Balance)
	for _, tx := range transactions {
		if -1*tx.Balance < balanceAdjustment {
			balanceAdjustment = tx.Balance + balanceAdjustment
			tx.Balance = 0
			logger.Info("tx to update", tx)
			if err := t.updateBalance(ctx, dbTxn, tx); err != nil {
				return err
			}
			continue
		}
		tx.Balance = tx.Balance + balanceAdjustment
		if err := t.updateBalance(ctx, dbTxn, tx); err != nil {
			return err
		}
	}

	return nil
}

func (t *transactionRepo) CreateTransaction(ctx *gin.Context, transaction entities.Transaction) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionRepository").
		WithField("Method", "CreateTransaction")

	if transaction.Balance < 0 {
		return t.createTransaction(ctx, t.db, &transaction)
	}

	txn := t.db.Begin()
	if txn.Error != nil {
		logger.Errorf("error while starting transaction - %v", txn.Error)
		return txn.Error
	}

	if err := t.createTransaction(ctx, txn, &transaction); err != nil {
		txn.Rollback()
		return err
	}

	if err := t.AdjustBalance(ctx, &transaction, txn); err != nil {
		txn.Rollback()
		return err
	}

	txn.Commit()
	return nil
}

func (t *transactionRepo) createTransaction(ctx *gin.Context, dbTxn *gorm.DB, transaction *entities.Transaction) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionRepository").
		WithField("Method", "createTransaction")

	result := dbTxn.Create(transaction)
	if result.Error != nil {
		logger.Errorf("error while inserting transaction - %v", result.Error)
		return result.Error
	}
	return nil
}

func (t *transactionRepo) updateBalance(ctx *gin.Context, dbTxn *gorm.DB, transaction entities.Transaction) error {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionRepository").
		WithField("Method", "createTransaction")

	result := dbTxn.Model(&transaction).Update("balance", transaction.Balance).Where("transaction_id = ?", transaction.ID)
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
