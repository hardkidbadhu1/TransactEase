package application

import (
	"TransactEase/configuration"
	"TransactEase/router"
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run the application
func Run(filePath string) {

	// parse the configuration
	serverConfig, err := configuration.Parse(filePath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
