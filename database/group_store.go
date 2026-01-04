package database

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateGroup(db *gorm.DB, name string) (*Group, error) {
	group := &Group{
		Name: name,
	}

	if err := db.Create(group).Error; err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}
	return group, nil
}

func (g *Group) AddUser(db *gorm.DB, user User) error {
	return db.Model(g).Association("Users").Append(user)
}
