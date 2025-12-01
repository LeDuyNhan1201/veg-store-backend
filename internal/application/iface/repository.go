package iface

import (
	"context"

	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/data"
)

type IEntity interface {
	Created() *model.AuditingModel
	Updated() *model.AuditingModel
	Deleted() *model.AuditingModel
}

type IRepository[TEntity IEntity, TId model.AllowedId] interface {
	Create(db *data.PostgresDB, ctx context.Context, entity TEntity) error
	FindById(db *data.PostgresDB, ctx context.Context, id TId, opt ...dto.FindByIDOption) (TEntity, error)
	FindAll(db *data.PostgresDB, ctx context.Context) ([]TEntity, error)
	Update(db *data.PostgresDB, ctx context.Context, entity TEntity) error
	SoftDelete(db *data.PostgresDB, ctx context.Context, id TId) error
	HardDelete(db *data.PostgresDB, ctx context.Context, id TId) error
	OffsetPage(
		db *data.PostgresDB, ctx context.Context,
		opt dto.OffsetPageOption,
	) (dto.OffsetPageResult[TEntity], error)
}

type UserRepository interface {
	IRepository[*model.User, model.UUID]
	Seed(db *data.PostgresDB, num int8) error
}

type TaskRepository interface {
	IRepository[*model.Task, model.UUID]
	Seed(db *data.PostgresDB, num int8, statusIds []model.UUID) error
}

type TaskStatusRepository interface {
	IRepository[*model.TaskStatus, model.UUID]
	Seed(db *data.PostgresDB) error
}
