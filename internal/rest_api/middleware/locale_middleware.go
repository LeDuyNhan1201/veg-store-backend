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

func (middleware *LocaleMiddleware) handler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		lang := ginContext.GetHeader("Accept-Language")

		middleware.Logger.Info(fmt.Sprintf("Accept-Language: %s", lang))

		if lang == "" {
			lang = middleware.Config.Server.DefaultLocale

		} else {
			// Just get the first language tag, e.g. "en-US,en;q=0.9" -> "en"
			lang = strings.Split(lang, ",")[0]
			lang = strings.Split(lang, "-")[0]
			lang = strings.TrimSpace(lang)
		}

		// Save locale to context
		ginContext.Set(util.LocaleContextKey, lang)

		middleware.Logger.Debug("[BEFORE] Locale invoked")
		ginContext.Next()
		middleware.Logger.Debug("[AFTER] Locale invoked")
	}
}

func (middleware *LocaleMiddleware) Setup() {
	middleware.Router.Engine.Use(middleware.handler())
}

func (middleware *LocaleMiddleware) Priority() uint8 {
	return util.LocaleMiddlewarePriority
}
