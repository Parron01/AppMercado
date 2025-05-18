package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"github.com/Parron01/AppMercado/backend/pkg/utils"
)

// PriceHistoryService handles business logic for price history
type PriceHistoryService struct {
	priceHistoryRepository *repositories.PriceHistoryRepository
	productService         *ProductService
	userService            *UserService
}

// NewPriceHistoryService creates a new instance of PriceHistoryService
func NewPriceHistoryService(
	priceHistoryRepo *repositories.PriceHistoryRepository,
	productService *ProductService,
	userService *UserService) *PriceHistoryService {
	return &PriceHistoryService{
		priceHistoryRepository: priceHistoryRepo,
		productService:         productService,
		userService:            userService,
	}
}

// createPriceHistory creates a new price history record (internal use only)
func (service *PriceHistoryService) createPriceHistory(productID uint, userID uint, purchaseDate time.Time, purchasePlace string, pricePaid float64) (*models.PriceHistory, error) {
	// Verify if product exists
	product, err := service.productService.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("CreatePriceHistory: produto não encontrado: " + err.Error())
	}

	// Create price history entry
	priceHistory := &models.PriceHistory{
		ProductID:     product.ID,
		UserID:        userID,
		PurchaseDate:  purchaseDate,
		PurchasePlace: purchasePlace,
		PricePaid:     utils.FormatDecimal(pricePaid),
	}

	if err := service.priceHistoryRepository.CreatePriceHistory(priceHistory); err != nil {
		return nil, err
	}

	return priceHistory, nil
}

// GetPriceHistoryByID retrieves a price history record by its ID
func (service *PriceHistoryService) GetPriceHistoryByID(priceHistoryID uint, userID uint, userRole string) (*models.PriceHistory, error) {
	// Get price history
	priceHistory, err := service.priceHistoryRepository.GetPriceHistoryByID(priceHistoryID)
	if err != nil {
		return nil, errors.New("GetPriceHistoryByID: registro de histórico de preço não encontrado")
	}

	// Only admins can see all price history entries
	// Regular users can only see their own entries
	if priceHistory.UserID != userID && userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetPriceHistoryByID: permissão negada: você não pode visualizar registros de histórico de preço de outros usuários")
	}

	return priceHistory, nil
}

// GetPriceHistoryByProductID retrieves price history for a specific product
func (service *PriceHistoryService) GetPriceHistoryByProductID(productID uint) ([]*models.PriceHistory, error) {
	// Verify if product exists
	_, err := service.productService.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("GetPriceHistoryByProductID: produto não encontrado: " + err.Error())
	}

	return service.priceHistoryRepository.GetPriceHistoryByProductID(productID)
}

// GetPriceHistoryByProductAndDateRange retrieves price history for a product in a date range
func (service *PriceHistoryService) GetPriceHistoryByProductAndDateRange(productID uint, startDate, endDate time.Time) ([]*models.PriceHistory, error) {
	// Verify if product exists
	_, err := service.productService.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("GetPriceHistoryByProductAndDateRange: produto não encontrado: " + err.Error())
	}

	return service.priceHistoryRepository.GetPriceHistoryByProductAndDateRange(productID, startDate, endDate)
}

// DeletePriceHistory deletes a price history record
func (service *PriceHistoryService) DeletePriceHistory(priceHistoryID uint, userID uint, userRole string) error {
	// Get price history to check ownership
	priceHistory, err := service.priceHistoryRepository.GetPriceHistoryByID(priceHistoryID)
	if err != nil {
		return errors.New("DeletePriceHistory: registro de histórico de preço não encontrado")
	}

	// Only the creator or admins can delete
	if priceHistory.UserID != userID && userRole != string(models.RoleAdmin) {
		return errors.New("DeletePriceHistory: permissão negada: você não pode excluir registros de histórico de preço de outros usuários")
	}

	return service.priceHistoryRepository.DeletePriceHistory(priceHistoryID)
}

// GetAllPriceHistory retrieves all price history (admin only)
func (service *PriceHistoryService) GetAllPriceHistory(userRole string) ([]*models.PriceHistory, error) {
	// Only admins can see all price history
	if userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetAllPriceHistory: permissão negada: apenas administradores podem listar todos os históricos de preço")
	}

	return service.priceHistoryRepository.GetAllPriceHistory()
}

