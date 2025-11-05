package middleware

import (
	"veg-store-backend/injection/core"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Writer.Header().Set("Access-Control-Allow-Origin", core.Configs.Cors.AllowOrigins)
		ginCtx.Writer.Header().Set("Access-Control-Allow-Methods", core.Configs.Cors.AllowMethods)
		ginCtx.Writer.Header().Set("Access-Control-Allow-Headers", core.Configs.Cors.AllowHeaders)

		if ginCtx.Request.Method == "OPTIONS" {
			ginCtx.AbortWithStatus(204)
			return
		}

		ginCtx.Next()
	}
}
