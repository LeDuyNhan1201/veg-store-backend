package middleware

import (
	"net/http"
	"strings"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"

	"github.com/gin-gonic/gin"
)

func JWT(jwtManager infra_interface.JWTManager) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		httpContext := core.GetHttpContext(ginContext)
		authHeader := ginContext.GetHeader("Authorization")
		if authHeader == "" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, dto.HttpResponse[any]{
				HttpStatus: http.StatusUnauthorized,
				Code:       core.Error.Auth.Unauthenticated.Code,
				Message:    core.Translator.Localize(httpContext.Locale(), core.Error.Auth.Unauthenticated.MessageKey),
				Data:       nil,
			})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtManager.Verify(rawToken)
		if err != nil {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, dto.HttpResponse[any]{
				HttpStatus: http.StatusUnauthorized,
				Code:       core.Error.Invalid.Token.Code,
				Message:    core.Translator.Localize("en", core.Error.Invalid.Token.MessageKey),
				Data:       nil,
			})
			return
		}

		// Register SecurityContext in HttpContext
		securityContext := &core.SecurityContext{
			Identity: claims.UserID,
			Roles:    claims.Roles,
		}
		httpContext.SetSecurityContext(securityContext)

		ginContext.Next()
	}
}
