package infra_interface

import "veg-store-backend/internal/domain/model"

type IAuditing[TId comparable] interface {
	model.IAuditingModel[TId]
}

type IRepository[TModel IAuditing[TId], TId comparable] interface {
	FindById(id TId) (TModel, bool, error)
	Save(entity TModel) (string, error)
	Delete(id TId) (string, error)
	FindAll() ([]TModel, error)
}

type UserRepository interface {
	IRepository[*model.User, string]
}
