package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username string) (*User, error) {
	user := &User{
		Username: username,
		LastSeen: time.Now(),
	}
	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
