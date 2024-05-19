package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"transact-api/api_error"
	"transact-api/model/dto/request"
	"transact-api/service"
	"transact-api/utils"
)

type TransactionController struct {
	svc service.TransactionService
}

func NewTransactionController(svc service.TransactionService) *TransactionController {
	return &TransactionController{
		svc: svc,
	}
}

// CreateTransaction godoc
// @Tags Transaction
// @Summary Create transaction endpoint
// @Description Create transaction endpoint
// @Accept json
// @Produce json
// @Param transaction body request.TransactionCreateRequest true "Transaction details"
// @Success 201
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 500 {object} api_error.ErrorResponse
// @Router /transaction [post]
func (t TransactionController) CreateTransaction(ctx *gin.Context) {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "TransactionController").
		WithField("Method", "CreateTransaction")

	var req request.TransactionCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("error while binding request - %v", err)
		ctx.AbortWithStatusJSON(api_error.InvalidParams.HttpStatusCode, api_error.InvalidParams)
		return
	}

	if err := t.svc.CreateTransaction(ctx, req); err != nil {
		logger.Errorf("error while creating transaction - %v", err)
		srvErr := api_error.NewInternalServerError(err.Error())
		ctx.AbortWithStatusJSON(srvErr.HttpStatusCode, srvErr)
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}
