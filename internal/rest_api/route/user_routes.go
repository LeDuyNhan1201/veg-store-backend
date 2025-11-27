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

func (r *UserRoute) Setup() {
	group := r.Router.AppGroup(r.Router.ApiPath + "/users")
	{
		r.Router.AppGET(group, "/hello", r.Handler.Hello)
		r.Router.AppGET(group, "/:id", r.Handler.Details)
		r.Router.AppGET(group, "", r.Handler.GetAllUsers)
	}
}
