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

func NewMiddleware(
	router *router.Router,
	jwtManager infra_interface.JWTManager,
) MiddlewaresCollection {
	middlewaresCollection := MiddlewaresCollection{
		Router:     router,
		JWTManager: jwtManager,
	}

	// Register all middlewares
	middlewaresCollection.Router.Engine.Use(
		Locale(util.DefaultLocale),
		//HttpContext(),
		TraceID(),
		ErrorHandler(),
		JWT(jwtManager),
	)
	return middlewaresCollection
}

var Module = fx.Options(fx.Invoke(NewMiddleware))
