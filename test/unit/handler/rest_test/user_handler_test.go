package rest_test

import (
	"net/http"
	"testing"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/rest_api/rest_handler"
	"veg-store-backend/test/service_test"
	"veg-store-backend/test/unit/injection_test"

	"github.com/stretchr/testify/assert"
)

type UserHandler struct {
	*HandlerTest[*rest_handler.UserHandler, *service_test.MockUserService]
}

func setupUserHandlerTest() *UserHandler {
	mockService := new(service_test.MockUserService)
	mockHandler := rest_handler.NewUserHandler(mockService)
	engine := injection_test.MockUserRoutes(mockHandler)

	handlerTest := NewHandlerTest[*rest_handler.UserHandler, *service_test.MockUserService](engine, mockHandler, mockService)
	return &UserHandler{
		HandlerTest: handlerTest,
	}
}

func (testHandler *UserHandler) TestHello_success(test *testing.T) {
	// GIVEN
	expectedGreeting := "Hello Ben"

	// WHEN
	testHandler.MockService.On("Greeting").Return(expectedGreeting)
	httpRecorder := testHandler.Get(test, AppURI("/user/hello"))

	// THEN & ASSERT
	assert.Equal(test, http.StatusOK, httpRecorder.Code)
	assert.Contains(test, httpRecorder.Body.String(), "Hello Ben")
	testHandler.MockService.AssertExpectations(test)
}

func (testHandler *UserHandler) TestDetails_withNotFoundID_fail(test *testing.T) {
	// GIVEN
	notFoundIdSample := "123"
	expectedError := core.Error.NotFound.User

	// WHEN
	testHandler.MockService.On("FindById", notFoundIdSample).Return(nil, expectedError)
	responseRecorder := testHandler.Get(test, AppURI("/user/details/123"))

	// THEN
	var response dto.HttpResponse[any]
	testHandler.DecodeResponse(test, responseRecorder, &response)

	// ASSERT
	assert.Equal(test, http.StatusNotFound, response.HttpStatus)
	assert.Equal(test, nil, response.Data)
	testHandler.MockService.AssertExpectations(test)
}

func TestUserHandler(test *testing.T) {
	// SETUP
	injection_test.MockGlobalComponents()
	mockHandler := setupUserHandlerTest()

	// RUN TESTS
	test.Run("TestHello_success", mockHandler.TestHello_success)
	test.Run("TestDetails_withNotFoundID_fail", mockHandler.TestDetails_withNotFoundID_fail)
}
