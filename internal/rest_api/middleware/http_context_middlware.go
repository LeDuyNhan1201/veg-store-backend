package middleware

import (
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

type HTTPMiddleware struct {
	*Middleware
}

func NewHTTPMiddleware(core *core.Core, router *router.HTTPRouter) *HTTPMiddleware {
	return &HTTPMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (middleware *HTTPMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		httpContext := &context.Http{
			Core: middleware.Core,
			Gin:  ginContext,
		}

		ginContext.Set(util.AppContextKey, httpContext)

		middleware.Logger.Debug("[BEFORE] Http invoked")
		ginContext.Next()
		middleware.Logger.Debug("[AFTER] Http invoked")
	}
}

func (middleware *HTTPMiddleware) Setup() {
	middleware.Router.Engine.Use(middleware.handler())
}

func (middleware *HTTPMiddleware) Priority() uint8 {
	return util.HTTPMiddlewarePriority
}
