package dto

// CreateUserCategoryProductDTO represents data needed to create a user-category-product relationship
type CreateUserCategoryProductDTO struct {
	CategoryID uint `json:"categoryId" binding:"required"`
	ProductID  uint `json:"productId" binding:"required"`
}

// UserCategoryProductResponseDTO represents the response data for a user-category-product relationship
type UserCategoryProductResponseDTO struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"userId"`
	UserName     string `json:"userName"`
	CategoryID   uint   `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	ProductID    uint   `json:"productId"`
	ProductName  string `json:"productName"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
