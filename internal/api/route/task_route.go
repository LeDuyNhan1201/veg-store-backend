package route

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/api/resthandler"
)

type TaskRoute struct {
	*Route[*resthandler.TaskHandler]
}

func NewTaskRoutes(authHandler *resthandler.TaskHandler, router *router.HTTPRouter) *TaskRoute {
	return &TaskRoute{
		Route: &Route[*resthandler.TaskHandler]{
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
