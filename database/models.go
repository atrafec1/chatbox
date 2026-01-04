package database

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	LastSeen  time.Time

	Messages []Message `gorm:"foreignKey:UserID"`
	Groups   []Group   `gorm:"many2many:user_groups;"`
}

type Group struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	Users    []User    `gorm:"many2many:user_groups;"`
	Messages []Message `gorm:"foreignKey:GroupID"`
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	GroupID   uint
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
