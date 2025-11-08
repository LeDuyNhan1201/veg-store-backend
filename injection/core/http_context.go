package core

import (
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

/*
This file defines the HttpContext struct, which holds shared resources for handling requests.
It includes:
- Translator: For handling localization and translations.
- Gin: The current Gin context for the request.
*/

type HttpContext struct {
	Translator      *Localizer
	Gin             *gin.Context
	SecurityContext *SecurityContext
}

// GetHttpContext - Usage: httpContext := core.GetHttpContext(c) to get the HttpContext in a rest_handler.
func GetHttpContext(ginContext *gin.Context) *HttpContext {
	val, exists := ginContext.Get(util.AppContextKey)
	if !exists {
		panic("HttpContext not found in Gin context — did you forget to register AppContextMiddleware?")
	}

	httpContext, ok := val.(*HttpContext)
	if !ok {
		panic("HttpContext type assertion failed")
	}
	return httpContext

	//return &HttpContext{
	//	Translator:      Translator,
	//	Gin:             ginContext,
	//	SecurityContext: nil,
	//}
}

// JSON - httpContext.JSON(statusCode, responseObject) to send a JSON response.
func (httpContext *HttpContext) JSON(status int, obj any) {
	httpContext.Gin.JSON(status, obj)
}

// T - Usage: httpContext.T("message_id", params) to get a localized message.
func (httpContext *HttpContext) T(msgID string, params ...map[string]interface{}) string {
	if httpContext.Translator == nil {
		return msgID
	}
	return httpContext.Translator.Localize(httpContext.Locale(), msgID, params...)
}

// Locale - Usage: locale := httpContext.Locale() to retrieve the locale string.
func (httpContext *HttpContext) Locale() string {
	if v, ok := httpContext.Gin.Get(util.LocaleContextKey); ok {
		return v.(string)
	}
	return "en"
}

// SetSecurityContext - Usage: httpContext.SetSecurityContext(securityContext) to set the SecurityContext.
func (httpContext *HttpContext) SetSecurityContext(securityContext *SecurityContext) {
	httpContext.SecurityContext = securityContext
}
