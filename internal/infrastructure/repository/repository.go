package repository

import (
	"errors"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"

	"gorm.io/gorm"
)

type Repository[TModel infra_interface.IAuditing[TId], TId comparable] struct {
	*core.Core
	postgres *data.PostgresDB
}

func NewRepository[TModel infra_interface.IAuditing[TId], TId comparable](core *core.Core, postgres *data.PostgresDB) *Repository[TModel, TId] {
	return &Repository[TModel, TId]{
		Core:     core,
		postgres: postgres,
	}
}

func (repository *Repository[TModel, TId]) FindById(id TId) (TModel, bool, error) {
	var entity TModel
	result := repository.postgres.DB.First(&entity, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var zero TModel
		return zero, false, nil
	}

	return entity, true, result.Error
}

func (repository *Repository[TModel, TId]) Save(model TModel) (TId, error) {
	var zeroId TId
	if model.GetId() == zeroId {
		model.Created()
		return repository.create(model)
	}

	_, found, err := repository.FindById(model.GetId())
	if err != nil {
		return model.GetId(), err
	}

	if !found {
		model.Created()
		return repository.create(model)
	}

	model.Updated()
	return repository.update(model)
}

func (repository *Repository[TModel, TId]) create(model TModel) (TId, error) {
	model.Created()
	if err := repository.postgres.DB.Create(model).Error; err != nil {
		return model.GetId(), err
	}
	return model.GetId(), nil
}

func (repository *Repository[TModel, TId]) update(model TModel) (TId, error) {
	model.Updated()
	if err := repository.postgres.DB.Save(model).Error; err != nil {
		return model.GetId(), err
	}
	return model.GetId(), nil
}

func (repository *Repository[TModel, TId]) Delete(id TId) (TId, error) {
	existing, found, err := repository.FindById(id)
	if err != nil {
		return id, err
	}

	if !found {
		return id, nil // Not found or already deleted
	}

	existing.Deleted()
	if err := repository.postgres.DB.Save(existing).Error; err != nil {
		return id, err
	}

	return existing.GetId(), nil
}

func (repository *Repository[TModel, TId]) FindAll() ([]TModel, error) {
	var list []TModel
	if err := repository.postgres.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
