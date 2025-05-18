package dto

// CreateProductDTO representa os dados para criar um novo produto
type CreateProductDTO struct {
	Name    string `json:"name" binding:"required"`
	Barcode string `json:"barcode" binding:"omitempty"` // Cliente envia "" para vazio ou omite
}

// UpdateProductDTO representa os dados para atualizar um produto existente
type UpdateProductDTO struct {
	Name    *string `json:"name,omitempty"`
	Barcode *string `json:"barcode,omitempty"` // Cliente pode enviar "", null, ou omitir
}

// ProductResponseDTO representa os dados de um produto para resposta HTTP
type ProductResponseDTO struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Barcode   *string `json:"barcode"` // Alterado para *string para refletir que pode ser null
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
