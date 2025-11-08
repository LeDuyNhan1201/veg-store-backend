package exception

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
	Size     SubError
	Max      SubError
	Min      SubError
}

type AppError struct {
	NotFound   NotFoundError
	Auth       AuthError
	Invalid    InvalidError
	Validation ValidationError

	errorMap map[string]SubError
}

func (appError *AppError) FindByCode(code string) (SubError, bool) {
	err, ok := appError.errorMap[code]
	return err, ok
}

func InitAppError() *AppError {
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
			Size: SubError{
				Code:       "validation/forbidden",
				MessageKey: "Validation.Size",
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

	appError.buildErrorMap()
	return appError
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
		appError.Validation.Size.Code:      appError.Validation.Size,
		appError.Validation.Max.Code:       appError.Validation.Max,
		appError.Validation.Min.Code:       appError.Validation.Min,
	}
}
