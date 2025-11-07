package repository

import (
	"context"
	"fmt"
	"veg-store-backend/injection/core"

	"go.uber.org/fx"
)

type UserRepository interface {
	Name() string
	Start() error
	Stop() error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

/*----------------------------------INJECTION--------------------------------------*/

func (repository *userRepository) Name() string { return "UserRepository" }
func (repository *userRepository) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", repository.Name()))
	return nil
}
func (repository *userRepository) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s destroyed", repository.Name()))
	return nil
}

func RegisterUserRepository(lifecycle fx.Lifecycle, repository UserRepository) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return repository.Start()
		},
		OnStop: func(context context.Context) error {
			return repository.Stop()
		},
	})
}

var UserRepositoryModule = fx.Options(fx.Provide(NewUserRepository), fx.Invoke(RegisterUserRepository))
