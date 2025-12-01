package validation

import (
	"slices"

	"veg-store-backend/internal/domain/model"

	"github.com/go-playground/validator/v10"
)

func AllTaskFields() []string {
	return []string{
		string(model.FieldTitle),
		string(model.FieldStatusID),
		string(model.FieldStartDay),
		string(model.FieldTargetDay),
		string(model.FieldEndDay),
		string(model.FieldPriority),
		string(model.FieldCreatedAt),
		string(model.FieldCreatedBy),
		string(model.FieldUpdatedAt),
		string(model.FieldUpdatedBy),
		string(model.FieldIsDeleted),
		string(model.FieldVersion),
	}
}

var FieldsValidator validator.Func = func(fl validator.FieldLevel) bool {
	// Pass nil values
	value := fl.Field().String()
	if value == "" {
		return true
	}

	return slices.Contains(AllTaskFields(), fl.Field().String())
}
