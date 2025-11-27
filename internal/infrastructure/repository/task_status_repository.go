package repository

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
)

type taskStatusRepository struct {
	*Repository[*model.TaskStatus, model.UUID]
}

func NewTaskStatusRepository(core *core.Core) infra_interface.TaskStatusRepository {
	return &taskStatusRepository{
		Repository: NewRepository[*model.TaskStatus, model.UUID](core),
	}
}

func (r *taskStatusRepository) Seed(db *data.PostgresDB) error {
	sampleTaskStatuses := make([]model.TaskStatus, 0, 3)

	ToDoStatus := model.TaskStatus{
		AuditingModel: model.AuditingModel{},
		Title:         "To do",
	}

	InProgressStatus := model.TaskStatus{
		AuditingModel: model.AuditingModel{},
		Title:         "In Progress",
	}

	DoneStatus := model.TaskStatus{
		AuditingModel: model.AuditingModel{},
		Title:         "Done",
	}

	ToDoStatus.Created()
	InProgressStatus.Created()
	DoneStatus.Created()
	sampleTaskStatuses = append(sampleTaskStatuses, ToDoStatus, InProgressStatus, DoneStatus)

	return db.Create(&sampleTaskStatuses).Error
}
