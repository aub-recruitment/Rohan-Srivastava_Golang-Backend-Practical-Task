package infrastructure

import (
	"fmt"
	"log"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var gormLogger logger.Interface
	if cfg.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	log.Println("Database connection established")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Content{},
		&domain.Plan{},
		&domain.Subscription{},
		&domain.WatchHistory{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}
