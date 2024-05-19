package router

import (
	docs "TransactEase/Docs"
	"TransactEase/configuration"
	"TransactEase/constants"
	"TransactEase/controller"
	"TransactEase/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	config configuration.Configuration
	engine *gin.Engine
}

func NewRouter(config configuration.Configuration, engine *gin.Engine) *Router {
	return &Router{
		config: config,
		engine: engine,
	}
}

func (r *Router) InitRouter() {
	group := r.engine.Group(constants.ApiVersion)

	r.engine.Use(middleware.ValidateHeader())

	//controller
	healthCtrl := controller.NewController()

	docs.SwaggerInfo.Version = constants.AppVersion
	docs.SwaggerInfo.Host = r.config.GetAppHost()
	docs.SwaggerInfo.Title = "TransactEase API"
	docs.SwaggerInfo.Description = "TransactEase server API documentation"
	docs.SwaggerInfo.BasePath = constants.ApiVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	group.GET("/healthz", healthCtrl.Healthz)
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
