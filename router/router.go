package router

import (
	docs "transact-api/Docs"
	"transact-api/configuration"
	"transact-api/constants"
	"transact-api/controller"
	"transact-api/database"
	"transact-api/middleware"
	"transact-api/repository"
	"transact-api/service"
	"transact-api/utils"

	"context"

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

	logger := utils.GetLogger(context.Background())
	db, err := database.ConnectDB(r.config)
	if err != nil {
		logger.Fatalf("error: connecting to db: %s", err.Error())
	}

	// repository
	accountRepo := repository.NewAccountRepository(db)
	transactRepo := repository.NewTransactionRepo(db)

	//SVC
	transactSVC := service.NewTransactionService(transactRepo)
	accountSvc := service.NewAccountService(accountRepo)

	//controller
	healthCtrl := controller.NewController()
	accCtrl := controller.NewAccountController(accountSvc)
	transactCtrl := controller.NewTransactionController(transactSVC)

	docs.SwaggerInfo.Version = constants.AppVersion
	docs.SwaggerInfo.Host = r.config.GetAppHost()
	docs.SwaggerInfo.Title = "transact-api"
	docs.SwaggerInfo.Description = "transact-api server API documentation"
	docs.SwaggerInfo.BasePath = constants.ApiVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	group.GET("/healthz", healthCtrl.Healthz)
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	group.POST("/accounts", accCtrl.CreateAccount)
	group.GET("/accounts/:documentNumber", accCtrl.GetAccount)
	group.POST("/transaction", transactCtrl.CreateTransaction)
}
