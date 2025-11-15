package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/rest_handler"
)

type UserRoute struct {
	*Route[*rest_handler.UserHandler]
}

func NewUserRoutes(userHandler *rest_handler.UserHandler, router *router.HTTPRouter) *UserRoute {
	return &UserRoute{
		Route: &Route[*rest_handler.UserHandler]{
			Handler: userHandler,
			Router:  router,
		},
	}
}

func (route *UserRoute) Setup() {
	group := route.Router.AppGroup(route.Router.ApiPath + "/user")
	{
		route.Router.AppGET(group, "/hello", route.Handler.Hello)
		route.Router.AppGET(group, "/details/:id", route.Handler.Details)
		route.Router.AppGET(group, "/", route.Handler.GetAllUsers)
	}
}
