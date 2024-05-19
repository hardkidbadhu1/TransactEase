package controller

import (
	"TransactEase/constants"
	"TransactEase/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HealthController struct {
}

func NewController() *HealthController {
	return new(HealthController)
}

// HealthCheck godoc
// @Tags Health
// @Summary Health check endpoint
// @Description Health check endpoint
// @Accept json
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /healthz [get]
func (h HealthController) Healthz(ctx *gin.Context) {
	logger := log.WithField("Class", "controller").
		WithField("Method", "Healthz")

	logger.Info("Health check method initiated")
	ctx.JSON(200, model.HealthResponse{
		Status:  "ok",
		Version: constants.AppVersion,
	})
	logger.Info("Health check call completed")
}
