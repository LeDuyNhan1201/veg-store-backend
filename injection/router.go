package injection

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/api/middleware"
	"veg-store-backend/internal/api/rest"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ApiPath string
	Engine  *gin.Engine
}

func (router *Router) Setup() {
	// Init Gin Engine
	engine := router.initGinEngine()

	// Register all middlewares
	engine.Use(
		middleware.CORS(),
		middleware.Locale(util.DefaultLocale),
		middleware.HttpContext(),
		middleware.TraceID(),
		middleware.ErrorHandler(),
	)

	router.ApiPath = core.Configs.Server.ApiPrefix + core.Configs.Server.ApiVersion
	router.Engine = engine
}

func (router *Router) initGinEngine() *gin.Engine {
	engine := gin.New()

	// Custom log for Gin per request
	core.UseGinRequestLogging(engine)

	// Register Custom recovery handler for Gin
	engine.Use(gin.CustomRecovery(func(ginCtx *gin.Context, recovered interface{}) {
		rest.CustomRecoveryHandler(core.GetHttpContext(ginCtx), recovered)
	}))

	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic("Failed to set trusted proxies" + err.Error())
	}

	return engine
}

func (router *Router) RegisterUserRoutes(userHandler *rest.UserHandler) {
	api := router.Engine.Group(router.ApiPath + "/user")
	{
		api.GET("/hello", func(ginContext *gin.Context) {
			userHandler.Hello(core.GetHttpContext(ginContext))
		})
		api.GET("/details/:id", func(ginContext *gin.Context) {
			userHandler.Details(core.GetHttpContext(ginContext))
		})
		api.GET("/ping", func(ginContext *gin.Context) {
			userHandler.HealthCheck(core.GetHttpContext(ginContext))
		})
		api.GET("/", func(ginContext *gin.Context) {
			userHandler.GetAllUsers(core.GetHttpContext(ginContext))
		})
	}
}
