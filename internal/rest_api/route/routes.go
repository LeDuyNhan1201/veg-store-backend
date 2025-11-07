package route

import (
	"veg-store-backend/internal/infrastructure/router"

	"go.uber.org/fx"
)

type RoutesCollection []Routes

type Route[THandler any] struct {
	Handler THandler
	Router  *router.Router
}

type Routes interface {
	Setup()
}

// NewRoutesCollection IMPORTANT: REMEMBER TO ADD NEW ROUTES TO RoutesCollection
func NewRoutesCollection(
	userRoutes *UserRoutes,
	authRoutes *AuthRoutes,
) RoutesCollection {
	return RoutesCollection{
		userRoutes,
		authRoutes,
	}
}

func (routesCollection RoutesCollection) Setup() {
	for _, routes := range routesCollection {
		routes.Setup()
	}
}

// RoutesModule IMPORTANT: REMEMBER TO ADD NEW ROUTES TO RoutesModule
var RoutesModule = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewRoutesCollection),
)
