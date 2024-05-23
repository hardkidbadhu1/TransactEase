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
	mockCtrl   *gomock.Controller
	context    *gin.Context
	sqlMock    sqlmock.Sqlmock
	repository TransactionRepository
	recorder   *httptest.ResponseRecorder
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
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransactionRepositoryTestSuite) TestInsertAccount() {
	var mockTransaction = entities.Transaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          1,
	}

	expectedSQL := "INSERT INTO \"transactions\" (.+) VALUES (.+)"
	addRow := sqlmock.NewRows([]string{"account_id"}).AddRow("1")
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(expectedSQL).WillReturnRows(addRow)
	suite.sqlMock.ExpectCommit()

	err := suite.repository.CreateTransaction(suite.context, mockTransaction)
	suite.NoError(err)
	suite.Assert().Nil(suite.sqlMock.ExpectationsWereMet())
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
