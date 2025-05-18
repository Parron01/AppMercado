package repositories

import (
	"time"

	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

// PriceHistoryRepository handles database operations for price history
type PriceHistoryRepository struct {
	database *gorm.DB
}

// NewPriceHistoryRepository creates a new instance of PriceHistoryRepository
func NewPriceHistoryRepository(db *gorm.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{database: db}
}

// CreatePriceHistory adds a new price history record to the database
func (repo *PriceHistoryRepository) CreatePriceHistory(priceHistory *models.PriceHistory) error {
	return repo.database.Create(priceHistory).Error
}

// GetPriceHistoryByID retrieves a price history record by its ID
func (repo *PriceHistoryRepository) GetPriceHistoryByID(id uint) (*models.PriceHistory, error) {
	var priceHistory models.PriceHistory
	if err := repo.database.Preload("Product").Preload("User").First(&priceHistory, id).Error; err != nil {
		return nil, err
	}
	return &priceHistory, nil
}

// GetPriceHistoryByProductID retrieves price history records for a specific product
func (repo *PriceHistoryRepository) GetPriceHistoryByProductID(productID uint) ([]*models.PriceHistory, error) {
	var priceHistories []*models.PriceHistory
	if err := repo.database.Where("product_id = ?", productID).
		Preload("Product").
		Preload("User").
		Order("purchase_date desc").
		Find(&priceHistories).Error; err != nil {
		return nil, err
	}
	return priceHistories, nil
}

// GetPriceHistoryByUserID retrieves price history records for a specific user
func (repo *PriceHistoryRepository) GetPriceHistoryByUserID(userID uint) ([]*models.PriceHistory, error) {
	var priceHistories []*models.PriceHistory
	if err := repo.database.Where("user_id = ?", userID).
		Preload("Product").
		Preload("User").
		Order("purchase_date desc").
		Find(&priceHistories).Error; err != nil {
		return nil, err
	}
	return priceHistories, nil
}

// GetPriceHistoryByProductAndDateRange retrieves price history records for a specific product in a date range
func (repo *PriceHistoryRepository) GetPriceHistoryByProductAndDateRange(
	productID uint, startDate, endDate time.Time) ([]*models.PriceHistory, error) {
	var priceHistories []*models.PriceHistory
	if err := repo.database.Where("product_id = ? AND purchase_date BETWEEN ? AND ?",
		productID, startDate, endDate).
		Preload("Product").
		Preload("User").
		Order("purchase_date desc").
		Find(&priceHistories).Error; err != nil {
		return nil, err
	}
	return priceHistories, nil
}

// UpdatePriceHistory updates an existing price history record in the database
func (repo *PriceHistoryRepository) UpdatePriceHistory(priceHistory *models.PriceHistory) error {
	return repo.database.Save(priceHistory).Error
}

// DeletePriceHistory removes a price history record from the database
func (repo *PriceHistoryRepository) DeletePriceHistory(id uint) error {
	return repo.database.Delete(&models.PriceHistory{}, id).Error
}

// GetAllPriceHistory retrieves all price history records
func (repo *PriceHistoryRepository) GetAllPriceHistory() ([]*models.PriceHistory, error) {
	var priceHistories []*models.PriceHistory
	if err := repo.database.Preload("Product").Preload("User").Find(&priceHistories).Error; err != nil {
		return nil, err
	}
	return priceHistories, nil
}

// GetPriceStatisticsByProductID retrieves price statistics for a product
func (repo *PriceHistoryRepository) GetPriceStatisticsByProductID(productID uint) (
	lowestPrice float64, highestPrice float64, firstDate time.Time, lastDate time.Time, count int64, err error) {

	// Get lowest price
	err = repo.database.Model(&models.PriceHistory{}).
		Where("product_id = ?", productID).
		Select("MIN(price_paid)").
		Row().
		Scan(&lowestPrice)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, 0, err
	}

	// Get highest price
	err = repo.database.Model(&models.PriceHistory{}).
		Where("product_id = ?", productID).
		Select("MAX(price_paid)").
		Row().
		Scan(&highestPrice)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, 0, err
	}

	// Get earliest date
	err = repo.database.Model(&models.PriceHistory{}).
		Where("product_id = ?", productID).
		Order("purchase_date ASC").
		Limit(1).
		Select("purchase_date").
		Row().
		Scan(&firstDate)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, 0, err
	}

	// Get latest date
	err = repo.database.Model(&models.PriceHistory{}).
		Where("product_id = ?", productID).
		Order("purchase_date DESC").
		Limit(1).
		Select("purchase_date").
		Row().
		Scan(&lastDate)
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, 0, err
	}

	// Get count of records
	err = repo.database.Model(&models.PriceHistory{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	if err != nil {
		return 0, 0, time.Time{}, time.Time{}, 0, err
	}

	return lowestPrice, highestPrice, firstDate, lastDate, count, nil
}

// CalculateAveragePriceForProduct calcula o preço médio de um produto baseado em todos os registros de histórico
func (repo *PriceHistoryRepository) CalculateAveragePriceForProduct(productID uint) (float64, error) {
	var result struct {
		AvgPrice float64
	}

	// Utilizamos o banco de dados para calcular a média diretamente, preservando precisão máxima
	// Nota: CAST para decimal(10,4) garante consistência nos tipos
	err := repo.database.Model(&models.PriceHistory{}).
		Select("COALESCE(AVG(CAST(price_paid AS DECIMAL(10,4))), 0.0) as avg_price").
		Where("product_id = ?", productID).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.AvgPrice, nil
}
