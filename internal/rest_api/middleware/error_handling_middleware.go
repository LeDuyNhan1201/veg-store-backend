package middleware

import (
	"errors"
	"net/http"
	"strings"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ErrorHandlingMiddleware struct {
	*Middleware
}

func NewErrorHandlingMiddleware(core *core.Core, router *router.HTTPRouter) *ErrorHandlingMiddleware {
	return &ErrorHandlingMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (middleware *ErrorHandlingMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		middleware.Logger.Debug("[BEFORE] ErrorHandler invoked")
		ginContext.Next()
		middleware.Logger.Debug("[AFTER] ErrorHandler invoked")

		if len(ginContext.Errors) == 0 {
			return
		}

		rootError := ginContext.Errors.Last().Err

		var validationErrors validator.ValidationErrors
		if errors.As(rootError, &validationErrors) {
			return
		}

		var subError exception.SubError
		var response dto.HttpResponse[any]
		var httpStatus int

		switch rootError.(type) {
		case exception.SubError:
			// Map code prefix -> HTTP httpStatus
			errors.As(rootError, &subError)
			httpStatus = mapErrorCodeToStatus(subError.Code)
			response = dto.HttpResponse[any]{
				HttpStatus: httpStatus,
				Code:       subError.Code,
				Message:    middleware.Localizer.T(util.GetLocale(ginContext), subError.MessageKey),
			}

		default:
			middleware.Logger.Error("Unhandled error", zap.Error(rootError))
			ginContext.AbortWithStatusJSON(http.StatusInternalServerError, dto.HttpResponse[any]{
				HttpStatus: http.StatusInternalServerError,
				Code:       "internal/server-error",
				Message:    "Internal Server Error",
			})
		}

		// Get trace_id
		traceID := util.GetTraceId(ginContext)

		middleware.Logger.Error("Request failed",
			zap.String("trace_id", traceID),
			zap.String("code", subError.Code),
			zap.String("message", subError.MessageKey),
			zap.String("path", ginContext.FullPath()),
			zap.String("method", ginContext.Request.Method),
		)

		ginContext.AbortWithStatusJSON(httpStatus, response)
	}
}

func (middleware *ErrorHandlingMiddleware) Setup() {
	middleware.Router.Engine.Use(middleware.handler())
}

func (middleware *ErrorHandlingMiddleware) Priority() uint8 {
	return util.ErrorHandlingMiddlewarePriority
}

func mapErrorCodeToStatus(code string) int {
	switch {
	case strings.HasPrefix(code, "invalid/"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "auth/unauthenticated"):
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "auth/forbidden"):
		return http.StatusForbidden
	case strings.HasPrefix(code, "not_found/"):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
