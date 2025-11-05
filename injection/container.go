package injection

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/api/rest"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/infrastructure/data"
)

/*
This file sets up the dependency injection container for the application.
Logic:
- Load configuration settings
- Initialize Logger and i18n Translator
- Inject handlers, services, and repositories
- setup routes with injected dependencies
- Provide a method to run the HTTP server
*/

type Container struct {
	Router *Router
}

func Inject(mode string) *Container {
	core.Configs.Mode = mode
	core.Logger = core.InitLogger()       // Initialize Logger
	core.Configs = core.Load()            // Load configuration
	core.Translator = core.InitI18n()     // Initialize i18n Translator
	core.Error = exception.InitAppError() // Initialize App Error

	// setup routes
	router := &Router{}
	router.Setup()

	return &Container{
		Router: router,
	}
}

func (container *Container) HttpRun() error {
	// Inject handlers
	container.Router.RegisterUserRoutes(injectUserHandler())
	// TODO: Inject other handlers as needed

	if core.Configs.Mode != "prod" && core.Configs.Mode != "production" {
		// Register Swagger UI in non-production modes
		container.Router.RegisterSwaggerUI()
	}

	// Run HTTP Server
	core.Logger.Info(fmt.Sprintf("Starting HTTP server on port %s...", core.Configs.Server.Port))
	httpServer := &http.Server{
		Addr:           ":" + core.Configs.Server.Port,
		Handler:        container.Router.Engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server startup failed: %v", err)
	}

	return nil
}

func injectUserHandler() *rest.UserHandler {
	userRepository := data.NewUserRepository()
	userService := service.NewUserService(userRepository)
	return rest.NewUserHandler(userService)
}

func injectANYHandler() /**handler.ANYHandler*/ {
	// TODO: Implement ANY handler injection following the same pattern as userHandler
}
