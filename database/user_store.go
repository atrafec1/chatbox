package database

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	var group Group
	if err := db.Where("name = ?", "general").First(&group).Error; err != nil {
		return nil, fmt.Errorf("failed to find default group: %w", err)
	}

	user := &User{
		Username: username,
		Password: string(hashedPassword),
		LastSeen: time.Now(),
		GroupID:  group.ID,
	}
	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
func UsernameExists(db *gorm.DB, username string) (bool, error) {
	var count int64
	if err := db.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

func Login(db *gorm.DB, username, password string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if error := CheckPassword(user.Password, password); error != nil {
		return nil, fmt.Errorf("invalid password")
	}

	user.LastSeen = time.Now()
	if err := db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update last seen: %w", err)
	}
	return &user, nil
}

func CheckPassword(hashPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(rawPassword))
}
