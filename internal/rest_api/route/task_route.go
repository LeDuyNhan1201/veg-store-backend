package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/rest_handler"
)

type TaskRoute struct {
	*Route[*rest_handler.TaskHandler]
}

func NewTaskRoutes(authHandler *rest_handler.TaskHandler, router *router.HTTPRouter) *TaskRoute {
	return &TaskRoute{
		Route: &Route[*rest_handler.TaskHandler]{
			Handler: authHandler,
			Router:  router,
		},
	}
}

func (r *TaskRoute) Setup() {
	group := r.Router.AppGroup(r.Router.ApiPath + "/tasks")
	{
		r.Router.AppGET(group, "/search", r.Handler.Search)
		r.Router.AppGET(group, "", r.Handler.FindAll)
		r.Router.AppGET(group, "/statuses", r.Handler.FindAllStatuses)
	}
}
