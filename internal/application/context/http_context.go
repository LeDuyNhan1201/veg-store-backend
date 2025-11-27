package context

import (
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

/*
This file defines the Http struct, which holds shared resources for handling requests.
It includes:
- Localizer: For handling localization and translations.
- Gin: The current Gin context for the request.
*/

type Http struct {
	*core.Core
	Gin             *gin.Context
	SecurityContext *SecurityContext
}

// GetHttpContext - Usage: httpContext := core.GetHttpContext(c) to get the Http in a rest_handler.
func GetHttpContext(ginContext *gin.Context) *Http {
	val, exists := ginContext.Get(util.AppContextKey)
	if !exists {
		panic("Http context not found in Gin context — did you forget to register HTTPMiddleware?")
	}

	httpContext, ok := val.(*Http)
	if !ok {
		panic("Http type assertion failed")
	}
	return httpContext

	//return &Http{
	//	Localizer:      Localizer,
	//	Gin:             ginContext,
	//	SecurityContext: nil,
	//}
}

// JSON - httpContext.JSON(statusCode, responseObject) to send a JSON response.
func (ctx *Http) JSON(status int, obj any) {
	ctx.Gin.JSON(status, obj)
}

// T - Usage: httpContext.T("message_id", params) to get a localized message.
func (ctx *Http) T(msgID string, params ...map[string]interface{}) string {
	if ctx.Localizer == nil {
		return msgID
	}
	return ctx.Localizer.Localize(ctx.Locale(), msgID, params...)
}

// Locale - Usage: locale := httpContext.Locale() to retrieve the locale string.
func (ctx *Http) Locale() string {
	if v, ok := ctx.Gin.Get(util.LocaleContextKey); ok {
		return v.(string)
	}
	return "en"
}

// SetSecurityContext - Usage: httpContext.SetSecurityContext(securityContext) to set the SecurityContext.
func (ctx *Http) SetSecurityContext(securityContext *SecurityContext) {
	ctx.SecurityContext = securityContext
}
