package repository

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"net/http/httptest"
	"testing"
	"transact-api/model/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
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

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
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
	var mockAccounts = entities.Account{
		DocumentNumber: "1232383239",
	}

	expectedSQL := "INSERT INTO \"accounts\" (.+) VALUES (.+)"
	addRow := sqlmock.NewRows([]string{"document_number"}).AddRow("1232383239")
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(expectedSQL).WillReturnRows(addRow)
	suite.sqlMock.ExpectCommit()

	err := suite.repository.InsertAccount(suite.context, mockAccounts)
	suite.NoError(err)
	suite.Assert().Nil(suite.sqlMock.ExpectationsWereMet())
}

func (suite *AccountRepositoryTestSuite) TestInsertAccount_Error() {
	var mockAccounts = entities.Account{
		DocumentNumber: "1232383239",
	}

	expectedSQL := "INSERT INTO \"accounts\" (.+) VALUES (.+)"
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery(expectedSQL).WillReturnError(errors.New("some error"))
	suite.sqlMock.ExpectCommit()

	err := suite.repository.InsertAccount(suite.context, mockAccounts)

	suite.Error(err)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *AccountRepositoryTestSuite) TestFindAccountByDocumentNumber() {
	mockAccount := entities.Account{
		DocumentNumber: "1232383239",
	}

	addRow := sqlmock.NewRows([]string{"document_number"}).AddRow("1232383239")

	expectedSQL := "SELECT (.+) FROM \"accounts\" WHERE document_number =(.+)"

	suite.sqlMock.ExpectQuery(expectedSQL).WillReturnRows(addRow)

	account, err := suite.repository.FindAccountByDocumentNumber(suite.context, mockAccount.DocumentNumber)

	suite.NoError(err)
	suite.Equal(mockAccount.DocumentNumber, account.DocumentNumber)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *AccountRepositoryTestSuite) TestFindAccountByDocumentNumber_Error() {
	mockAccount := entities.Account{
		DocumentNumber: "1232383239",
	}

	expectedSQL := "SELECT (.+) FROM \"accounts\" WHERE document_number =(.+)"

	suite.sqlMock.ExpectQuery(expectedSQL).WillReturnError(errors.New("some error"))

	_, err := suite.repository.FindAccountByDocumentNumber(suite.context, mockAccount.DocumentNumber)

	suite.Error(err)
	suite.sqlMock.ExpectationsWereMet()
}
