package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWT(jwtManager infra_interface.JWTManager) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		path := ginContext.FullPath()
		// If path empty (don't match any routes) fallback to RequestURI
		if path == "" {
			path = ginContext.Request.URL.Path
		}
		core.Logger.Info("Verifying JWT", zap.String("path", path))

		// Bypass JWT verification for public endpoints
		for _, endpoint := range core.Configs.Security.PublicEndpoints {
			if matchPath(fmt.Sprintf("%s%s%s",
				core.Configs.Server.ApiPrefix,
				core.Configs.Server.ApiVersion, endpoint,
			), path) {
				core.Logger.Debug("[JWT] Skipped public endpoint", zap.String("path", path))
				ginContext.Next()
				return
			}
		}

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

		core.Logger.Debug("[BEFORE] JWT invoked")
		ginContext.Next()
		core.Logger.Debug("[AFTER] JWT invoked")
	}
}

func matchPath(pattern, path string) bool {
	if strings.HasSuffix(pattern, "/*any") {
		prefix := strings.TrimSuffix(pattern, "/*any")
		return strings.HasPrefix(path, prefix)
	}
	return pattern == path
}
