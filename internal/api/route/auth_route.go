package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/api/resthandler"
)

type AuthRoute struct {
	*Route[*resthandler.AuthHandler]
}

func NewAuthRoutes(authHandler *resthandler.AuthHandler, router *router.HTTPRouter) *AuthRoute {
	return &AuthRoute{
		Route: &Route[*resthandler.AuthHandler]{
			Handler: authHandler,
			Router:  router,
		},
	}
}

func (r *AuthRoute) Setup() {
	group := r.Router.AppGroup(r.Router.ApiPath + "/auth")
	{
		r.Router.AppPOST(group, "", r.Handler.SignIn)
		r.Router.AppGET(group, "/me", r.Handler.Info)
	}
}

//func (route *AuthRoute) Setup() {
//	api := route.Router.Engine.Group(route.Router.ApiPath + "/auth")
//	{
//		api.POST("/", func(ginContext *gin.Context) {
//			route.Handler.Search(context.GetHttpContext(ginContext))
//		})
//		api.GET("/me", func(ginContext *gin.Context) {
//			route.Handler.Info(context.GetHttpContext(ginContext))
//		})
//	}
//}
