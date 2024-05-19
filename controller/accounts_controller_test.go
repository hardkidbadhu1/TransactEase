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
	"transact-api/model/dto/response"
	"transact-api/service/mocks"
)

type AccountControllerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	context     *gin.Context
	mockSvc     *mocks.MockAccountService
	recorder    *httptest.ResponseRecorder
	accountCtrl *AccountController
}

func (suite *AccountControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockSvc = mocks.NewMockAccountService(suite.mockCtrl)
	suite.accountCtrl = NewAccountController(suite.mockSvc)
}

func (suite *AccountControllerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestAccountController(t *testing.T) {
	suite.Run(t, new(AccountControllerTestSuite))
}

func (suite *AccountControllerTestSuite) TestCreateAccount_FailOnEmptyPayload() {
	suite.context.Request = httptest.NewRequest("POST", "/create-account", nil)
	suite.accountCtrl.CreateAccount(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *AccountControllerTestSuite) TestCreateAccount_FailOnServiceError() {
	req := request.AccountCreateRequest{
		DocumentNumber: "123456789",
	}

	reqStr, _ := json.Marshal(req)
	suite.mockSvc.EXPECT().InsertAccount(suite.context, req).Return(errors.New("some service error"))

	suite.context.Request = httptest.NewRequest("POST", "/create-account",
		bytes.NewBufferString(string(reqStr)))
	suite.accountCtrl.CreateAccount(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *AccountControllerTestSuite) TestCreateAccount_Success() {
	req := request.AccountCreateRequest{
		DocumentNumber: "123456789",
	}

	reqStr, _ := json.Marshal(req)
	suite.mockSvc.EXPECT().InsertAccount(suite.context, gomock.Any()).Return(nil)

	suite.context.Request = httptest.NewRequest("POST", "/create-account",
		bytes.NewBufferString(string(reqStr)))
	suite.accountCtrl.CreateAccount(suite.context)
	suite.Equal(http.StatusCreated, suite.recorder.Code)
}

func (suite *AccountControllerTestSuite) TestGetAccount_Success() {
	documentNumber := "1234567890"
	mockAccount := &response.AccountResponse{
		DocumentNumber: documentNumber,
	}

	suite.mockSvc.EXPECT().GetAccount(suite.context, documentNumber).Return(mockAccount, nil)

	suite.context.Request = httptest.NewRequest("GET", "/get-account/1234567890", nil)
	suite.context.Params = gin.Params{{Key: "documentNumber", Value: documentNumber}}
	suite.accountCtrl.GetAccount(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	var account response.AccountResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &account)
	suite.NoError(err)
	suite.Equal(mockAccount.DocumentNumber, account.DocumentNumber)
}

func (suite *AccountControllerTestSuite) TestGetAccount_Error() {
	documentNumber := "1234567890"
	suite.mockSvc.EXPECT().GetAccount(suite.context, documentNumber).Return(nil, errors.New("some error"))

	suite.context.Request = httptest.NewRequest("GET", "/get-account/1234567890", nil)
	suite.context.Params = gin.Params{{Key: "documentNumber", Value: documentNumber}}

	suite.accountCtrl.GetAccount(suite.context)

	suite.NotEqual(http.StatusOK, suite.recorder.Code)
}
