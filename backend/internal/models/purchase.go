package models

import (
	"time"

	"gorm.io/gorm"
)

// Purchase defines the structure for purchase records
type Purchase struct {
	gorm.Model
	PurchaseDate     time.Time `gorm:"not null;index:idx_purchase_date_location_user"`
	PurchaseLocation string    `gorm:"size:255;index:idx_purchase_date_location_user"`
	UserID           uint      `gorm:"not null;index:idx_purchase_date_location_user"`
	User             User      `gorm:"foreignKey:UserID"`

	// Relationships
	Items []PurchaseItem `gorm:"foreignKey:PurchaseID"`
	Total float64        `gorm:"type:decimal(10,4)"` // Aumentado para decimal(10,4)
}
