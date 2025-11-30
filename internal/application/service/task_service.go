package service

import (
	"context"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/iface"
	"veg-store-backend/internal/application/mapper"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
	"veg-store-backend/util"
)

type TaskService interface {
	Search(ctx context.Context, opt dto.OffsetPageOption) (*dto.OffsetPageResult[dto.TaskItem], error)
	FindAll(ctx context.Context) ([]dto.TaskItem, error)
	Create(ctx context.Context, request dto.CreateTaskRequest) (string, error)
	UpdateStatus(ctx context.Context, request dto.UpdateTaskStatusRequest) (string, error)
	HardDelete(ctx context.Context, id string) (string, error)
}

type taskService struct {
	Service[*data.PostgresDB, iface.TaskRepository]
}

func NewTaskService(
	core *core.Core,
	db *data.PostgresDB,
	repository iface.TaskRepository,
) TaskService {
	return &taskService{
		Service[*data.PostgresDB, iface.TaskRepository]{
			Core:       core,
			DB:         db,
			Repository: repository,
		},
	}
}

func (s *taskService) Search(
	ctx context.Context,
	opt dto.OffsetPageOption,
) (*dto.OffsetPageResult[dto.TaskItem], error) {
	result := &dto.OffsetPageResult[dto.TaskItem]{
		Items: make([]dto.TaskItem, 0),
		Page:  opt.Page,
		Size:  opt.Size,
		Total: 0,
	}

	rawPage, err := s.Repository.OffsetPage(s.DB, ctx, opt)
	if err != nil {
		return result, s.Error.NotFound.Task.MoreInfo(map[string]interface{}{
			"Count": 2,
		})
	}

	result.Items = util.Map(rawPage.Items, func(t *model.Task) dto.TaskItem { return mapper.ToTaskItem(t) })
	result.Total = rawPage.Total
	result.Page = rawPage.Page
	result.Size = rawPage.Size
	return result, nil
}

func (s *taskService) FindAll(ctx context.Context) ([]dto.TaskItem, error) {
	statuses, err := s.Repository.FindAll(s.DB, ctx)
	if err != nil {
		return make([]dto.TaskItem, 0), s.Error.NotFound.Task.MoreInfo(map[string]interface{}{
			"Count": 2,
		})
	}

	return util.Map(statuses, func(t *model.Task) dto.TaskItem { return mapper.ToTaskItem(t) }), nil
}

func (s *taskService) Create(ctx context.Context, request dto.CreateTaskRequest) (string, error) {
	// Map DTO -> ORM Model
	newTask := mapper.ToTask(request)

	// Handle create new task
	err := s.Repository.Create(s.DB, ctx, &newTask)
	if err != nil {
		return "", s.Error.Fail.CreateTask
	}

	return newTask.ID.String(), nil
}

func (s *taskService) UpdateStatus(ctx context.Context, request dto.UpdateTaskStatusRequest) (string, error) {
	// Find existing task then set task's status with new value
	existing, err := s.Repository.FindById(s.DB, ctx, model.ToUUID(request.ID))
	if err != nil {
		return "", s.Error.NotFound.Task
	}
	existing.StatusID = model.ToUUID(request.StatusID)

	// Update existing task if found
	err = s.Repository.Update(s.DB, ctx, existing)
	if err != nil {
		return "", s.Error.Fail.UpdateTask
	}

	return request.ID, nil
}

func (s *taskService) HardDelete(ctx context.Context, id string) (string, error) {
	// Update existing task if found
	err := s.Repository.HardDelete(s.DB, ctx, model.ToUUID(id))
	if err != nil {
		return "", s.Error.Fail.DeleteTask
	}

	return id, nil
}
