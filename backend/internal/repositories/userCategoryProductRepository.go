package repositories

import (
	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

// UserCategoryProductRepository handles database operations for user-category-product relationships
type UserCategoryProductRepository struct {
	database *gorm.DB
}

// NewUserCategoryProductRepository creates a new instance of UserCategoryProductRepository
func NewUserCategoryProductRepository(db *gorm.DB) *UserCategoryProductRepository {
	return &UserCategoryProductRepository{database: db}
}

// CreateUserCategoryProduct adds a new user-category-product relationship to the database
func (repo *UserCategoryProductRepository) CreateUserCategoryProduct(ucp *models.UserCategoryProduct) error {
	return repo.database.Create(ucp).Error
}

// GetUserCategoryProductByID retrieves a user-category-product relationship by its ID
func (repo *UserCategoryProductRepository) GetUserCategoryProductByID(id uint) (*models.UserCategoryProduct, error) {
	var ucp models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").First(&ucp, id).Error; err != nil {
		return nil, err
	}
	return &ucp, nil
}

// GetUserCategoryProductsByUserID retrieves all user-category-product relationships for a specific user
func (repo *UserCategoryProductRepository) GetUserCategoryProductsByUserID(userID uint) ([]*models.UserCategoryProduct, error) {
	var ucps []*models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").Where("user_id = ?", userID).Find(&ucps).Error; err != nil {
		return nil, err
	}
	return ucps, nil
}

// GetUserCategoryProductsByCategoryID retrieves all user-category-product relationships for a specific category
func (repo *UserCategoryProductRepository) GetUserCategoryProductsByCategoryID(categoryID uint) ([]*models.UserCategoryProduct, error) {
	var ucps []*models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").Where("category_id = ?", categoryID).Find(&ucps).Error; err != nil {
		return nil, err
	}
	return ucps, nil
}

// GetUserCategoryProductsByProductID retrieves all user-category-product relationships for a specific product
func (repo *UserCategoryProductRepository) GetUserCategoryProductsByProductID(productID uint) ([]*models.UserCategoryProduct, error) {
	var ucps []*models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").Where("product_id = ?", productID).Find(&ucps).Error; err != nil {
		return nil, err
	}
	return ucps, nil
}

// GetUserCategoryProduct retrieves a specific user-category-product relationship
func (repo *UserCategoryProductRepository) GetUserCategoryProduct(userID, categoryID, productID uint) (*models.UserCategoryProduct, error) {
	var ucp models.UserCategoryProduct
	if err := repo.database.Where("user_id = ? AND category_id = ? AND product_id = ?",
		userID, categoryID, productID).First(&ucp).Error; err != nil {
		return nil, err
	}
	return &ucp, nil
}

// ReloadUserCategoryProductWithRelations recarrega o objeto com todos os seus relacionamentos
func (repo *UserCategoryProductRepository) ReloadUserCategoryProductWithRelations(id uint) (*models.UserCategoryProduct, error) {
	var ucp models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").First(&ucp, id).Error; err != nil {
		return nil, err
	}
	return &ucp, nil
}

// DeleteUserCategoryProduct deletes a user-category-product relationship
func (repo *UserCategoryProductRepository) DeleteUserCategoryProduct(id uint) error {
	return repo.database.Delete(&models.UserCategoryProduct{}, id).Error
}

// DeleteUserCategoryProductByFields deletes a user-category-product relationship by its composite fields
func (repo *UserCategoryProductRepository) DeleteUserCategoryProductByFields(userID, categoryID, productID uint) error {
	return repo.database.Where("user_id = ? AND category_id = ? AND product_id = ?",
		userID, categoryID, productID).Delete(&models.UserCategoryProduct{}).Error
}

// GetAllUserCategoryProducts retrieves all user-category-product relationships
func (repo *UserCategoryProductRepository) GetAllUserCategoryProducts() ([]*models.UserCategoryProduct, error) {
	var ucps []*models.UserCategoryProduct
	if err := repo.database.Preload("User").Preload("Category").Preload("Product").Find(&ucps).Error; err != nil {
		return nil, err
	}
	return ucps, nil
}
