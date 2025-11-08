package service_test

import (
	"fmt"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mockService *MockUserService) FindByUsername(username string) (*model.User, error) {
	args := mockService.Called(username)

	// Nếu return value đầu tiên là nil -> return nil
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}

	return user, args.Error(1)
}

func (mockService *MockUserService) FindById(id string) (*model.User, error) {
	args := mockService.Called(id)

	// Nếu return value đầu tiên là nil -> return nil
	var user *model.User
	if u := args.Get(0); u != nil {
		user = u.(*model.User)
	}

	return user, args.Error(1)
}

func (mockService *MockUserService) Greeting() string {
	args := mockService.Called()
	return args.String(0)
}

/*----------------------------------INJECTION--------------------------------------*/

func (mockService *MockUserService) Name() string { return "MockUserService" }
func (mockService *MockUserService) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", mockService.Name()))
	return nil
}
func (mockService *MockUserService) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s destroyed", mockService.Name()))
	return nil
}
