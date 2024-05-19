package controller

import (
	"github.com/gin-gonic/gin"
	"transact-api/constants"
	"transact-api/model/dto/response"
	"transact-api/utils"
)

type HealthController struct {
}

func NewController() *HealthController {
	return &HealthController{}
}

// HealthCheck godoc
// @Tags Health
// @Summary Health check endpoint
// @Description Health check endpoint
// @Accept json
// @Produce json
// @Success 200 {object} response.HealthResponse
// @Router /healthz [get]
func (h HealthController) Healthz(ctx *gin.Context) {
	logger := utils.GetLogger(ctx)
	logger = logger.WithField("Class", "controller").
		WithField("Method", "Healthz")

	logger.Info("Health check method initiated")
	ctx.JSON(200, response.HealthResponse{
		Status:  "ok",
		Version: constants.AppVersion,
	})
	logger.Info("Health check call completed")
}
