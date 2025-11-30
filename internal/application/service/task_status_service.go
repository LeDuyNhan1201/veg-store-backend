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

type TaskStatusService interface {
	FindAll(ctx context.Context) ([]dto.PreviewStatus, error)
}

type taskStatusService struct {
	Service[*data.PostgresDB, iface.TaskStatusRepository]
}

func NewTaskStatusService(
	core *core.Core,
	db *data.PostgresDB,
	repository iface.TaskStatusRepository,
) TaskStatusService {
	return &taskStatusService{
		Service[*data.PostgresDB, iface.TaskStatusRepository]{
			Core:       core,
			DB:         db,
			Repository: repository,
		},
	}
}

func (s *taskStatusService) FindAll(ctx context.Context) ([]dto.PreviewStatus, error) {
	statuses, err := s.Repository.FindAll(s.DB, ctx)
	if err != nil {
		return make([]dto.PreviewStatus, 0), s.Error.NotFound.Task.MoreInfo(map[string]interface{}{
			"Count": 2,
		})
	}

	return util.Map(statuses, func(s *model.TaskStatus) dto.PreviewStatus { return mapper.ToPreviewTaskStatus(s) }), nil
}
