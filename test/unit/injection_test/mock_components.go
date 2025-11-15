package injection_test

import (
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/infrastructure/config"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/localizer"
	"veg-store-backend/internal/infrastructure/logger"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/middleware"
	"veg-store-backend/internal/rest_api/rest_handler"
)

func MockCore() *core.Core {
	return &core.Core{
		Error:     exception.Init(),
		Localizer: localizer.Init("test"),
		Logger:    logger.Init("test"),
		Config:    config.Init("test"),
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

func (router *TestHTTPRouter) MockUserRoute(handler *rest_handler.UserHandler) {
	group := router.AppGroup(router.ApiPath + "/user")
	{
		router.AppGET(group, "/hello", handler.Hello)
		router.AppGET(group, "/details/:id", handler.Details)
		router.AppGET(group, "/", handler.GetAllUsers)
	}
}
