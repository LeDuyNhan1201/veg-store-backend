package rest_handler

import (
	"net/http"
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/service"
)

type AuthHandler struct {
	service service.AuthenticationService
}

func NewAuthHandler(authenticationService service.AuthenticationService) *AuthHandler {
	return &AuthHandler{
		service: authenticationService,
	}
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate user and return a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.SignInRequest true "User credentials"
// @Success 200 {object} dto.HttpResponse[dto.Tokens]
// @Failure 401 {object} dto.HttpResponse[string]
// @Router /auth [post]
func (h *AuthHandler) SignIn(context *context.Http) {
	var request dto.SignInRequest
	err := context.Gin.ShouldBindJSON(&request)
	if err != nil {
		context.Gin.Error(err)
		return
	}

	tokens, err := h.service.Tokens(request)
	if err != nil {
		context.Gin.Error(context.Error.Auth.Unauthenticated)
		return
	}

	context.JSON(http.StatusOK, dto.HttpResponse[dto.Tokens]{
		HttpStatus: http.StatusOK,
		Data:       *tokens,
	})
}

// Info godoc
// @Security BearerAuth
// @Summary Me
// @Description Get name of a current user
// @Tags Auth
// @Produce json
// @Success 200 {object} dto.HttpResponse[string]
// @Failure 401 {object} dto.HttpResponse[any]
// @Router /auth/me [get]
func (h *AuthHandler) Info(context *context.Http) {
	id := context.SecurityContext.Identity
	username, err := h.service.Me(id)
	if err != nil {
		context.Gin.Error(err)
		return
	}
	context.JSON(http.StatusOK, dto.HttpResponse[string]{
		HttpStatus: http.StatusOK,
		Data:       username,
	})
}
