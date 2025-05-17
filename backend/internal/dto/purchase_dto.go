package dto

import "time"

// PurchaseItemDTO represents an item in a purchase
type PurchaseItemDTO struct {
	ProductID uint    `json:"productId" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unitPrice" binding:"required,gt=0"` // Renomeado de PricePaid para UnitPrice para maior clareza
}

// CreatePurchaseDTO represents data needed to create a purchase
type CreatePurchaseDTO struct {
	PurchaseDate     time.Time         `json:"purchaseDate" binding:"required"`
	PurchaseLocation string            `json:"purchaseLocation" binding:"required"`
	Items            []PurchaseItemDTO `json:"items" binding:"required,dive"`
}

// UpdatePurchaseDTO represents data needed to update a purchase
type UpdatePurchaseDTO struct {
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	PurchaseLocation *string            `json:"purchaseLocation,omitempty"`
	Items            *[]PurchaseItemDTO `json:"items,omitempty" binding:"omitempty,dive"`
}

// PurchaseItemResponseDTO represents the response data for a purchase item
type PurchaseItemResponseDTO struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	TotalPrice  float64 `json:"totalPrice"` // Renomeado de PricePaid para TotalPrice para maior clareza
	// Removido o campo subtotal por ser redundante com totalPrice
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// PurchaseResponseDTO represents the response data for a purchase
type PurchaseResponseDTO struct {
	ID               uint                      `json:"id"`
	PurchaseDate     string                    `json:"purchaseDate"`
	PurchaseLocation string                    `json:"purchaseLocation"`
	UserID           uint                      `json:"userId"`
	Items            []PurchaseItemResponseDTO `json:"items"`
	Total            float64                   `json:"total"`
	CreatedAt        string                    `json:"createdAt"`
	UpdatedAt        string                    `json:"updatedAt"`
}
