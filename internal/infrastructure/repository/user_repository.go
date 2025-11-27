package repository

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
	"veg-store-backend/util"

	"github.com/go-faker/faker/v4"
)

type userRepository struct {
	*Repository[*model.User, model.UUID]
}

func NewUserRepository(core *core.Core) infra_interface.UserRepository {
	return &userRepository{
		Repository: NewRepository[*model.User, model.UUID](core),
	}
}

func (r *userRepository) Seed(db *data.PostgresDB, num int8) error {
	fakeUsers := make([]model.User, 0, num)

	for i := int8(0); i < num; i++ {
		// Random age 18-120
		randomAge, _ := faker.RandomInt(18, 120, 1)

		// Random sex
		randomSex, _ := faker.RandomInt(0, 1, 1)
		sex := randomSex[0] == 0

		// Create fake user
		fakeUser := model.User{
			AuditingModel: model.AuditingModel{
				CreatedAt: util.RandomDateTime(),
			},
			Name:     faker.Name(),
			Age:      int8(randomAge[0]),
			Sex:      sex,
			Email:    faker.Email(),
			Password: util.HashPassword("password123"),
		}
		fakeUser.Created()

		fakeUsers = append(fakeUsers, fakeUser)
	}

	return db.Create(&fakeUsers).Error
}
