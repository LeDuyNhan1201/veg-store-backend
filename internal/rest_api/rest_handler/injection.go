package rest_handler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserHandler),
	fx.Provide(NewAuthHandler),
	fx.Provide(NewTaskHandler),
)
