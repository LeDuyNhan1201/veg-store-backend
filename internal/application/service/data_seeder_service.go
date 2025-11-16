package service

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/infrastructure/core"

	"go.uber.org/zap"
)

type DataSeederService interface {
	SeedData()
	SeedUsers(num int8)
}

type dataSeederService struct {
	*core.Core
	userRepository infra_interface.UserRepository
}

func NewDataSeederService(
	core *core.Core,
	userRepository infra_interface.UserRepository,
) DataSeederService {
	return &dataSeederService{
		Core:           core,
		userRepository: userRepository,
	}
}

func (service *dataSeederService) SeedData() {
	if service.isDataSeedingEnabled() {
		service.SeedUsers(10)
	}
}

func (service *dataSeederService) SeedUsers(num int8) {
	err := service.userRepository.Seed(num)
	if err != nil {
		service.Logger.Fatal("Failed to seed users", zap.Error(err))
	}
}

func (service *dataSeederService) isDataSeedingEnabled() bool {
	if !service.Config.Data.EnableDataSeeding {
		service.Logger.Warn("Data seeding is disabled in configuration.")
		return false
	}
	return true
}
