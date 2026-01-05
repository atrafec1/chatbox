package database

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	LastSeen  time.Time

	Messages []Message `gorm:"foreignKey:UserID"`
	Group    Group     `gorm:"foreignKey:GroupID"`
	GroupID  uint
}

type Group struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	Users    []User    `gorm:"foreignKey:GroupID"`
	Messages []Message `gorm:"foreignKey:GroupID"`
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	GroupID   uint
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
