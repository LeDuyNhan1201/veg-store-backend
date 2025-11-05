package rest

import (
	"net/http"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

func (handler *UserHandler) Hello(context *core.HttpContext) {
	context.JSON(http.StatusOK, gin.H{
		"message": context.T(handler.service.Greeting(), map[string]interface{}{
			"name": "Ben",
		}),
	})
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate user and return a token
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  dto.SignInRequest  true  "User credentials"
// @Success 200 {object} dto.Tokens
// @Failure 401 {object} string
// @Router /user/sign-in [post]
func (handler *UserHandler) SignIn(context *core.HttpContext) {
	var request dto.SignInRequest

	if err := context.Gin.ShouldBindJSON(&request); err != nil {
		context.Gin.Error(core.Error.Auth.Unauthenticated)
		return
	}

	context.JSON(http.StatusOK, dto.HttpResponse[dto.Tokens]{
		HttpStatus: http.StatusOK,
		Data: dto.Tokens{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
		},
	})
}

func (handler *UserHandler) Details(context *core.HttpContext) {
	id := context.Gin.Param("id")
	user, err := handler.service.FindById(id)
	if err != nil {
		context.Gin.Error(err)
	} else {
		context.JSON(http.StatusOK, dto.HttpResponse[*model.User]{
			HttpStatus: http.StatusOK,
			Data:       user,
		})
	}
}

func (handler *UserHandler) HealthCheck(ctx *core.HttpContext) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (handler *UserHandler) GetAllUsers(ctx *core.HttpContext) {
	ctx.JSON(http.StatusOK, gin.H{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
}
