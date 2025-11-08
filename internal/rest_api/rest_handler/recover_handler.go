package rest_handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

/*
This is a global exception rest_handler for error from panic()
*/

func CustomRecoveryHandler(httpContext *core.HttpContext, recovered interface{}) {
	traceID := httpContext.Gin.GetString(util.TraceIDContextKey) // Require middleware to attach trace ID
	stack := string(debug.Stack())
	core.Logger.Warn(fmt.Sprintf("[PANIC] trace_id=%s error=%v stack trace:\n%s", traceID, recovered, stack))

	httpContext.JSON(http.StatusInternalServerError, dto.HttpResponse[any]{
		HttpStatus: http.StatusInternalServerError,
		Code:       "internal/server-error",
		Message:    "Internal Server Error",
	})

	httpContext.JSON(http.StatusInternalServerError, gin.H{
		"error":                "internal server error",
		util.TraceIDContextKey: traceID,
	})
	httpContext.Gin.Abort()
}
