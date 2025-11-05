package injection_test

import (
	"veg-store-backend/injection/core"

	"github.com/gin-gonic/gin"
)

func MockHttpContext(
	ginCtx *gin.Context,
) *core.HttpContext {
	return &core.HttpContext{
		Translator: core.Translator,
		Gin:        ginCtx,
	}
}
