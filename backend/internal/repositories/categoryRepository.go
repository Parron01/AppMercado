package repositories

import (
	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	database *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{database: db}
}

func (repository *CategoryRepository) CreateCategory(category *models.Category) error {
	return repository.database.Create(category).Error
}

func (repository *CategoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := repository.database.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (repository *CategoryRepository) GetCategoriesByUserID(userID uint) ([]*models.Category, error) {
	var categories []*models.Category
	if err := repository.database.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repository *CategoryRepository) UpdateCategory(category *models.Category) error {
	return repository.database.Save(category).Error
}

func (repository *CategoryRepository) DeleteCategory(id uint) error {
	return repository.database.Delete(&models.Category{}, id).Error
}

func (repository *CategoryRepository) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := repository.database.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
