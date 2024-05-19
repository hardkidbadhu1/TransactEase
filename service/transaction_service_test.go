package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
	"transact-api/model/dto/request"
	"transact-api/repository/mocks"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockTransactionRepository
	context  *gin.Context
	service  TransactionService
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockTransactionRepository(suite.mockCtrl)
	suite.context, _ = gin.CreateTestContext(nil)
	suite.service = NewTransactionService(suite.mockRepo)
}

func (suite *TransactionServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestTransactionService(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

func (suite *TransactionServiceTestSuite) TestInsertTransaction() {
	suite.mockRepo.EXPECT().CreateTransaction(suite.context, gomock.Any()).Return(nil)
	err := suite.service.CreateTransaction(suite.context, request.TransactionCreateRequest{
		AccountID:       12,
		OperationTypeID: 23,
		Amount:          100,
	})
	suite.NoError(err)
}
