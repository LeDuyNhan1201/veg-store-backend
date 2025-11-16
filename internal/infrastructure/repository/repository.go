package repository

import (
	"context"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
)

type Repository[TEntity infra_interface.IEntity[TId], TId model.AllowedId] struct {
	*core.Core
	Postgres *data.PostgresDB
}

func NewRepository[TEntity infra_interface.IEntity[TId], TId model.AllowedId](core *core.Core, postgres *data.PostgresDB) *Repository[TEntity, TId] {
	return &Repository[TEntity, TId]{
		Core:     core,
		Postgres: postgres,
	}
}

func (r *Repository[TEntity, TId]) Create(ctx context.Context, entity *TEntity) error {
	(*entity).Created()
	return r.Postgres.DB.WithContext(ctx).Create(entity).Error
}

func (r *Repository[TEntity, TId]) FindById(ctx context.Context, id TId) (TEntity, error) {
	var entity TEntity
	err := r.Postgres.DB.WithContext(ctx).First(&entity, id).Error
	return entity, err
}

func (r *Repository[TEntity, TId]) FindAll(ctx context.Context) ([]TEntity, error) {
	var result []TEntity
	err := r.Postgres.DB.WithContext(ctx).Where("is_deleted = false").Find(&result).Error
	return result, err
}

func (r *Repository[TEntity, TId]) Update(ctx context.Context, entity *TEntity) error {
	(*entity).Updated()
	return r.Postgres.DB.WithContext(ctx).Save(entity).Error
}

func (r *Repository[TEntity, TId]) SoftDelete(ctx context.Context, id TId) error {
	entity, err := r.FindById(ctx, id)
	if err != nil {
		return err
	}

	entity.Deleted()
	return r.Postgres.DB.WithContext(ctx).Save(entity).Error
}

func (r *Repository[TEntity, TId]) HardDelete(ctx context.Context, id TId) error {
	return r.Postgres.DB.WithContext(ctx).Delete(new(TEntity), id).Error
}
