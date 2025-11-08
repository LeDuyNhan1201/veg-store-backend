package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validation() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		core.Logger.Debug("[BEFORE] Validation invoked")
		ginContext.Next()
		core.Logger.Debug("[AFTER] Validation invoked")

		err := ginContext.Errors.Last()
		if err == nil {
			return
		}

		// Only handle validation errors
		var validationErrors validator.ValidationErrors
		if !errors.As(err.Err, &validationErrors) {
			return
		}

		httpContext := core.GetHttpContext(ginContext)
		var validationErrorDTO []dto.ValidationError
		if errors.As(err.Err, &validationErrors) {
			validationErrorDTO = make([]dto.ValidationError, 0)
			for _, fieldError := range validationErrors {
				field := handleField(httpContext.Locale(), fieldError.Field())
				errorParam := fieldError.Param()
				errorMessage := handleFieldError(httpContext.Locale(), field, fieldError.Tag(), errorParam)
				validationErrorDTO = append(validationErrorDTO, dto.ValidationError{
					Field: field,
					Error: errorMessage,
				})
			}

			ginContext.JSON(http.StatusBadRequest, dto.HttpResponse[[]dto.ValidationError]{
				HttpStatus: http.StatusBadRequest,
				Code:       core.Error.Invalid.Fields.Code,
				Message:    core.Translator.T(httpContext.Locale(), core.Error.Invalid.Fields.MessageKey),
				Data:       validationErrorDTO,
			})
			ginContext.Abort()
		}
	}
}

func handleField(locale, field string) string {
	messageKey := fmt.Sprintf("Field.%s", field)
	return core.Translator.T(locale, messageKey)
}

func handleFieldError(locale, field, errorKey string, params string) string {
	messageKey := core.ValidationMessageKeys[errorKey]
	return core.Translator.T(locale, messageKey, validation.HandleParamForMessageKey(messageKey, field, params))
}
