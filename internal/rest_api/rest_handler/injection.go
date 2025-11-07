package rest_handler

import "go.uber.org/fx"

var Module = fx.Options(
	UserHandlerModule,
	AuthHandlerModule,
)
