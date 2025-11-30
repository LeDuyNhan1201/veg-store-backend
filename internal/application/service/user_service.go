package service

import (
	"veg-store-backend/internal/application/iface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
)

type UserService interface {
	Greeting() string
	FindById(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type userService struct {
	Service[*data.PostgresDB, iface.UserRepository]
}

func NewUserService(
	core *core.Core,
	db *data.PostgresDB,
	repository iface.UserRepository,
) UserService {
	return &userService{
		Service[*data.PostgresDB, iface.UserRepository]{
			Core:       core,
			DB:         db,
			Repository: repository,
		},
	}
}

func (s *userService) Greeting() string {
	return "Hello"
}

func (s *userService) FindById(id string) (*model.User, error) {
	entity, err := s.Repository.FindById(s.DB, nil, model.ToUUID(id))
	if err != nil {
		return nil, s.Error.NotFound.User
	}
	return entity, nil
}

func (s *userService) FindByUsername(username string) (*model.User, error) {
	if username == "test" {
		return nil, s.Error.NotFound.User

	} else {
		return &model.User{
			ID:            model.NewUUID(),
			AuditingModel: model.AuditingModel{},
			Name:          "Test",
			Age:           18,
			Sex:           true,
		}, nil
	}
}
