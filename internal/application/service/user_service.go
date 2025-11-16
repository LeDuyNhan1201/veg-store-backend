package service

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"

	"github.com/google/uuid"
)

type UserService interface {
	Greeting() string
	FindById(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type userService struct {
	*core.Core
	repository infra_interface.UserRepository
}

func NewUserService(core *core.Core, repository infra_interface.UserRepository) UserService {
	return &userService{
		Core:       core,
		repository: repository,
	}
}

func (service *userService) Greeting() string {
	return "Hello"
}

func (service *userService) FindById(id string) (*model.User, error) {
	entity, err := service.repository.FindById(nil, model.ToUUID(id))
	if err != nil {
		return nil, service.Error.NotFound.User
	}
	return entity, nil
}

func (service *userService) FindByUsername(username string) (*model.User, error) {
	if username == "test" {
		return nil, service.Error.NotFound.User

	} else {
		return &model.User{
			AuditingModel: model.AuditingModel[model.UUID]{Id: model.ToUUID(uuid.New().String())},
			Name:          "Test",
			Age:           18,
			Sex:           true,
		}, nil
	}
}
