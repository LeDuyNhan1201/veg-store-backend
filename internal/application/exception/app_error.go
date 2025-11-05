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
}

type AuthError struct {
	Unauthenticated SubError
	WrongPassword   SubError
	Forbidden       SubError
}

type AppError struct {
	NotFound NotFoundError
	Auth     AuthError
	Invalid  InvalidError

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
		appError.Auth.Unauthenticated.Code: appError.Auth.Unauthenticated,
		appError.Auth.WrongPassword.Code:   appError.Auth.WrongPassword,
		appError.Auth.Forbidden.Code:       appError.Auth.Forbidden,
	}
}
