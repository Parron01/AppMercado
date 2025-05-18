package models

import (
	"time"

	"gorm.io/gorm"
)

// PriceHistory defines the structure for product price history records
type PriceHistory struct {
	gorm.Model
	ProductID     uint      `gorm:"not null;index:idx_price_history_product"`
	Product       Product   `gorm:"foreignKey:ProductID"`
	UserID        uint      `gorm:"not null;index:idx_price_history_user"`
	User          User      `gorm:"foreignKey:UserID"`
	PurchaseDate  time.Time `gorm:"not null;index:idx_price_history_date"`
	PurchasePlace string    `gorm:"size:255"`                    // Store where the product was purchased
	PricePaid     float64   `gorm:"type:decimal(10,4);not null"` // Aumentado para decimal(10,4)
}
