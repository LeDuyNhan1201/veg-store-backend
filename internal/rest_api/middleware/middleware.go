package middleware

import (
	"sort"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/router"

	"go.uber.org/fx"
)

type Middleware struct {
	*core.Core
	Router *router.HTTPRouter
}

type IMiddleware interface {
	Priority() uint8
	Setup()
}

type MiddlewaresCollection []IMiddleware

func NewMiddlewaresCollection(
	localMiddleware *LocaleMiddleware,
	httpMiddleware *HTTPMiddleware,
	jwtMiddleware *JWTMiddleware,
	traceIDMiddleware *TraceIDMiddleware,
	validationMiddleware *ValidationMiddleware,
	errorHandlingMiddleware *ErrorHandlingMiddleware,
) *MiddlewaresCollection {
	middlewares := MiddlewaresCollection{
		localMiddleware,
		httpMiddleware,
		jwtMiddleware,
		traceIDMiddleware,
		validationMiddleware,
		errorHandlingMiddleware,
	}

	// Sort by increment priority
	sort.Slice(middlewares, func(i, j int) bool {
		return middlewares[i].Priority() < middlewares[j].Priority()
	})

	return &middlewares
}

func (c *MiddlewaresCollection) Setup() {
	for _, middleware := range *c {
		middleware.Setup()
	}
}

// Module IMPORTANT: REMEMBER TO ADD NEW MIDDLEWARE TO Module
var Module = fx.Options(
	fx.Provide(NewLocaleMiddleware),
	fx.Provide(NewHTTPMiddleware),
	fx.Provide(NewJWTMiddleware),
	fx.Provide(NewTraceIDMiddleware),
	fx.Provide(NewValidationMiddleware),
	fx.Provide(NewErrorHandlingMiddleware),
	fx.Provide(NewMiddlewaresCollection),
)
