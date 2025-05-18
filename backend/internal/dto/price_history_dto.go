package dto

import "time"

// CreatePriceHistoryDTO represents data needed to manually create a price history record
type CreatePriceHistoryDTO struct {
	ProductID     uint      `json:"productId" binding:"required"`
	PurchaseDate  time.Time `json:"purchaseDate" binding:"required"`
	PurchasePlace string    `json:"purchasePlace" binding:"required"`
	PricePaid     float64   `json:"pricePaid" binding:"required,gt=0"`
}

// PriceHistoryResponseDTO represents the response data for a price history record
type PriceHistoryResponseDTO struct {
	ID            uint    `json:"id"`
	ProductID     uint    `json:"productId"`
	ProductName   string  `json:"productName"`
	UserID        uint    `json:"userId"`
	UserName      string  `json:"userName"`
	PurchaseDate  string  `json:"purchaseDate"`
	PurchasePlace string  `json:"purchasePlace"`
	PricePaid     float64 `json:"pricePaid"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

// PriceHistoryStatisticsDTO represents statistical data about a product's price history
type PriceHistoryStatisticsDTO struct {
	ProductID       uint    `json:"productId"`
	ProductName     string  `json:"productName"`
	CurrentAvgPrice float64 `json:"currentAvgPrice"`
	LowestPrice     float64 `json:"lowestPrice"`
	HighestPrice    float64 `json:"highestPrice"`
	PriceVariation  float64 `json:"priceVariation"` // Percentage variation between lowest and highest
	RecordsCount    int     `json:"recordsCount"`
	FirstRecordDate string  `json:"firstRecordDate"`
	LastRecordDate  string  `json:"lastRecordDate"`
}
