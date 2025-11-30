package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/api/resthandler"
)

type UserRoute struct {
	*Route[*resthandler.UserHandler]
}

func NewUserRoutes(userHandler *resthandler.UserHandler, router *router.HTTPRouter) *UserRoute {
	return &UserRoute{
		Route: &Route[*resthandler.UserHandler]{
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