// RegisterPurchaseInPriceHistory creates price history entries for all items in a purchase
func (service *PriceHistoryService) RegisterPurchaseInPriceHistory(purchase *models.Purchase) error {
	for _, item := range purchase.Items {
		_, err := service.createPriceHistory(
			item.ProductID,
			purchase.UserID,
			purchase.PurchaseDate,
			purchase.PurchaseLocation,
			item.UnitPrice,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// CalculateAveragePriceForProduct calcula o preço médio de um produto baseado em todos os registros de histórico
func (service *PriceHistoryService) CalculateAveragePriceForProduct(productID uint) (float64, error) {
	// Calcular o preço médio usando o repositório
	avgPrice, err := service.priceHistoryRepository.CalculateAveragePriceForProduct(productID)
	if err != nil {
		return 0, err
	}

	// Formatar para quatro casas decimais para cálculos
	return utils.FormatDecimal(avgPrice), nil
}

// GetProductPriceStatistics retrieves price statistics for a product
func (service *PriceHistoryService) GetProductPriceStatistics(productID uint) (*dto.PriceHistoryStatisticsDTO, error) {
	// Verify if product exists
	product, err := service.productService.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("GetProductPriceStatistics: produto não encontrado: " + err.Error())
	}

	// Get statistics
	lowestPrice, highestPrice, firstDate, lastDate, count, err := service.priceHistoryRepository.GetPriceStatisticsByProductID(productID)
	if err != nil {
		return nil, err
	}

	// Calcular o preço médio atual usando o método centralizado
	avgPrice, err := service.CalculateAveragePriceForProduct(productID)
	if err != nil {
		return nil, err
	}

	// Calculate price variation as percentage
	var priceVariation float64 = 0
	if lowestPrice > 0 {
		priceVariation = ((highestPrice - lowestPrice) / lowestPrice) * 100
	}

	// Format dates
	firstDateStr := ""
	lastDateStr := ""
	if !firstDate.IsZero() {
		firstDateStr = firstDate.Format(time.RFC3339)
	}
	if !lastDate.IsZero() {
		lastDateStr = lastDate.Format(time.RFC3339)
	}

	return &dto.PriceHistoryStatisticsDTO{
		ProductID:       product.ID,
		ProductName:     product.Name,
		CurrentAvgPrice: utils.FormatForDisplay(avgPrice), // Preço médio calculado do histórico
		LowestPrice:     utils.FormatForDisplay(lowestPrice),
		HighestPrice:    utils.FormatForDisplay(highestPrice),
		PriceVariation:  utils.FormatForDisplay(priceVariation),
		RecordsCount:    int(count),
		FirstRecordDate: firstDateStr,
		LastRecordDate:  lastDateStr,
	}, nil
}

// ToPriceHistoryResponseDTO converts a PriceHistory model to PriceHistoryResponseDTO
func (service *PriceHistoryService) ToPriceHistoryResponseDTO(priceHistory *models.PriceHistory) dto.PriceHistoryResponseDTO {
	return dto.PriceHistoryResponseDTO{
		ID:            priceHistory.ID,
		ProductID:     priceHistory.ProductID,
		ProductName:   priceHistory.Product.Name,
		UserID:        priceHistory.UserID,
		UserName:      priceHistory.User.Name,
		PurchaseDate:  priceHistory.PurchaseDate.Format(time.RFC3339),
		PurchasePlace: priceHistory.PurchasePlace,
		PricePaid:     utils.FormatForDisplay(priceHistory.PricePaid), // Formatar para exibição
		CreatedAt:     priceHistory.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     priceHistory.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPriceHistoryResponseDTOList converts a list of PriceHistory models to PriceHistoryResponseDTOs
func (service *PriceHistoryService) ToPriceHistoryResponseDTOList(priceHistories []*models.PriceHistory) []dto.PriceHistoryResponseDTO {
	dtos := make([]dto.PriceHistoryResponseDTO, len(priceHistories))
	for i, priceHistory := range priceHistories {
		dtos[i] = service.ToPriceHistoryResponseDTO(priceHistory)
	}
	return dtos
}
