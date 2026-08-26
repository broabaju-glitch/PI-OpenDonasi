package config

import (
	"log"
	"os"

	"opendonasi-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/opendonasi?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Campaign{},
		&models.Donation{},
		&models.Report{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}

	return db, nil
}
