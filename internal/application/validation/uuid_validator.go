package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var UUIDValidator validator.Func = func(fieldLevel validator.FieldLevel) bool {
	param := fieldLevel.Param()
	if param == "" {
		return false
	}
	_, err := uuid.Parse(param)
	if err != nil {
		return false
	}
	return true
}
