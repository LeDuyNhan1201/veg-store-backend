package exception

import (
	"veg-store-backend/internal/application/validation"
)

type SubError struct {
	Code       string
	MessageKey string
}

func (subError SubError) Error() string {
	return subError.Code
}

type NotFoundError struct {
	User    SubError
	Product SubError
}

type InvalidError struct {
	Token    SubError
	Email    SubError
	Username SubError
	Fields   SubError
}

type AuthError struct {
	Unauthenticated SubError
	WrongPassword   SubError
	Forbidden       SubError
}

type ValidationError struct {
	Email    SubError
	Required SubError
	Range    SubError
	Max      SubError
	Min      SubError
}

type AppError struct {
	NotFound   NotFoundError
	Auth       AuthError
	Invalid    InvalidError
	Validation ValidationError

	ValidationMessages map[string]string

	errorMap map[string]SubError
}

func (appError *AppError) FindByCode(code string) (SubError, bool) {
	err, ok := appError.errorMap[code]
	return err, ok
}

func Init() *AppError {
	appError := &AppError{
		NotFound: NotFoundError{
			User: SubError{
				Code:       "not_found/user",
				MessageKey: "NotFound.User",
			},
			Product: SubError{
				Code:       "not_found/product",
				MessageKey: "NotFound.Product",
			},
		},
		Invalid: InvalidError{
			Token: SubError{
				Code:       "invalid/token",
				MessageKey: "Invalid.Token",
			},
			Email: SubError{
				Code:       "invalid/email",
				MessageKey: "Invalid.Email",
			},
			Username: SubError{
				Code:       "invalid/username",
				MessageKey: "Invalid.Username",
			},
			Fields: SubError{
				Code:       "invalid/fields",
				MessageKey: "Invalid.Fields",
			},
		},
		Auth: AuthError{
			Unauthenticated: SubError{
				Code:       "auth/unauthenticated",
				MessageKey: "Auth.Unauthenticated",
			},
			WrongPassword: SubError{
				Code:       "auth/wrong-password",
				MessageKey: "Auth.WrongPassword",
			},
			Forbidden: SubError{
				Code:       "auth/forbidden",
				MessageKey: "Auth.Forbidden",
			},
		},
		Validation: ValidationError{
			Email: SubError{
				Code:       "validation/email",
				MessageKey: "Validation.Email",
			},
			Required: SubError{
				Code:       "validation/required",
				MessageKey: "Validation.Required",
			},
			Range: SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Range",
			},
			Max: SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Max",
			},
			Min: SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Min",
			},
		},
	}
	appError.initValidationMessageKeys()
	appError.buildErrorMap()
	return appError
}

func (appError *AppError) initValidationMessageKeys() {
	var validationMessages = map[string]string{
		"email":    appError.Validation.Required.MessageKey,
		"required": appError.Validation.Required.MessageKey,
		"min":      appError.Validation.Min.MessageKey,
		"max":      appError.Validation.Max.MessageKey,
		"range":    appError.Validation.Range.MessageKey,
	}
	appError.ValidationMessages = validationMessages
}

func (appError *AppError) HandleParamForMessageKey(messageKey, field, param string) map[string]interface{} {
	params := make(map[string]interface{})
	switch messageKey {
	case appError.Validation.Min.MessageKey:
		params["Min"] = param
	case appError.Validation.Max.MessageKey:
		params["Max"] = param
	case appError.Validation.Range.MessageKey:
		minParam, maxParam := validation.ParseRangeParam(param)
		params["Min"] = minParam
		params["Max"] = maxParam
	}
	params["Field"] = field
	return params
}

func (appError *AppError) buildErrorMap() {
	appError.errorMap = map[string]SubError{
		appError.NotFound.User.Code:        appError.NotFound.User,
		appError.NotFound.Product.Code:     appError.NotFound.Product,
		appError.Invalid.Token.Code:        appError.Invalid.Token,
		appError.Invalid.Email.Code:        appError.Invalid.Email,
		appError.Invalid.Username.Code:     appError.Invalid.Username,
		appError.Invalid.Fields.Code:       appError.Invalid.Fields,
		appError.Auth.Unauthenticated.Code: appError.Auth.Unauthenticated,
		appError.Auth.WrongPassword.Code:   appError.Auth.WrongPassword,
		appError.Auth.Forbidden.Code:       appError.Auth.Forbidden,
		appError.Validation.Email.Code:     appError.Validation.Email,
		appError.Validation.Required.Code:  appError.Validation.Required,
		appError.Validation.Range.Code:     appError.Validation.Range,
		appError.Validation.Max.Code:       appError.Validation.Max,
		appError.Validation.Min.Code:       appError.Validation.Min,
	}
}
