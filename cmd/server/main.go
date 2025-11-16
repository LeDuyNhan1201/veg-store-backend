package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
	"veg-store-backend/internal/infrastructure/identity"
	"veg-store-backend/internal/infrastructure/repository"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/middleware"
	"veg-store-backend/internal/rest_api/rest_handler"
	"veg-store-backend/internal/rest_api/route"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title RESTFUL API VERSION 1.0
// @version 1.0
// @description This is an API document for veg-store-backend.
// @termsOfService http://example.com/terms/

// @contact.name Nhan Le
// @contact.url http://example.com/support
// @contact.email benlun1201@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

// @schemes http https

func main() {
	dependencies := fx.Options(
		fx.Provide(core.Init),
		fx.Provide(data.InitPostgresDB),
		repository.Module,
		identity.Module,
		service.Module,
		rest_handler.Module,
		fx.Provide(router.InitHTTPRouter),
		middleware.Module,
		route.Module,
	)

	app := fx.New(
		dependencies,

		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			router *router.HTTPRouter,
			middlewaresCollection *middleware.MiddlewaresCollection,
			routesCollection *route.RoutesCollection,
			postgresDB *data.PostgresDB,
			dataSeeder service.DataSeederService,
		) {
			// Setup middlewares and routes
			middlewaresCollection.Setup()
			routesCollection.Setup()

			// Run DB migrations
			if err := postgresDB.Migrate(); err != nil {
				zap.NewExample().Fatal("ORM failed to migrate: ", zap.Error(err))
			}
			dataSeeder.SeedData()

			lifecycle.Append(fx.Hook{
				OnStart: func(context context.Context) error {
					go func() {
						if err := router.HttpRun(); err != nil {
							zap.NewExample().Fatal("Server failed to start: ", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(context context.Context) error {
					zap.NewExample().Info("Shutting down server...")
					return nil
				},
			})
		}),
	)

	// Graceful shutdown
	startContext, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Start(startContext); err != nil {
		zap.NewExample().Fatal("Server failed to start: ", zap.Error(err))
	}

	// Wait for OS signal to terminate
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	stopContext, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Stop(stopContext); err != nil {
		zap.NewExample().Fatal("Failed to stop application: ", zap.Error(err))
	}
}
