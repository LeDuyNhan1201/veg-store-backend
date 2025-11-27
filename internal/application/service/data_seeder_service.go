package service

import (
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/internal/infrastructure/data"
	"veg-store-backend/util"

	"go.uber.org/zap"
)

type DataSeederService interface {
	SeedData()
	SeedUsers(num int8)
}

type dataSeederService struct {
	Service[*data.PostgresDB, infra_interface.UserRepository]
	taskStatusRepository infra_interface.TaskStatusRepository
	taskRepository       infra_interface.TaskRepository
}

func NewDataSeederService(
	core *core.Core,
	db *data.PostgresDB,
	userRepository infra_interface.UserRepository,
	taskStatusRepository infra_interface.TaskStatusRepository,
	taskRepository infra_interface.TaskRepository,
) DataSeederService {
	return &dataSeederService{
		Service: Service[*data.PostgresDB, infra_interface.UserRepository]{
			Core:       core,
			DB:         db,
			Repository: userRepository,
		},
		taskStatusRepository: taskStatusRepository,
		taskRepository:       taskRepository,
	}
}

func (s *dataSeederService) SeedData() {
	if s.isDataSeedingEnabled() {
		s.SeedUsers(10)

		s.SeedTaskStatuses()
		statuses, err := s.taskStatusRepository.FindAll(s.DB, nil)
		if err != nil {
			s.Logger.Fatal("Failed to get statuses", zap.Error(err))
		}
		s.Logger.Info("Successfully seeded all statuses", zap.Int("num", len(statuses)))
		s.SeedTasks(20, util.Map(statuses, func(s *model.TaskStatus) model.UUID { return s.ID }))

	}
}

func (s *dataSeederService) SeedUsers(num int8) {
	err := s.Repository.Seed(s.DB, num)
	if err != nil {
		s.Logger.Fatal("Failed to seed users", zap.Error(err))
	}
}

func (s *dataSeederService) SeedTaskStatuses() {
	err := s.taskStatusRepository.Seed(s.DB)
	if err != nil {
		s.Logger.Fatal("Failed to seed users", zap.Error(err))
	}
}

func (s *dataSeederService) SeedTasks(num int8, statusIds []model.UUID) {
	err := s.taskRepository.Seed(s.DB, num, statusIds)
	if err != nil {
		s.Logger.Fatal("Failed to seed users", zap.Error(err))
	}
}

func (s *dataSeederService) isDataSeedingEnabled() bool {
	if !s.AppConfig.Data.EnableDataSeeding {
		s.Logger.Warn("Data seeding is disabled in configuration.")
		return false
	}
	return true
}
