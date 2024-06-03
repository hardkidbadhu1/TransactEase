package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http/httptest"
	"testing"
	"transact-api/model/entities"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	context     *gin.Context
	sqlMock     sqlmock.Sqlmock
	repository  TransactionRepository
	recorder    *httptest.ResponseRecorder
	transaction *gorm.DB
}

func TestTransactionRepository(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, suite.sqlMock, err = sqlmock.New()
	if err != nil {
		suite.Fail("an error was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		suite.Fail("an error was not expected when opening gorm database", err)
	}

	suite.recorder = httptest.NewRecorder()
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.repository = NewTransactionRepo(gormDB)
	suite.transaction = gormDB
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransactionRepositoryTestSuite) TestCreateTransaction() {
	mockTransaction := entities.Transaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          -50,
		Balance:         -50,
	}

	suite.transaction = suite.transaction.Begin()
	defer suite.transaction.Commit()
	suite.sqlMock.ExpectBegin()
	addRow := sqlmock.NewRows([]string{"account_id", "operation_type_id", "amount", "balance"}).
		AddRow("1", "1", "-50", "-50")
	suite.sqlMock.ExpectQuery("INSERT INTO \"transactions\" (.+) VALUES (.+)").
		WillReturnRows(addRow)
	suite.sqlMock.ExpectCommit()
	err := suite.repository.CreateTransaction(suite.context, mockTransaction)
	suite.NoError(err)
	err = suite.sqlMock.ExpectationsWereMet()
	suite.NoError(err, "there were unfulfilled expectations")
}

func (suite *TransactionRepositoryTestSuite) TestInsertAccountFailure() {
	mockTransaction := entities.Transaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          100,
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO").
		WithArgs(mockTransaction.AccountID, mockTransaction.OperationTypeID, mockTransaction.Amount).
		WillReturnError(errors.New("insert error"))
	suite.sqlMock.ExpectRollback()

	err := suite.repository.CreateTransaction(suite.context, mockTransaction)

	suite.Error(err)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *TransactionRepositoryTestSuite) TestAdjustBalance() {
	creditTransaction := entities.Transaction{
		ID:              4,
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          80,
	}

	suite.sqlMock.ExpectBegin()

	suite.transaction = suite.transaction.Begin()

	rows := sqlmock.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "balance"}).
		AddRow("1", "1", "1", "-50", "-50").
		AddRow("2", "1", "1", "-30", "-30")
	suite.sqlMock.ExpectQuery("SELECT (.+) FROM \"transactions\" WHERE account_id =(.+) AND balance < 0").
		WillReturnRows(rows)
	suite.sqlMock.ExpectExec("UPDATE \"transactions\" SET (.+) WHERE \"transaction_id\" = (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec("UPDATE \"transactions\" SET (.+) WHERE \"transaction_id\" = (.+)").
		WillReturnResult(sqlmock.NewResult(2, 1))
	suite.sqlMock.ExpectCommit()

	err := suite.repository.AdjustBalance(suite.context, &creditTransaction, suite.transaction)

	suite.NoError(err)
	suite.transaction.Commit()
	err = suite.sqlMock.ExpectationsWereMet()
	suite.NoError(err, "there were unfulfilled expectations")
}
