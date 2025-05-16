package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string `gorm:"size:100;not null"`
	UserID uint   `gorm:"not null"`
	User   User   `gorm:"foreignKey:UserID"`
}
