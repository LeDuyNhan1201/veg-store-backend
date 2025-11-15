package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationMiddleware struct {
	*Middleware
}

func NewValidationMiddleware(core *core.Core, router *router.HTTPRouter) *ValidationMiddleware {
	return &ValidationMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (middleware *ValidationMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		middleware.Logger.Debug("[BEFORE] Validation invoked")
		ginContext.Next()
		middleware.Logger.Debug("[AFTER] Validation invoked")

		err := ginContext.Errors.Last()
		if err == nil {
			return
		}

		// Only handle validation errors
		var validationErrors validator.ValidationErrors
		if !errors.As(err.Err, &validationErrors) {
			return
		}

		httpContext := context.GetHttpContext(ginContext)
		var validationErrorDTO []dto.ValidationError
		if errors.As(err.Err, &validationErrors) {
			validationErrorDTO = make([]dto.ValidationError, 0)
			for _, fieldError := range validationErrors {
				field := middleware.handleField(httpContext.Locale(), fieldError.Field())
				errorParam := fieldError.Param()
				errorMessage := middleware.handleFieldError(httpContext.Locale(), field, fieldError.Tag(), errorParam)
				validationErrorDTO = append(validationErrorDTO, dto.ValidationError{
					Field: field,
					Error: errorMessage,
				})
			}

			ginContext.AbortWithStatusJSON(http.StatusBadRequest, dto.HttpResponse[[]dto.ValidationError]{
				HttpStatus: http.StatusBadRequest,
				Code:       middleware.Error.Invalid.Fields.Code,
				Message:    middleware.Localizer.T(httpContext.Locale(), middleware.Error.Invalid.Fields.MessageKey),
				Data:       validationErrorDTO,
			})
		}
	}
}

func (middleware *ValidationMiddleware) Setup() {
	middleware.Router.Engine.Use(middleware.handler())
}

func (middleware *ValidationMiddleware) Priority() uint8 {
	return util.ValidationMiddlewarePriority
}

func (middleware *ValidationMiddleware) handleField(locale, field string) string {
	messageKey := fmt.Sprintf("Field.%s", field)
	return middleware.Localizer.T(locale, messageKey)
}

func (middleware *ValidationMiddleware) handleFieldError(locale, field, errorKey string, params string) string {
	messageKey := middleware.Error.ValidationMessages[errorKey]
	return middleware.Localizer.T(locale, messageKey, middleware.Error.HandleParamForMessageKey(messageKey, field, params))
}
