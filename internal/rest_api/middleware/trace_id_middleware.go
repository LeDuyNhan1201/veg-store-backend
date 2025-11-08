package middleware

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceID() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		traceID := ginContext.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		ginContext.Set(util.TraceIDContextKey, traceID)
		ginContext.Writer.Header().Set("X-Request-ID", traceID)

		core.Logger.Debug("[BEFORE] TraceID invoked")
		ginContext.Next()
		core.Logger.Debug("[AFTER] TraceID invoked")
	}
}
