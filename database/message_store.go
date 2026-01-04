package database

import (
	"fmt"

	"gorm.io/gorm"
)

func SaveMessage(db *gorm.DB, msg string, userId, groupId uint) error {
	message := &Message{
		UserID:  userId,
		GroupID: groupId,
		Content: msg,
	}
	if err := db.Create(message).Error; err != nil {
		return fmt.Errorf("error saving message to db: %w", err)
	}
	return nil
}
