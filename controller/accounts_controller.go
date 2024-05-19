package controller

import (
	"transact-api/api_error"
	"transact-api/model/dto/request"
	"transact-api/service"
	"transact-api/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	svc service.AccountService
}

func NewAccountController(svc service.AccountService) *AccountController {
	return &AccountController{
		svc: svc,
	}
}

// CreateAccount godoc
// @Tags Account
// @Summary Create account endpoint
// @Description Create account endpoint
// @Accept json
// @Produce json
// @Param account body request.AccountCreateRequest true "Account details"
// @Success 201
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 500 {object} api_error.ErrorResponse
// @Router /accounts [post]
func (a *AccountController) CreateAccount(ctx *gin.Context) {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "controller").
		WithField("Method", "CreateAccount")
	var req request.AccountCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Error binding request: %s", err.Error())
		ctx.AbortWithStatusJSON(api_error.InvalidDocumentNumber.HttpStatusCode, api_error.InvalidDocumentNumber)
		return
	}

	if err := a.svc.InsertAccount(ctx, req); err != nil {
		logger.Errorf("Error inserting account: %s", err.Error())
		srvErr := api_error.NewInternalServerError(err.Error())
		ctx.AbortWithStatusJSON(srvErr.HttpStatusCode, srvErr)
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

// GetAccount godoc
// @Tags Account
// @Summary Get account endpoint
// @Description Get account endpoint
// @Accept json
// @Produce json
// @Param documentNumber path string true "Document number"
// @Success 200 {object} response.AccountResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Failure 500 {object} api_error.ErrorResponse
// @Router /accounts/{documentNumber} [get]
func (a *AccountController) GetAccount(ctx *gin.Context) {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "controller").
		WithField("Method", "GetAccount")
	documentNumber := ctx.Param("documentNumber")
	account, err := a.svc.GetAccount(ctx, documentNumber)
	if err != nil {
		logger.Errorf("Error fetching account: %s", err.Error())
		srvErr := api_error.NewInternalServerError(err.Error())
		ctx.AbortWithStatusJSON(srvErr.HttpStatusCode, srvErr)
		return
	}

	ctx.AbortWithStatusJSON(200, account)
}
