package models

import "gorm.io/gorm"

// Product define a estrutura da tabela de produtos
type Product struct {
	gorm.Model
	Name    string  `gorm:"size:255;not null"`
	Barcode *string `gorm:"size:100;uniqueIndex"` // Alterado para *string para permitir NULL

	// Relacionamentos (serão mais explorados ao criar as tabelas de junção e PriceHistory)
	// UserCategoryProducts []UserCategoryProduct `gorm:"foreignKey:ProductID"`
	// PriceHistories       []PriceHistory        `gorm:"foreignKey:ProductID"`
}
