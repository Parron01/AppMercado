package models

import "gorm.io/gorm"

// PurchaseItem represents a product in a purchase with its quantity and price
type PurchaseItem struct {
	gorm.Model
	PurchaseID uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   float64 `gorm:"type:decimal(10,4);not null"` // Aumentado para decimal(10,4)
	UnitPrice  float64 `gorm:"type:decimal(10,4);not null"` // Aumentado para decimal(10,4)
	TotalPrice float64 `gorm:"type:decimal(10,4);not null"` // Aumentado para decimal(10,4)

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
