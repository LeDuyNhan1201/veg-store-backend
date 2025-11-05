package injection

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/api/middleware"
	"veg-store-backend/internal/api/rest"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// @title Example API
// @version 1.0
// @description This is a sample server that uses JWT authentication.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

// @schemes http

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
	engine.Use(gin.CustomRecovery(func(ginContext *gin.Context, recovered interface{}) {
		rest.CustomRecoveryHandler(core.GetHttpContext(ginContext), recovered)
	}))

	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic("Failed to set trusted proxies" + err.Error())
	}

	return engine
}

func (router *Router) RegisterSwaggerUI() {
	docs.SwaggerInfo.BasePath = router.ApiPath
	router.Engine.GET(router.ApiPath+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) /*,
	ginSwagger.URL("http://localhost:2345"+router.ApiPath+"/swagger/doc.json"),
	ginSwagger.DefaultModelsExpandDepth(-1)*/
}

func (router *Router) RegisterUserRoutes(userHandler *rest.UserHandler) {
	core.Logger.Info("Api Path: " + router.ApiPath)
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
