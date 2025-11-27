package router

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"veg-store-backend/docs"
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/validation"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/logger"
	"veg-store-backend/internal/rest_api/rest_handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Summary Health Check
// @Description Check if the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @HTTPRouter /heath [get]

type HTTPRouter struct {
	*core.Core
	*gin.Engine
	ApiPath string
}

type HttpHandlerFunc func(context *context.Http)

func InitHTTPRouter(core *core.Core) *HTTPRouter {
	router := &HTTPRouter{
		Core:    core,
		Engine:  initGinEngine(core),
		ApiPath: core.AppConfig.Server.ApiPrefix + core.AppConfig.Server.ApiVersion,
	}

	router.Engine.GET(router.ApiPath+"/heath", func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	if core.AppConfig.Mode != "prod" && core.AppConfig.Mode != "production" {
		// Register Swagger UI in non-production modes
		router.registerSwaggerUI()
	}

	return router
}

func (r *HTTPRouter) HttpRun() error {
	// Run HTTP Server
	r.Logger.Info(fmt.Sprintf("Starting HTTP server on port %s...", r.AppConfig.Server.Port))
	httpServer := &http.Server{
		Addr:           ":" + r.AppConfig.Server.Port,
		Handler:        r.Engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server startup failed %w", err)
	}

	return nil
}

func initGinEngine(core *core.Core) *gin.Engine {
	engine := gin.New()

	// Custom log for Gin per request
	logger.UseGinRequestLogging(core.AppConfig.Mode, engine)

	// Register Custom recovery rest_handler for Gin
	engine.Use(gin.CustomRecovery(func(ginContext *gin.Context, recovered interface{}) {
		rest_handler.NewRecoveryHandler(core).Setup(context.GetHttpContext(ginContext), recovered)
	}))

	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic("Failed to set trusted proxies" + err.Error())
	}

	// Configure CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     core.AppConfig.Cors.AllowOrigins,
		AllowMethods:     core.AppConfig.Cors.AllowMethods,
		AllowHeaders:     core.AppConfig.Cors.AllowHeaders,
		AllowCredentials: core.AppConfig.Cors.AllowCredentials,
	}))

	validation.Init()

	return engine
}

func (r *HTTPRouter) registerSwaggerUI() {
	docs.SwaggerInfo.Host = r.AppConfig.Swagger.Host
	docs.SwaggerInfo.BasePath = r.ApiPath

	swaggerUiPrefix := docs.SwaggerInfo.BasePath + "/swagger-ui/*any"
	r.Engine.GET(swaggerUiPrefix, ginSwagger.WrapHandler(swaggerFiles.Handler)) /*,
	ginSwagger.URL("http://localhost:8080"+r.ApiPath+"/swagger-ui/doc.json"),
	ginSwagger.DefaultModelsExpandDepth(1)*/
}

// ====================================== //
// ====== CUSTOMIZE GIN FUNCTIONS ======= //
// ====================================== //

func (r *HTTPRouter) AppUse(middlewares ...HttpHandlerFunc) gin.IRoutes {
	return r.Use(adaptHandlers(middlewares...)...)
}

func (r *HTTPRouter) AppGroup(relativePath string, handlers ...HttpHandlerFunc) *gin.RouterGroup {
	return r.Group(relativePath, adaptHandlers(handlers...)...)
}

func (r *HTTPRouter) AppPOST(group *gin.RouterGroup, relativePath string, handlers ...HttpHandlerFunc) gin.IRoutes {
	return group.POST(relativePath, adaptHandlers(handlers...)...)
}

func (r *HTTPRouter) AppGET(group *gin.RouterGroup, relativePath string, handlers ...HttpHandlerFunc) gin.IRoutes {
	return group.GET(relativePath, adaptHandlers(handlers...)...)
}

func (r *HTTPRouter) AppDELETE(group *gin.RouterGroup, relativePath string, handlers ...HttpHandlerFunc) gin.IRoutes {
	return group.DELETE(relativePath, adaptHandlers(handlers...)...)
}

func (r *HTTPRouter) AppPATCH(group *gin.RouterGroup, relativePath string, handlers ...HttpHandlerFunc) gin.IRoutes {
	return group.PATCH(relativePath, adaptHandlers(handlers...)...)
}

func (r *HTTPRouter) AppPUT(group *gin.RouterGroup, relativePath string, handlers ...HttpHandlerFunc) gin.IRoutes {
	return group.PUT(relativePath, adaptHandlers(handlers...)...)
}

func adaptHandler(handler HttpHandlerFunc) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		httpContext := context.GetHttpContext(ginContext)
		handler(httpContext)
	}
}

func adaptHandlers(handlers ...HttpHandlerFunc) []gin.HandlerFunc {
	out := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		out[i] = adaptHandler(handler)
	}
	return out
}
