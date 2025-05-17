package models

import "gorm.io/gorm"

// PurchaseItem represents a product in a purchase with its quantity and price
type PurchaseItem struct {
	gorm.Model
	PurchaseID uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   float64 `gorm:"type:decimal(10,2);not null"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null"` // Preço por unidade
	TotalPrice float64 `gorm:"type:decimal(10,2);not null"` // Preço total do item (quantidade * preço unitário)

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
