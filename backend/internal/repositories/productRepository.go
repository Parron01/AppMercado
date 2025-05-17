package repositories

import (
	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

// ProductRepository define a interface para operações de banco de dados de produtos
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository cria uma nova instância de ProductRepository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct adiciona um novo produto ao banco de dados
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetProductByID busca um produto pelo ID
func (r *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductByBarcode busca um produto pelo código de barras
func (r *ProductRepository) GetProductByBarcode(barcode string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("barcode = ?", barcode).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProducts retorna todos os produtos
func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct atualiza um produto existente no banco de dados
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

// DeleteProduct remove um produto do banco de dados pelo ID
func (r *ProductRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
