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

func (m *ValidationMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		m.Logger.Debug("[BEFORE] Validation invoked")
		ginContext.Next()
		m.Logger.Debug("[AFTER] Validation invoked")

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
				field := m.handleField(httpContext.Locale(), fieldError.Field())
				errorParam := fieldError.Param()
				errorMessage := m.handleFieldError(httpContext.Locale(), field, fieldError.Tag(), errorParam)
				validationErrorDTO = append(validationErrorDTO, dto.ValidationError{
					Field: field,
					Error: errorMessage,
				})
			}

			ginContext.AbortWithStatusJSON(http.StatusBadRequest, dto.HttpResponse[[]dto.ValidationError]{
				HttpStatus: http.StatusBadRequest,
				Code:       m.Error.Invalid.Fields.Code,
				Message:    m.Localizer.T(httpContext.Locale(), m.Error.Invalid.Fields.MessageKey),
				Data:       validationErrorDTO,
			})
		}
	}
}

func (m *ValidationMiddleware) Setup() {
	m.Router.Engine.Use(m.handler())
}

func (m *ValidationMiddleware) Priority() uint8 {
	return util.ValidationMiddlewarePriority
}

func (m *ValidationMiddleware) handleField(locale, field string) string {
	messageKey := fmt.Sprintf("Field.%s", field)
	return m.Localizer.T(locale, messageKey)
}

func (m *ValidationMiddleware) handleFieldError(locale, field, errorKey string, params string) string {
	messageKey := m.Error.ValidationMessages[errorKey]
	return m.Localizer.T(locale, messageKey, m.Error.HandleParamForMessageKey(messageKey, field, params))
}
