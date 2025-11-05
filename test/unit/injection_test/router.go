package injection_test

import (
	"veg-store-backend/injection"
	"veg-store-backend/internal/api/rest"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *injection.Router {
	return InjectMock().Router
}

func MockUserRoutes(handler *rest.UserHandler) *gin.Engine {
	router := setupTestRouter()
	api := router.Engine.Group(router.ApiPath + "/user")
	{
		api.GET("/hello", func(ginCtx *gin.Context) {
			handler.Hello(MockHttpContext(ginCtx))
		})
		api.GET("/details/:id", func(ginCtx *gin.Context) {
			handler.Details(MockHttpContext(ginCtx))
		})
		api.GET("/ping", func(ginCtx *gin.Context) {
			handler.HealthCheck(MockHttpContext(ginCtx))
		})
		api.GET("/", func(ginCtx *gin.Context) {
			handler.GetAllUsers(MockHttpContext(ginCtx))
		})
	}
	return router.Engine
}
