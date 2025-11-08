package middleware

import (
	"errors"
	"net/http"
	"strings"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		core.Logger.Debug("[BEFORE] ErrorHandler invoked")
		ginContext.Next()
		core.Logger.Debug("[AFTER] ErrorHandler invoked")

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
				Message:    core.Translator.T(util.GetLocale(ginContext), subError.MessageKey),
			}

		default:
			core.Logger.Error("Unhandled error", zap.Error(rootError))
			ginContext.JSON(http.StatusInternalServerError, dto.HttpResponse[any]{
				HttpStatus: http.StatusInternalServerError,
				Code:       "internal/server-error",
				Message:    "Internal Server Error",
			})
			return
		}

		// Get trace_id
		traceID := util.GetTraceId(ginContext)

		core.Logger.Error("Request failed",
			zap.String("trace_id", traceID),
			zap.String("code", subError.Code),
			zap.String("message", subError.MessageKey),
			zap.String("path", ginContext.FullPath()),
			zap.String("method", ginContext.Request.Method),
		)

		ginContext.JSON(httpStatus, response)
		ginContext.Abort()
	}
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
