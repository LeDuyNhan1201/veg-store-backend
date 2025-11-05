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
			"name": "Phat",
		}),
	})
}

func (handler *UserHandler) Details(context *core.HttpContext) {
	id := context.GinContext.Param("id")
	user, err := handler.service.FindById(id)
	if err != nil {
		context.GinContext.Error(err)
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
