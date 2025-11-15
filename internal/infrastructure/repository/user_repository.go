package repository

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
)

type userRepository struct {
	*Repository[*model.User, string]
}

func NewUserRepository(core *core.Core, postgres *data.PostgresDB) infra_interface.UserRepository {
	return &userRepository{
		Repository: NewRepository[*model.User, string](core, postgres),
	}
}
