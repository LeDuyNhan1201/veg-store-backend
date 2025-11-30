package resthandler

import (
	"net/http"
	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

// Hello godoc
// @Summary Anh trai say hi
// @Description Anh trai say gex
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.HttpResponse[string]
// @Router /users/hello [get]
func (h *UserHandler) Hello(ctx *context.Http) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": ctx.T(h.service.Greeting(), map[string]interface{}{
			"DBName": "Ben",
			"Count":  1,
		}),
	})
}

// Details godoc
// @Security BearerAuth
// @Summary User details
// @Description Get details of a user by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} dto.HttpResponse[model.User]
// @Failure 400 {object} dto.HttpResponse[any]
// @Router /users/{id} [get]
func (h *UserHandler) Details(ctx *context.Http) {
	id := ctx.Gin.Param("id")
	user, err := h.service.FindById(id)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.HttpResponse[*model.User]{
		HttpStatus: http.StatusOK,
		Data:       user,
	})
}

func (h *UserHandler) GetAllUsers(ctx *context.Http) {
	ctx.JSON(http.StatusOK, gin.H{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
}
