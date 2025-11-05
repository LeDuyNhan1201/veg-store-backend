package rest_test

import (
	"net/http"
	"testing"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/api/rest"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/test/mock"
	"veg-store-backend/test/unit/injection_test"

	"github.com/stretchr/testify/assert"
)

type UserHandler struct {
	*HandlerTest[*rest.UserHandler, *mock.UserService]
}

func setupUserHandlerTest() *UserHandler {
	mockService := new(mock.UserService)
	handler := rest.NewUserHandler(mockService)
	engine := injection_test.MockUserRoutes(handler)

	handlerTest := NewHandlerTest[*rest.UserHandler, *mock.UserService](engine, handler, mockService)
	return &UserHandler{
		HandlerTest: handlerTest,
	}
}

func (testHandler *UserHandler) TestHello_success(test *testing.T) {
	testHandler.MockService.On("Greeting").Return("Hello Ben")
	httpRecorder := testHandler.Get(test, "/api/v1/user/hello")
	assert.Equal(test, http.StatusOK, httpRecorder.Code)
	assert.Contains(test, httpRecorder.Body.String(), "Hello Ben")
	testHandler.MockService.AssertExpectations(test)
}

func (testHandler *UserHandler) TestDetails_withNotFoundID_fail(test *testing.T) {
	testHandler.MockService.On("FindById", "123").Return(nil, core.Error.NotFound.User)
	responseRecorder := testHandler.Get(test, "/api/v1/user/details/123")

	var response dto.HttpResponse[any]
	testHandler.DecodeResponse(test, responseRecorder, &response)
	assert.Equal(test, http.StatusNotFound, response.HttpStatus)
	assert.Equal(test, nil, response.Data)
	testHandler.MockService.AssertExpectations(test)
}

func TestUserHandler(test *testing.T) {
	handler := setupUserHandlerTest()
	test.Run("TestHello_success", handler.TestHello_success)
	test.Run("TestDetails_withNotFoundID_fail", handler.TestDetails_withNotFoundID_fail)
}
