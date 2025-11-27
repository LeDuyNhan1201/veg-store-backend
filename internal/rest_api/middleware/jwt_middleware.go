package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JWTMiddleware struct {
	*Middleware
	jwtManager infra_interface.JWTManager
}

func NewJWTMiddleware(core *core.Core, router *router.HTTPRouter) *JWTMiddleware {
	return &JWTMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (m *JWTMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		path := ginContext.FullPath()
		// If path empty (don't match any routes) fallback to RequestURI
		if path == "" {
			path = ginContext.Request.URL.Path
		}
		m.Logger.Info("Verifying JWT", zap.String("path", path))

		// Bypass JWT verification for public endpoints
		for _, endpoint := range m.AppConfig.Security.PublicEndpoints {
			if matchPath(fmt.Sprintf("%s%s%s",
				m.AppConfig.Server.ApiPrefix,
				m.AppConfig.Server.ApiVersion, endpoint,
			), path) {
				m.Logger.Debug("[JWT] Skipped public endpoint", zap.String("path", path))
				ginContext.Next()
				return
			}
		}

		httpContext := context.GetHttpContext(ginContext)
		authHeader := ginContext.GetHeader("Authorization")
		if authHeader == "" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, dto.HttpResponse[any]{
				HttpStatus: http.StatusUnauthorized,
				Code:       m.Error.Auth.Unauthenticated.Code,
				Message:    m.Localizer.Localize(httpContext.Locale(), m.Error.Auth.Unauthenticated.MessageKey),
				Data:       nil,
			})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwtManager.Verify(rawToken)
		if err != nil {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, dto.HttpResponse[any]{
				HttpStatus: http.StatusUnauthorized,
				Code:       m.Error.Invalid.Token.Code,
				Message: m.Localizer.Localize(
					m.AppConfig.Server.DefaultLocale,
					m.Error.Invalid.Token.MessageKey,
				),
				Data: nil,
			})
			return
		}

		// Register SecurityContext in Http
		securityContext := &context.SecurityContext{
			Identity: claims.UserId,
			Roles:    claims.Roles,
		}
		httpContext.SetSecurityContext(securityContext)

		m.Logger.Debug("[BEFORE] JWT invoked")
		ginContext.Next()
		m.Logger.Debug("[AFTER] JWT invoked")
	}
}

func matchPath(pattern, path string) bool {
	if strings.HasSuffix(pattern, "/*any") {
		prefix := strings.TrimSuffix(pattern, "/*any")
		return strings.HasPrefix(path, prefix)
	}
	return pattern == path
}

func (m *JWTMiddleware) Setup() {
	m.Router.Engine.Use(m.handler())
}

func (m *JWTMiddleware) Priority() uint8 {
	return util.JWTMiddlewarePriority
}
