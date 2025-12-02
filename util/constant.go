package util

const (
	DefaultLocale                   = "en"
	DefaultTimezone                 = "Asia/Ho_Chi_Minh"
	AppContextKey                   = "app_context"
	LocaleContextKey                = "locale"
	TraceIDContextKey               = "trace_id"
	LocaleMiddlewarePriority        = 0
	HTTPMiddlewarePriority          = 1
	JWTMiddlewarePriority           = 2
	TraceIDMiddlewarePriority       = 3
	ValidationMiddlewarePriority    = 4
	ErrorHandlingMiddlewarePriority = 5
)
