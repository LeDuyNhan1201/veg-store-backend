package middleware

import (
	"fmt"
	"strings"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

type LocaleMiddleware struct {
	*Middleware
}

func NewLocaleMiddleware(core *core.Core, router *router.HTTPRouter) *LocaleMiddleware {
	return &LocaleMiddleware{
		Middleware: &Middleware{
			Core:   core,
			Router: router,
		},
	}
}

func (m *LocaleMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		lang := ginContext.GetHeader("Accept-Language")

		m.Logger.Info(fmt.Sprintf("Accept-Language: %s", lang))

		if lang == "" {
			lang = m.AppConfig.Server.DefaultLocale

		} else {
			// Just get the first language tag, e.g. "en-US,en;q=0.9" -> "en"
			lang = strings.Split(lang, ",")[0]
			lang = strings.Split(lang, "-")[0]
			lang = strings.TrimSpace(lang)
		}

		// Save locale to context
		ginContext.Set(util.LocaleContextKey, lang)

		m.Logger.Debug("[BEFORE] Locale invoked")
		ginContext.Next()
		m.Logger.Debug("[AFTER] Locale invoked")
	}
}

func (m *LocaleMiddleware) Setup() {
	m.Router.Engine.Use(m.handler())
}

func (m *LocaleMiddleware) Priority() uint8 {
	return util.LocaleMiddlewarePriority
}
