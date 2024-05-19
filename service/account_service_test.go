package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
	"transact-api/model/dto/request"
	"transact-api/model/entities"
	"transact-api/repository/mocks"
)

type AccountServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockAccountRepository
	context  *gin.Context
	service  AccountService
}

func (suite *AccountServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockAccountRepository(suite.mockCtrl)
	suite.context, _ = gin.CreateTestContext(nil)
	suite.service = NewAccountService(suite.mockRepo)
}

func (suite *AccountServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestAccountService(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}

func (suite *AccountServiceTestSuite) TestInsertAccount() {
	suite.mockRepo.EXPECT().InsertAccount(suite.context, gomock.Any()).Return(nil)
	err := suite.service.InsertAccount(suite.context, request.AccountCreateRequest{
		DocumentNumber: "1234444444",
	})
	suite.NoError(err)
}

func (suite *AccountServiceTestSuite) TestGetAccount() {
	suite.mockRepo.EXPECT().FindAccountByDocumentNumber(suite.context, "1234444444").Return(&entities.Account{
		DocumentNumber: "1234444444",
	}, nil)
	_, err := suite.service.GetAccount(suite.context, "1234444444")
	suite.NoError(err)
}

func (suite *AccountServiceTestSuite) TestGetAccountError() {
	suite.mockRepo.EXPECT().FindAccountByDocumentNumber(suite.context, "1234444444").Return(nil, errors.New("account not found"))
	_, err := suite.service.GetAccount(suite.context, "1234444444")
	suite.Error(err)
}
