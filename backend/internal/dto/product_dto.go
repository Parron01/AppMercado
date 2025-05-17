package dto

// CreateProductDTO representa os dados para criar um novo produto
type CreateProductDTO struct {
	Name         string  `json:"name" binding:"required"`
	AveragePrice float64 `json:"averagePrice" binding:"omitempty,gte=0"`
	Barcode      string  `json:"barcode" binding:"omitempty"` // Cliente envia "" para vazio ou omite
}

// UpdateProductDTO representa os dados para atualizar um produto existente
type UpdateProductDTO struct {
	Name         *string  `json:"name,omitempty"`
	AveragePrice *float64 `json:"averagePrice,omitempty" binding:"omitempty,gte=0"`
	Barcode      *string  `json:"barcode,omitempty"` // Cliente pode enviar "", null, ou omitir
}

// ProductResponseDTO representa os dados de um produto para resposta HTTP
type ProductResponseDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	AveragePrice float64 `json:"averagePrice"`
	Barcode      *string `json:"barcode"` // Alterado para *string para refletir que pode ser null
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}
