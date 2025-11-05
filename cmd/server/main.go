package main

import (
	"log"
	"os"
	"veg-store-backend/injection"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

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
