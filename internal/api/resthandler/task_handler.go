package resthandler

import (
	"net/http"
	"strconv"

	"veg-store-backend/internal/application/context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/mapper"
	"veg-store-backend/internal/application/service"
)

type TaskHandler struct {
	taskService       service.TaskService
	taskStatusService service.TaskStatusService
}

func NewTaskHandler(
	taskService service.TaskService,
	taskStatusService service.TaskStatusService,
) *TaskHandler {
	return &TaskHandler{
		taskService:       taskService,
		taskStatusService: taskStatusService,
	}
}

// Search godoc
// @Summary Search
// @Description Search tasks and return offset page
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body dto.AdvancedFilterTaskRequest true "Filters"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} dto.HttpResponse[dto.OffsetPageResult[dto.TaskItem]]
// @Failure 401 {object} dto.HttpResponse[string]
// @Router /tasks/search [post]
func (h *TaskHandler) Search(ctx *context.Http) {
	pageNum, _ := strconv.Atoi(ctx.Gin.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Gin.DefaultQuery("size", "10"))
	var request dto.AdvancedFilterTaskRequest
	err := ctx.Gin.ShouldBindJSON(&request)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}

	pageOptions := mapper.ToOffsetPageOption(request, pageNum, pageSize)
	page, err := h.taskService.Search(ctx.Gin, pageOptions)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.HttpResponse[dto.OffsetPageResult[dto.TaskItem]]{
		HttpStatus: http.StatusOK,
		Data:       *page,
	})
}

// FindAllStatuses godoc
// @Summary Get All Statuses
// @Description Get All Preview Statuses
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} dto.HttpResponse[[]dto.PreviewStatus]
// @Failure 400 {object} dto.HttpResponse[any]
// @Router /tasks/statuses [get]
func (h *TaskHandler) FindAllStatuses(ctx *context.Http) {
	previewStatuses, err := h.taskStatusService.FindAll(ctx.Gin)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.HttpResponse[[]dto.PreviewStatus]{
		HttpStatus: http.StatusOK,
		Data:       previewStatuses,
	})
}

// FindAll godoc
// @Summary Get All Tasks
// @Description Get All Tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} dto.HttpResponse[[]dto.TaskItem]
// @Failure 400 {object} dto.HttpResponse[any]
// @Router /tasks [get]
func (h *TaskHandler) FindAll(ctx *context.Http) {
	taskItems, err := h.taskService.FindAll(ctx.Gin)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.HttpResponse[[]dto.TaskItem]{
		HttpStatus: http.StatusOK,
		Data:       taskItems,
	})
}

// Create godoc
// @Summary Create a task
// @Description Create new task and return task id
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task info"
// @Success 201 {object} dto.HttpResponse[string]
// @Failure 401 {object} dto.HttpResponse[string]
// @Router /tasks [post]
func (h *TaskHandler) Create(ctx *context.Http) {
	var request dto.CreateTaskRequest
	err := ctx.Gin.ShouldBindJSON(&request)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}

	taskID, err := h.taskService.Create(ctx.Gin, request)
	if err != nil {
		ctx.Gin.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.HttpResponse[string]{
		HttpStatus: http.StatusCreated,
		Data:       taskID,
	})
}
