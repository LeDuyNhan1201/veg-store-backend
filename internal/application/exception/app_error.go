package exception

import (
	"veg-store-backend/internal/application/validation"
)

type SubError struct {
	Code       string
	MessageKey string
	Args       []map[string]interface{}
}

func (e *SubError) Error() string {
	return e.Code
}

func (e *SubError) MoreInfo(Args ...map[string]interface{}) *SubError {
	e.Args = append(e.Args, Args...)
	return e
}

type NotFoundError struct {
	User    *SubError
	Product *SubError
	Task    *SubError
}

type FailError struct {
	CreateUser    *SubError
	CreateProduct *SubError
	CreateTask    *SubError
	UpdateTask    *SubError
	DeleteTask    *SubError
}

type InvalidError struct {
	Token    *SubError
	Email    *SubError
	Username *SubError
	Fields   *SubError
}

type AuthError struct {
	Unauthenticated *SubError
	WrongPassword   *SubError
	Forbidden       *SubError
}

type ValidationError struct {
	Email    *SubError
	Required *SubError
	Range    *SubError
	Max      *SubError
	Min      *SubError
}

type AppError struct {
	NotFound   NotFoundError
	Fail       FailError
	Auth       AuthError
	Invalid    InvalidError
	Validation ValidationError

	ValidationMessages map[string]string

	errorMap map[string]*SubError
}

func (e *AppError) FindByCode(code string) (*SubError, bool) {
	err, ok := e.errorMap[code]
	return err, ok
}

func Init() *AppError {
	appError := &AppError{
		NotFound: NotFoundError{
			User: &SubError{
				Code:       "not_found/user",
				MessageKey: "NotFound.User",
			},
			Product: &SubError{
				Code:       "not_found/product",
				MessageKey: "NotFound.Product",
			},
			Task: &SubError{
				Code:       "not_found/task",
				MessageKey: "NotFound.Task",
			},
		},
		Fail: FailError{
			CreateUser: &SubError{
				Code:       "fail/create_user",
				MessageKey: "Fail.User",
			},
			CreateProduct: &SubError{
				Code:       "fail/create_product",
				MessageKey: "Fail.CreateProduct",
			},
			CreateTask: &SubError{
				Code:       "fail/create_task",
				MessageKey: "Fail.CreateTask",
			},
			UpdateTask: &SubError{
				Code:       "fail/update_task",
				MessageKey: "Fail.UpdateTask",
			},
			DeleteTask: &SubError{
				Code:       "fail/delete_task",
				MessageKey: "Fail.DeleteTask",
			},
		},
		Invalid: InvalidError{
			Token: &SubError{
				Code:       "invalid/token",
				MessageKey: "Invalid.Token",
			},
			Email: &SubError{
				Code:       "invalid/email",
				MessageKey: "Invalid.Email",
			},
			Username: &SubError{
				Code:       "invalid/username",
				MessageKey: "Invalid.Username",
			},
			Fields: &SubError{
				Code:       "invalid/fields",
				MessageKey: "Invalid.Fields",
			},
		},
		Auth: AuthError{
			Unauthenticated: &SubError{
				Code:       "auth/unauthenticated",
				MessageKey: "Auth.Unauthenticated",
			},
			WrongPassword: &SubError{
				Code:       "auth/wrong-password",
				MessageKey: "Auth.WrongPassword",
			},
			Forbidden: &SubError{
				Code:       "auth/forbidden",
				MessageKey: "Auth.Forbidden",
			},
		},
		Validation: ValidationError{
			Email: &SubError{
				Code:       "validation/email",
				MessageKey: "Validation.Email",
			},
			Required: &SubError{
				Code:       "validation/required",
				MessageKey: "Validation.Required",
			},
			Range: &SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Range",
			},
			Max: &SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Max",
			},
			Min: &SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Min",
			},
		},
	}
	appError.initValidationMessageKeys()
	appError.buildErrorMap()
	return appError
}

func (e *AppError) initValidationMessageKeys() {
	var validationMessages = map[string]string{
		"email":    e.Validation.Required.MessageKey,
		"required": e.Validation.Required.MessageKey,
		"min":      e.Validation.Min.MessageKey,
		"max":      e.Validation.Max.MessageKey,
		"range":    e.Validation.Range.MessageKey,
	}
	e.ValidationMessages = validationMessages
}

func (e *AppError) HandleParamForMessageKey(messageKey, field, param string) map[string]interface{} {
	params := make(map[string]interface{})
	switch messageKey {
	case e.Validation.Min.MessageKey:
		params["Min"] = param
	case e.Validation.Max.MessageKey:
		params["Max"] = param
	case e.Validation.Range.MessageKey:
		minParam, maxParam := validation.ParseRangeParam(param)
		params["Min"] = minParam
		params["Max"] = maxParam
	}
	params["Field"] = field
	return params
}

func (e *AppError) buildErrorMap() {
	e.errorMap = map[string]*SubError{
		e.NotFound.User.Code:        e.NotFound.User,
		e.NotFound.Product.Code:     e.NotFound.Product,
		e.Invalid.Token.Code:        e.Invalid.Token,
		e.Invalid.Email.Code:        e.Invalid.Email,
		e.Invalid.Username.Code:     e.Invalid.Username,
		e.Invalid.Fields.Code:       e.Invalid.Fields,
		e.Auth.Unauthenticated.Code: e.Auth.Unauthenticated,
		e.Auth.WrongPassword.Code:   e.Auth.WrongPassword,
		e.Auth.Forbidden.Code:       e.Auth.Forbidden,
		e.Validation.Email.Code:     e.Validation.Email,
		e.Validation.Required.Code:  e.Validation.Required,
		e.Validation.Range.Code:     e.Validation.Range,
		e.Validation.Max.Code:       e.Validation.Max,
		e.Validation.Min.Code:       e.Validation.Min,
	}
}
