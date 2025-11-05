package service

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/data"
)

type UserService interface {
	Greeting() string
	FindById(id string) (*model.User, error)
}

type userService struct {
	repo data.UserRepository
}

func NewUserService(repo data.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) Greeting() string {
	return "hello"
}

func (service *userService) FindById(id string) (*model.User, error) {
	if id == "1" {
		return nil, core.Error.NotFound.User

	} else {
		return &model.User{
			Name: "Test",
			Age:  18,
			Sex:  true,
		}, nil
	}
}
