package models

import "gorm.io/gorm"

// UserCategoryProduct represents a relationship between a user, a category, and a product
type UserCategoryProduct struct {
	gorm.Model
	UserID     uint     `gorm:"not null;index:idx_user_category_product"`
	CategoryID uint     `gorm:"not null;index:idx_user_category_product"`
	ProductID  uint     `gorm:"not null;index:idx_user_category_product"`
	User       User     `gorm:"foreignKey:UserID"`
	Category   Category `gorm:"foreignKey:CategoryID"`
	Product    Product  `gorm:"foreignKey:ProductID"`
}
