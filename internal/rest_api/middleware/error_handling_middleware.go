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

func (m *ErrorHandlingMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		m.Logger.Debug("[BEFORE] ErrorHandler invoked")
		ginContext.Next()
		m.Logger.Debug("[AFTER] ErrorHandler invoked")

		if len(ginContext.Errors) == 0 {
			return
		}

		rootError := ginContext.Errors.Last().Err

		var validationErrors validator.ValidationErrors
		if errors.As(rootError, &validationErrors) {
			return
		}

		var subError *exception.SubError
		var response dto.HttpResponse[any]
		var httpStatus int

		if errors.As(rootError, &subError) {
			// Map code prefix -> HTTP httpStatus
			httpStatus = mapErrorCodeToStatus(subError.Code)
			response = dto.HttpResponse[any]{
				HttpStatus: httpStatus,
				Code:       subError.Code,
				Message:    m.Localizer.T(util.GetLocale(ginContext), subError.MessageKey, subError.Args...),
			}
		} else {
			m.Logger.Error("Unhandled error", zap.Error(rootError))
			ginContext.AbortWithStatusJSON(http.StatusInternalServerError, dto.HttpResponse[any]{
				HttpStatus: http.StatusInternalServerError,
				Code:       "internal/server-error",
				Message:    "Internal Server Error",
			})
		}

		// Get trace_id
		traceID := util.GetTraceId(ginContext)

		m.Logger.Error("Request failed",
			zap.String("trace_id", traceID),
			zap.String("code", subError.Code),
			zap.String("message", subError.MessageKey),
			zap.String("path", ginContext.FullPath()),
			zap.String("method", ginContext.Request.Method),
		)

		ginContext.AbortWithStatusJSON(httpStatus, response)
	}
}

func (m *ErrorHandlingMiddleware) Setup() {
	m.Router.Engine.Use(m.handler())
}

func (m *ErrorHandlingMiddleware) Priority() uint8 {
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
	case strings.HasPrefix(code, "fail/create_"):
		return http.StatusNotExtended
	case strings.HasPrefix(code, "fail/update_"):
		return http.StatusNotModified
	case strings.HasPrefix(code, "fail/delete_"):
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}
