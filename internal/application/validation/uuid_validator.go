package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var UUIDValidator validator.Func = func(fl validator.FieldLevel) bool {
	// Pass nil values
	value := fl.Field().String()
	if value == "" {
		return true
	}

	param := fl.Param()
	if param == "" {
		return false
	}
	_, err := uuid.Parse(param)
	return err == nil
}
