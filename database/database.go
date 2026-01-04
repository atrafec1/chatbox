package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	URL := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = MigrateDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Group{}, &Message{})
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}
	return nil
}
