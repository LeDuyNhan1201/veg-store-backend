package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewAuthenticationService),
	fx.Provide(NewTaskStatusService),
	fx.Provide(NewTaskService),
	fx.Provide(NewDataSeederService),
)
