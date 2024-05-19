package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"transact-api/model/dto/request"
	"transact-api/service/mocks"
)

type TransactionControllerTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	context      *gin.Context
	mockSvc      *mocks.MockTransactionService
	recorder     *httptest.ResponseRecorder
	transactCtrl *TransactionController
}

func (suite *TransactionControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockSvc = mocks.NewMockTransactionService(suite.mockCtrl)
	suite.transactCtrl = NewTransactionController(suite.mockSvc)
}

func (suite *TransactionControllerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestTransactionController(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}

func (suite *TransactionControllerTestSuite) TestCreateAccount_FailOnEmptyPayload() {
	suite.context.Request = httptest.NewRequest("POST", "/create-account", nil)
	suite.transactCtrl.CreateTransaction(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *TransactionControllerTestSuite) TestCreateAccount_FailOnServiceError() {
	req := request.TransactionCreateRequest{
		AccountID:       1,
		Amount:          100,
		OperationTypeID: 1,
	}

	reqStr, _ := json.Marshal(req)
	suite.mockSvc.EXPECT().CreateTransaction(suite.context, req).Return(errors.New("some service error"))

	suite.context.Request = httptest.NewRequest("POST", "/create-account",
		bytes.NewBufferString(string(reqStr)))
	suite.transactCtrl.CreateTransaction(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *TransactionControllerTestSuite) TestCreateAccount_Success() {
	req := request.TransactionCreateRequest{
		AccountID:       1,
		Amount:          100,
		OperationTypeID: 1,
	}

	reqStr, _ := json.Marshal(req)
	suite.mockSvc.EXPECT().CreateTransaction(suite.context, gomock.Any()).Return(nil)

	suite.context.Request = httptest.NewRequest("POST", "/create-account",
		bytes.NewBufferString(string(reqStr)))
	suite.transactCtrl.CreateTransaction(suite.context)
	suite.Equal(http.StatusCreated, suite.recorder.Code)
}
