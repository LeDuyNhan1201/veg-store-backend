package injection_test

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/rest_api/middleware"
	"veg-store-backend/internal/rest_api/rest_handler"
	"veg-store-backend/test/identity_test"

	"github.com/gin-gonic/gin"
)

func MockGlobalComponents() {
	core.Configs.Mode = "test"
	core.Logger = core.InitLogger()       // Initialize Logger
	core.Configs = core.Load()            // Load configuration
	core.Translator = core.InitI18n()     // Initialize i18n Translator
	core.Error = exception.InitAppError() // Initialize App Error
}

func MockRouter() *router.Router {
	mockRouter := router.NewRouter()
	middleware.NewMiddleware(mockRouter, new(identity_test.MockJWTManager))
	return mockRouter
}

func MockUserRoutes(handler *rest_handler.UserHandler) *gin.Engine {
	mockRouter := MockRouter()
	api := mockRouter.Engine.Group(mockRouter.ApiPath + "/user")
	{
		api.GET("/hello", func(ginCtx *gin.Context) {
			handler.Hello(mockHttpContext(ginCtx))
		})
		api.GET("/details/:id", func(ginCtx *gin.Context) {
			handler.Details(mockHttpContext(ginCtx))
		})
		api.GET("/ping", func(ginCtx *gin.Context) {
			handler.HealthCheck(mockHttpContext(ginCtx))
		})
		api.GET("/", func(ginCtx *gin.Context) {
			handler.GetAllUsers(mockHttpContext(ginCtx))
		})
	}
	return mockRouter.Engine
}

func mockHttpContext(
	ginCtx *gin.Context,
) *core.HttpContext {
	return &core.HttpContext{
		Translator: core.Translator,
		Gin:        ginCtx,
	}
}
