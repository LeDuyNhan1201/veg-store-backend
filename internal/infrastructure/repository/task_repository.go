package repository

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
	"veg-store-backend/util"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

type taskRepository struct {
	*Repository[*model.Task, model.UUID]
}

func NewTaskRepository(core *core.Core) infra_interface.TaskRepository {
	return &taskRepository{
		Repository: NewRepository[*model.Task, model.UUID](core),
	}
}

func (r *taskRepository) Seed(db *data.PostgresDB, num int8, statusIds []model.UUID) error {
	fakeTasks := make([]model.Task, 0, num)

	for i := int8(0); i < num; i++ {
		randomStatusIndex, err := faker.RandomInt(0, len(statusIds)-1, 1)
		startDay := util.RandomDateTime()
		targetDay := util.RandomTimeBetween(startDay, startDay.AddDate(0, 0, 7))
		endDay := util.RandomTimeBetween(startDay.AddDate(0, 0, 1), targetDay)

		if err != nil {
			return err
		}

		fakeTask := model.Task{
			AuditingModel: model.AuditingModel{
				CreatedAt: util.RandomDateTime(),
			},
			Title:     faker.Sentence(options.WithRandomStringLength(100)),
			StartDay:  &startDay,
			TargetDay: &targetDay,
			EndDay:    &endDay,
			StatusID:  statusIds[randomStatusIndex[0]],
		}

		fakeTask.Created()
		fakeTasks = append(fakeTasks, fakeTask)
	}
	return db.Create(&fakeTasks).Error
}
