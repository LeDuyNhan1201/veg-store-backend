package validation

import "veg-store-backend/injection/core"

func InitValidationMessageKeys() map[string]string {
	var validationMessages = map[string]string{
		"email":    core.Error.Validation.Required.MessageKey,
		"required": core.Error.Validation.Required.MessageKey,
		"min":      core.Error.Validation.Min.MessageKey,
		"max":      core.Error.Validation.Max.MessageKey,
		"size":     core.Error.Validation.Size.MessageKey,
	}
	return validationMessages
}

func HandleParamForMessageKey(messageKey, field, param string) map[string]interface{} {
	params := make(map[string]interface{})
	switch messageKey {
	case core.Error.Validation.Min.MessageKey:
		params["Min"] = param
	case core.Error.Validation.Max.MessageKey:
		params["Max"] = param
	case core.Error.Validation.Size.MessageKey:
		params["Min"] = param
		params["Max"] = param
	}
	params["Field"] = field
	return params
}
