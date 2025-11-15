package route

import (
	"veg-store-backend/internal/infrastructure/router"

	"go.uber.org/fx"
)

type RoutesCollection []IRoute

type Route[THandler any] struct {
	Handler THandler
	Router  *router.HTTPRouter
}

type IRoute interface {
	Setup()
}

// NewRoutesCollection IMPORTANT: REMEMBER TO ADD NEW ROUTES TO RoutesCollection
func NewRoutesCollection(
	userRoutes *UserRoute,
	authRoutes *AuthRoute,
) *RoutesCollection {
	return &RoutesCollection{
		userRoutes,
		authRoutes,
	}
}

func (routesCollection *RoutesCollection) Setup() {
	for _, routes := range *routesCollection {
		routes.Setup()
	}
}

// Module IMPORTANT: REMEMBER TO ADD NEW ROUTES TO Module
var Module = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewRoutesCollection),
)
