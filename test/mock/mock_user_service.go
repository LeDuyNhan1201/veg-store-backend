package mock

import (
	"veg-store-backend/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (mockService *UserService) FindById(id string) (*model.User, error) {
	args := mockService.Called(id)

	// Nếu return value đầu tiên là nil -> return nil
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}

	return user, args.Error(1)
}

func (mockService *UserService) Greeting() string {
	args := mockService.Called()
	return args.String(0)
}
