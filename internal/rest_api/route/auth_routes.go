package route

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/rest_handler"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	*Route[*rest_handler.AuthHandler]
}

func NewAuthRoutes(authHandler *rest_handler.AuthHandler, router *router.Router) *AuthRoutes {
	return &AuthRoutes{
		Route: &Route[*rest_handler.AuthHandler]{
			Handler: authHandler,
			Router:  router,
		},
	}
}

func (routes *AuthRoutes) Setup() {
	api := routes.Router.Engine.Group(routes.Router.ApiPath + "/auth")
	{
		api.POST("/sign-in", func(ginContext *gin.Context) {
			routes.Handler.SignIn(core.GetHttpContext(ginContext))
		})
	}
}
