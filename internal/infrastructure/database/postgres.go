package database

import (
	"log"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Register the explicit join-table model so GORM never tries to
	// drop-and-recreate the implicit post_tags table during AutoMigrate.
	if err := db.SetupJoinTable(&domain.Post{}, "Tags", &domain.PostTag{}); err != nil {
		log.Fatalf("Failed to setup join table: %v", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Tag{},
		&domain.PostTag{},
		&domain.Post{},
		&domain.Follow{},
		&domain.Like{},
		&domain.Comment{},
		&domain.Message{},
		&domain.Bookmark{},
		&domain.Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}
