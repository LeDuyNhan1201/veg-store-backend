package rest_test

import (
	"fmt"
	"net/http"
	"testing"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/rest_api/rest_handler"
	"veg-store-backend/test/service_mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type UserHandler struct {
	*HandlerTest[*rest_handler.UserHandler, *service_mock.MockUserService]
}

func setupUserHandlerTest() *UserHandler {
	mockService := new(service_mock.MockUserService)
	mockHandler := rest_handler.NewUserHandler(mockService)
	handlerTest := NewHandlerTest[*rest_handler.UserHandler, *service_mock.MockUserService](mockHandler, mockService)
	handlerTest.MockUserRoute(mockHandler)
	return &UserHandler{
		HandlerTest: handlerTest,
	}
}

// FunctionName_Condition1_Condition2_ExpectedResult
func (testHandler *UserHandler) hello_success(test *testing.T) {
	// GIVEN
	expectedGreeting := "Hello Ben"

	// WHEN
	testHandler.MockService.On("Greeting").Return(expectedGreeting)
	httpRecorder := testHandler.Get(test, testHandler.AppURI("/user/hello"))

	// THEN & ASSERT
	assert.Equal(test, http.StatusOK, httpRecorder.Code)
	assert.Contains(test, httpRecorder.Body.String(), "Hello Ben")
	testHandler.MockService.AssertExpectations(test)
}

func (testHandler *UserHandler) details_withNotFoundID_fail(test *testing.T) {
	// GIVEN
	notFoundIdSample := "1b332625-8949-4e3b-a10e-b3291101a341"
	expectedError := testHandler.Error.NotFound.User

	// WHEN
	testHandler.MockService.On("FindById", notFoundIdSample).Return(nil, expectedError)
	responseRecorder := testHandler.Get(test, testHandler.AppURI(fmt.Sprintf("/user/%s", notFoundIdSample)))
	testHandler.Logger.Debug("Response Body:", zap.Any("body", responseRecorder))
	// THEN
	var response dto.HttpResponse[any]
	testHandler.DecodeResponse(test, responseRecorder, &response) // map responseRecorder -> response struct

	// ASSERT
	assert.Equal(test, http.StatusNotFound, response.HttpStatus)
	assert.Equal(test, nil, response.Data)
	testHandler.MockService.AssertExpectations(test)
}

func TestUserHandler(test *testing.T) {
	// SETUP
	mockHandler := setupUserHandlerTest()

	// RUN TESTS
	test.Run("hello_success", mockHandler.hello_success)
	test.Run("details_withNotFoundID_fail", mockHandler.details_withNotFoundID_fail)
}
