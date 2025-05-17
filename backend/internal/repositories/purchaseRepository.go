package repositories

import (
	"time"

	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

// PurchaseRepository handles database operations for purchases
type PurchaseRepository struct {
	database *gorm.DB
}

// NewPurchaseRepository creates a new instance of PurchaseRepository
func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{database: db}
}

// GetPurchaseByDateAndLocation busca uma compra pelo local e data
func (repo *PurchaseRepository) GetPurchaseByDateAndLocation(date time.Time, location string, userID uint) (*models.Purchase, error) {
	var purchase models.Purchase

	// Precisamos definir um intervalo de tempo pequeno (1 segundo) para a verificação,
	// já que a comparação exata de timestamps pode ser problemática
	startTime := date.Add(-1 * time.Second)
	endTime := date.Add(1 * time.Second)

	err := repo.database.Where("purchase_location = ? AND purchase_date BETWEEN ? AND ? AND user_id = ?",
		location, startTime, endTime, userID).First(&purchase).Error

	return &purchase, err
}

// CreatePurchase adds a new purchase to the database
func (repo *PurchaseRepository) CreatePurchase(purchase *models.Purchase) error {
	if err := repo.database.Create(purchase).Error; err != nil {
		return err
	}
	// Reload the purchase with its relationships after creation
	return repo.database.Preload("Items.Product").First(purchase, purchase.ID).Error
}

// GetPurchaseByID retrieves a purchase by its ID
func (repo *PurchaseRepository) GetPurchaseByID(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	if err := repo.database.Preload("Items.Product").First(&purchase, id).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}

// GetPurchasesByUserID retrieves all purchases for a specific user
func (repo *PurchaseRepository) GetPurchasesByUserID(userID uint) ([]*models.Purchase, error) {
	var purchases []*models.Purchase
	if err := repo.database.Preload("Items.Product").Where("user_id = ?", userID).Find(&purchases).Error; err != nil {
		return nil, err
	}
	return purchases, nil
}

// UpdatePurchase updates an existing purchase in the database
func (repo *PurchaseRepository) UpdatePurchase(purchase *models.Purchase) error {
	return repo.database.Save(purchase).Error
}

// DeletePurchase removes a purchase from the database
func (repo *PurchaseRepository) DeletePurchase(id uint) error {
	return repo.database.Transaction(func(tx *gorm.DB) error {
		// Delete associated purchase items first
		if err := tx.Where("purchase_id = ?", id).Delete(&models.PurchaseItem{}).Error; err != nil {
			return err
		}
		// Then delete the purchase
		if err := tx.Delete(&models.Purchase{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAllPurchases retrieves all purchases (admin only)
func (repo *PurchaseRepository) GetAllPurchases() ([]*models.Purchase, error) {
	var purchases []*models.Purchase
	if err := repo.database.Preload("Items.Product").Find(&purchases).Error; err != nil {
		return nil, err
	}
	return purchases, nil
}
