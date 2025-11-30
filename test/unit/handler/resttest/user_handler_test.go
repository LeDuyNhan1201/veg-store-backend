package resttest

import (
	"fmt"
	"net/http"
	"testing"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/api/resthandler"
	"veg-store-backend/test/mockservice"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type UserHandler struct {
	*HandlerTest[*resthandler.UserHandler, *mockservice.MockUserService]
}

func setupUserHandlerTest() *UserHandler {
	mockService := new(mockservice.MockUserService)
	mockHandler := resthandler.NewUserHandler(mockService)
	handlerTest := NewHandlerTest[*resthandler.UserHandler, *mockservice.MockUserService](mockHandler, mockService)
	handlerTest.MockUserRoute(mockHandler)
	return &UserHandler{
		HandlerTest: handlerTest,
	}
}

// FunctionName_Condition1_Condition2_ExpectedResult
func (h *UserHandler) hello_success(test *testing.T) {
	// GIVEN
	expectedGreeting := "Hello Ben"

	// WHEN
	h.MockService.On("Greeting").Return(expectedGreeting)
	httpRecorder := h.Get(test, h.AppURI("/users/hello"))

	// THEN & ASSERT
	assert.Equal(test, http.StatusOK, httpRecorder.Code)
	assert.Contains(test, httpRecorder.Body.String(), "Hello Ben")
	h.MockService.AssertExpectations(test)
}

func (h *UserHandler) details_withNotFoundID_fail(test *testing.T) {
	// GIVEN
	notFoundIdSample := "1b332625-8949-4e3b-a10e-b3291101a341"
	expectedError := h.Error.NotFound.User

	// WHEN
	h.MockService.On("FindById", notFoundIdSample).Return(nil, expectedError)
	responseRecorder := h.Get(test, h.AppURI(fmt.Sprintf("/users/%s", notFoundIdSample)))
	h.Logger.Debug("Response Body:", zap.Any("body", responseRecorder))
	// THEN
	var response dto.HttpResponse[any]
	h.DecodeResponse(test, responseRecorder, &response) // map responseRecorder -> response struct

	// ASSERT
	assert.Equal(test, http.StatusNotFound, response.HttpStatus)
	assert.Equal(test, nil, response.Data)
	h.MockService.AssertExpectations(test)
}

func TestUserHandler(test *testing.T) {
	// SETUP
	mockHandler := setupUserHandlerTest()

	// RUN TESTS
	test.Run("hello_success", mockHandler.hello_success)
	test.Run("details_withNotFoundID_fail", mockHandler.details_withNotFoundID_fail)
}
