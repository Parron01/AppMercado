package dto

type CreateCategoryDTO struct {
	Name string `json:"name" binding:"required" example:"Frutas"`
}

type UpdateCategoryDTO struct {
	Name string `json:"name" binding:"required" example:"Legumes"`
}

type CategoryResponseDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	UserID    uint   `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
