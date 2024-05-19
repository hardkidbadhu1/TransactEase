package application

import (
	"errors"
	"transact-api/configuration"
	"transact-api/router"
	"transact-api/utils"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Run the application
func Run(filePath string) {

	logger := utils.GetLogger(context.Background())
	// parse the configuration
	serverConfig, err := configuration.Parse(filePath)
	if err != nil {
		logger.Fatalf("Failed to parse config: %v", err)
	}

	// create a new server
	engine := gin.New(func(e *gin.Engine) {
		e.Use(gin.Recovery())
	})

	gin.SetMode(serverConfig.GetAppMode())

	// init router
	router.NewRouter(serverConfig, engine).InitRouter()

	srv := &http.Server{
		Addr:    ":" + serverConfig.GetPort(),
		Handler: engine,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Println("Server exiting")
}
