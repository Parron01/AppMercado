package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name         string `gorm:"size:100"`
    Email        string `gorm:"uniqueIndex;size:100"`
    PasswordHash string `gorm:"size:255"`
    Role         string `gorm:"size:20"` // Admin, Standard, Guest
}