package database

import (
	"checkjop-be/internal/config"
	"checkjop-be/internal/model"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	const maxRetries = 10
	const retryInterval = 2 * time.Second

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
		if err == nil {
			// test ping
			sqlDB, pingErr := db.DB()
			if pingErr == nil && sqlDB.Ping() == nil {
				log.Printf("✅ Connected to database on attempt %d\n", i)
				break
			}
		}

		log.Printf("❌ Failed to connect to database (attempt %d/%d): %v\n", i, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("🔥 could not connect to database after %d retries: %w", maxRetries, err)
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&model.Curriculum{},
		&model.Category{},
		&model.Course{},
		&model.PrerequisiteGroup{},
		&model.PrerequisiteCourseLink{},
		&model.SetDefault{},
	); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
