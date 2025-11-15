package middleware

import (
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TraceIDMiddleware struct {
	*Middleware
}

func NewTraceIDMiddleware(core *core.Core, router *router.HTTPRouter) *TraceIDMiddleware {
	return &TraceIDMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (middleware *TraceIDMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		traceID := ginContext.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		ginContext.Set(util.TraceIDContextKey, traceID)
		ginContext.Writer.Header().Set("X-Request-ID", traceID)

		middleware.Logger.Debug("[BEFORE] TraceID invoked")
		ginContext.Next()
		middleware.Logger.Debug("[AFTER] TraceID invoked")
	}
}

func (middleware *TraceIDMiddleware) Setup() {
	middleware.Router.Engine.Use(middleware.handler())
}

func (middleware *TraceIDMiddleware) Priority() uint8 {
	return util.TraceIDMiddlewarePriority
}
