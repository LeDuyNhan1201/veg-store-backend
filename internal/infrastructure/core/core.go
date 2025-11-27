package core

import (
	"os"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/infrastructure/config"
	"veg-store-backend/internal/infrastructure/localizer"
	"veg-store-backend/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Core struct {
	AppConfig *config.Config
	Localizer *localizer.Localizer
	Logger    *zap.Logger
	Error     *exception.AppError
}

func Init() *Core {
	mode := determineMode()
	return &Core{
		Error:     exception.Init(),
		Localizer: localizer.Init(mode),
		Logger:    logger.Init(mode),
		AppConfig: config.Init(mode),
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
