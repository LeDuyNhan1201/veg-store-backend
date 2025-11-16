package infra_interface

import (
	"context"
	"veg-store-backend/internal/domain/model"
)

type IEntity[TId model.AllowedId] interface {
	GetId() TId
	Created() *model.AuditingModel[TId]
	Updated() *model.AuditingModel[TId]
	Deleted() *model.AuditingModel[TId]
}

type IRepository[TEntity IEntity[TId], TId model.AllowedId] interface {
	Create(ctx context.Context, entity *TEntity) error
	FindById(ctx context.Context, id TId) (TEntity, error)
	FindAll(ctx context.Context) ([]TEntity, error)
	Update(ctx context.Context, entity *TEntity) error
	SoftDelete(ctx context.Context, id TId) error
	HardDelete(ctx context.Context, id TId) error
}

type UserRepository interface {
	IRepository[*model.User, model.UUID]
	Seed(num int8) error
}
