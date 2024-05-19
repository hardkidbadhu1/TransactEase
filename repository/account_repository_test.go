package repository

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
	"transact-api/model/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	mockCtrl   *gomock.Controller
	context    *gin.Context
	sqlMock    sqlmock.Sqlmock
	repository AccountRepository
	recorder   *httptest.ResponseRecorder
}

func TestAccountRepository(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}

func (suite *AccountRepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, suite.sqlMock, err = sqlmock.New()
	if err != nil {
		suite.Fail("an error was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		suite.Fail("an error was not expected when opening gorm database", err)
	}

	suite.recorder = httptest.NewRecorder()
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.repository = NewAccountRepository(gormDB)
}

func (suite *AccountRepositoryTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *AccountRepositoryTestSuite) TestInsertAccount() {
	mockAccount := entities.Account{
		DocumentNumber: "1234567890",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO").
		WithArgs(mockAccount.DocumentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	err := suite.repository.InsertAccount(suite.context, mockAccount)

	suite.NoError(err)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *AccountRepositoryTestSuite) TestInsertAccount_Error() {
	mockAccount := entities.Account{
		DocumentNumber: "1234567890",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("INSERT INTO").
		WithArgs(mockAccount.DocumentNumber).
		WillReturnError(errors.New("some error"))
	suite.sqlMock.ExpectRollback()

	err := suite.repository.InsertAccount(suite.context, mockAccount)

	suite.Error(err)
	suite.EqualError(err, "some error")
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *AccountRepositoryTestSuite) TestFindAccountByDocumentNumber() {
	mockAccount := entities.Account{
		DocumentNumber: "1234567890",
	}

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM").
		WithArgs(mockAccount.DocumentNumber).
		WillReturnRows(sqlmock.NewRows([]string{"document_number"}).AddRow(mockAccount.DocumentNumber))

	account, err := suite.repository.FindAccountByDocumentNumber(suite.context, mockAccount.DocumentNumber)

	suite.NoError(err)
	suite.Equal(mockAccount.DocumentNumber, account.DocumentNumber)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *AccountRepositoryTestSuite) TestFindAccountByDocumentNumber_Error() {
	mockAccount := entities.Account{
		DocumentNumber: "1234567890",
	}

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM").
		WithArgs(mockAccount.DocumentNumber).
		WillReturnError(errors.New("some error"))

	_, err := suite.repository.FindAccountByDocumentNumber(suite.context, mockAccount.DocumentNumber)

	suite.Error(err)
	suite.EqualError(err, "some error")
	suite.sqlMock.ExpectationsWereMet()
}
