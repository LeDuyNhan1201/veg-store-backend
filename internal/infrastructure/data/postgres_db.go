package data

import (
	"fmt"
	"log"
	"time"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/util"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	*core.Core
	*gorm.DB
}

func InitPostgresDB(core *core.Core) *PostgresDB {
	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		core.AppConfig.Data.Postgres.Host,
		core.AppConfig.Data.Postgres.Port,
		core.AppConfig.Data.Postgres.User,
		core.AppConfig.Data.Postgres.Password,
		core.AppConfig.Data.Postgres.DBName,
		core.AppConfig.Data.Postgres.SSL.Mode,
	)

	// Add SSL cert paths if present
	certsPath := util.GetConfigPathFromGoMod("secrets/certs")
	if core.AppConfig.Data.Postgres.SSL.RootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s/%s", certsPath, core.AppConfig.Data.Postgres.SSL.RootCert)
	}
	if core.AppConfig.Data.Postgres.SSL.Cert != "" {
		dsn += fmt.Sprintf(" sslcert=%s/%s", certsPath, core.AppConfig.Data.Postgres.SSL.Cert)
	}
	if core.AppConfig.Data.Postgres.SSL.Key != "" {
		dsn += fmt.Sprintf(" sslkey=%s/%s", certsPath, core.AppConfig.Data.Postgres.SSL.Key)
	}

	// Open GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // optional: log SQL queries
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresDB{
		Core: core,
		DB:   db,
	}
}

func (pdb *PostgresDB) Migrate() error {
	models := []interface{}{
		&model.User{},
		&model.TaskStatus{},
		&model.Task{},
		&model.Product{},
		&model.Order{},
		&model.OrderDetail{},
		// Add other models here
	}

	switch pdb.AppConfig.Data.Postgres.DDLMode {
	case "create-drop":
		pdb.Logger.Info("Dropping all tables...")
		err := pdb.DB.Migrator().DropTable(models...)
		if err != nil {
			pdb.Logger.Fatal("Failed to drop tables: ", zap.Error(err))
		}
		return pdb.DB.AutoMigrate(models...)

	case "update":
		pdb.Logger.Info("Updating database schema...")
		return pdb.DB.AutoMigrate(models...)

	default:
		pdb.Logger.Info("No DDL operations will be performed.")
		return nil
	}
}
