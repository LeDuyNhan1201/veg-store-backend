package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/rest_handler"
)

type AuthRoute struct {
	*Route[*rest_handler.AuthHandler]
}

func NewAuthRoutes(authHandler *rest_handler.AuthHandler, router *router.HTTPRouter) *AuthRoute {
	return &AuthRoute{
		Route: &Route[*rest_handler.AuthHandler]{
			Handler: authHandler,
			Router:  router,
		},
	}
}

func (route *AuthRoute) Setup() {
	group := route.Router.AppGroup(route.Router.ApiPath + "/auth")
	{
		route.Router.AppPOST(group, "/", route.Handler.SignIn)
		route.Router.AppGET(group, "/me", route.Handler.Info)
	}
}

//func (route *AuthRoute) Setup() {
//	api := route.Router.Engine.Group(route.Router.ApiPath + "/auth")
//	{
//		api.POST("/", func(ginContext *gin.Context) {
//			route.Handler.SignIn(context.GetHttpContext(ginContext))
//		})
//		api.GET("/me", func(ginContext *gin.Context) {
//			route.Handler.Info(context.GetHttpContext(ginContext))
//		})
//	}
//}
