package injectiontest

import (
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/infrastructure/config"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/localizer"
	"veg-store-backend/internal/infrastructure/logger"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/api/middleware"
	"veg-store-backend/internal/api/resthandler"
)

func MockCore() *core.Core {
	return &core.Core{
		Error:     exception.Init(),
		Localizer: localizer.Init("test"),
		Logger:    logger.Init("test"),
		AppConfig: config.Init("test"),
	}
}

func InitTestHTTPRouter() *TestHTTPRouter {
	mockCore := MockCore()
	mockRouter := router.InitHTTPRouter(mockCore)
	middlewaresCollection := middleware.NewMiddlewaresCollection(
		middleware.NewLocaleMiddleware(mockCore, mockRouter),
		middleware.NewHTTPMiddleware(mockCore, mockRouter),
		middleware.NewJWTMiddleware(mockCore, mockRouter),
		middleware.NewTraceIDMiddleware(mockCore, mockRouter),
		middleware.NewValidationMiddleware(mockCore, mockRouter),
		middleware.NewErrorHandlingMiddleware(mockCore, mockRouter),
	)
	middlewaresCollection.Setup()
	return &TestHTTPRouter{mockRouter}
}

type TestHTTPRouter struct {
	*router.HTTPRouter
}

func (r *TestHTTPRouter) MockUserRoute(handler *resthandler.UserHandler) {
	group := r.AppGroup(r.ApiPath + "/users")
	{
		r.AppGET(group, "/hello", handler.Hello)
		r.AppGET(group, "/:id", handler.Details)
		r.AppGET(group, "", handler.GetAllUsers)
	}
}
