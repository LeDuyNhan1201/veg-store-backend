package mockservice

import (
	"veg-store-backend/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mck *MockUserService) FindByUsername(username string) (*model.User, error) {
	args := mck.Called(username)

	// Nếu return value đầu tiên là nil -> return nil
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}

	return user, args.Error(1)
}

func (mck *MockUserService) FindById(id string) (*model.User, error) {
	args := mck.Called(id)

	// Nếu return value đầu tiên là nil -> return nil
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}

	return user, args.Error(1)
}

func (mck *MockUserService) Greeting() string {
	args := mck.Called()
	return args.String(0)
}
