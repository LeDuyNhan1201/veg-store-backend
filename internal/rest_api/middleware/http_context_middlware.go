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

func (m *HTTPMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		httpContext := &context.Http{
			Core: m.Core,
			Gin:  ginContext,
		}

		ginContext.Set(util.AppContextKey, httpContext)

		m.Logger.Debug("[BEFORE] Http invoked")
		ginContext.Next()
		m.Logger.Debug("[AFTER] Http invoked")
	}
}

func (m *HTTPMiddleware) Setup() {
	m.Router.Engine.Use(m.handler())
}

func (m *HTTPMiddleware) Priority() uint8 {
	return util.HTTPMiddlewarePriority
}
