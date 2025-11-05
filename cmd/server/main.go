package main

import (
	"log"
	"os"
	"veg-store-backend/injection"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

// @title Example API
// @version 1.0
// @description This is a sample server that uses JWT authentication.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

// @schemes http

var (
	appGroup errgroup.Group
)

func main() {
	// Dependencies injection
	container := injection.Inject(determineMode())

	// HttpRun app
	appGroup.Go(func() error {
		return container.HttpRun()
	})

	if err := appGroup.Wait(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

func determineMode() string {
	mode := os.Getenv("MODE")
	switch mode {
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if mode == "" {
		zap.NewExample().Warn("No 'MODE' is defined. Server will run in 'dev' mode by default.")
		return "dev"
	}
	return mode
}
