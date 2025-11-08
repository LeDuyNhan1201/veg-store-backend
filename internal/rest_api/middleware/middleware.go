package middleware

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"go.uber.org/fx"
)

type MiddlewaresCollection struct {
	Router     *router.Router
	JWTManager infra_interface.JWTManager
}

func NewMiddlewareCollections(
	router *router.Router,
	jwtManager infra_interface.JWTManager,
) MiddlewaresCollection {
	middlewaresCollection := MiddlewaresCollection{
		Router:     router,
		JWTManager: jwtManager,
	}

	// Register all middlewares
	// IMPORTANT: THE ORDER OF MIDDLEWARES MATTERS (FIRST IN, FIRST OUT)
	middlewaresCollection.Router.Engine.Use(
		Locale(util.DefaultLocale), // Locale middleware should be the first one to set the locale
		HttpContext(),              // HttpContext middleware should be after Locale to have access to the locale
		JWT(jwtManager),            // JWT middleware should be after HttpContext to have access to the HttpContext
		TraceID(),
		Validation(),   // Validation middleware only for binding and validating request data
		ErrorHandler(), // ErrorHandler middleware should be after all other middlewares to catch errors
	)
	return middlewaresCollection
}

var Module = fx.Options(fx.Invoke(NewMiddlewareCollections))
