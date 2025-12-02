package route

import (
	"veg-store-backend/internal/api/resthandler"
	"veg-store-backend/internal/infrastructure/router"
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
		r.Router.AppGET(group, "/statuses", r.Handler.FindAllStatuses)
		r.Router.AppPOST(group, "/search", r.Handler.Search)
		r.Router.AppGET(group, "", r.Handler.FindAll)
		r.Router.AppGET(group, "/:id", r.Handler.FindByID)
		r.Router.AppPOST(group, "", r.Handler.Create)
		r.Router.AppPUT(group, "/:id", r.Handler.Update)
		r.Router.AppPATCH(group, "", r.Handler.UpdateStatus)
	}
}
