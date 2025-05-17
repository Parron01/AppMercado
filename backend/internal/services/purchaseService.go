package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"gorm.io/gorm"
)

// PurchaseService handles business logic for purchases
type PurchaseService struct {
	purchaseRepository *repositories.PurchaseRepository
	productService     *ProductService
}

// NewPurchaseService creates a new instance of PurchaseService
func NewPurchaseService(purchaseRepo *repositories.PurchaseRepository, productService *ProductService) *PurchaseService {
	return &PurchaseService{
		purchaseRepository: purchaseRepo,
		productService:     productService,
	}
}

// CreatePurchase creates a new purchase with its items
func (service *PurchaseService) CreatePurchase(purchaseDTO dto.CreatePurchaseDTO, userID uint) (*models.Purchase, error) {
	// Basic validation
	if len(purchaseDTO.Items) == 0 {
		return nil, errors.New("CreatePurchase: pelo menos um item é necessário")
	}

	// Verificar se já existe uma compra com mesmo local e data
	_, err := service.purchaseRepository.GetPurchaseByDateAndLocation(
		purchaseDTO.PurchaseDate,
		purchaseDTO.PurchaseLocation,
		userID)

	if err == nil {
		// Se não houver erro, significa que encontramos uma compra com os mesmos dados
		return nil, errors.New("CreatePurchase: já existe uma compra com o mesmo local e data")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Se o erro não for "registro não encontrado", é um erro de banco de dados
		return nil, err
	}

	// Create the purchase
	purchase := &models.Purchase{
		PurchaseDate:     purchaseDTO.PurchaseDate,
		PurchaseLocation: purchaseDTO.PurchaseLocation,
		UserID:           userID,
		Items:            make([]models.PurchaseItem, len(purchaseDTO.Items)),
		Total:            0,
	}

	// Map para acumular valores para atualização de preço médio
	productPriceUpdates := make(map[uint]struct {
		TotalQuantity float64
		TotalSpent    float64
	})

	// Add items to the purchase
	var total float64 = 0
	for i, itemDTO := range purchaseDTO.Items {
		// Get product to check if it exists
		_, err := service.productService.GetProductByID(itemDTO.ProductID)
		if err != nil {
			return nil, errors.New("CreatePurchase: produto não encontrado: " + err.Error())
		}

		// Usar o preço informado pelo usuário como preço unitário
		unitPrice := formatDecimal(itemDTO.UnitPrice)
		// Calcular o preço total do item
		totalPrice := formatDecimal(itemDTO.Quantity * unitPrice)

		// Create purchase item with the provided price
		purchase.Items[i] = models.PurchaseItem{
			ProductID:  itemDTO.ProductID,
			Quantity:   itemDTO.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: totalPrice,
		}
		total += totalPrice

		// Acumular dados para recálculo de preço médio
		if update, exists := productPriceUpdates[itemDTO.ProductID]; exists {
			update.TotalQuantity += itemDTO.Quantity
			update.TotalSpent += totalPrice
			productPriceUpdates[itemDTO.ProductID] = update
		} else {
			productPriceUpdates[itemDTO.ProductID] = struct {
				TotalQuantity float64
				TotalSpent    float64
			}{
				TotalQuantity: itemDTO.Quantity,
				TotalSpent:    totalPrice,
			}
		}
	}

	purchase.Total = formatDecimal(total)

	// Save to database
	if err := service.purchaseRepository.CreatePurchase(purchase); err != nil {
		return nil, err
	}

	// Atualizar preços médios dos produtos
	for productID, update := range productPriceUpdates {
		service.updateProductAveragePrice(productID, update.TotalQuantity, update.TotalSpent)
	}

	return purchase, nil
}

// updateProductAveragePrice updates the average price of a product based on purchases
func (service *PurchaseService) updateProductAveragePrice(productID uint, newQuantity float64, newTotal float64) error {
	product, err := service.productService.GetProductByID(productID)
	if err != nil {
		return err
	}

	// Calcular o preço unitário dessa compra
	newUnitPrice := formatDecimal(newTotal / newQuantity)

	var newAveragePrice float64

	// Se o produto ainda não tem preço médio (primeira compra), use o preço desta compra
	if product.AveragePrice == 0 {
		newAveragePrice = newUnitPrice
	} else {
		// Implementação de média ponderada
		// Fórmula: (preço_médio_atual * 80% + novo_preço * 20%)
		// Isso dá mais peso ao histórico e evita flutuações bruscas

		historicalWeight := 0.8  // 80% de peso para o histórico
		newPurchaseWeight := 0.2 // 20% de peso para a nova compra

		newAveragePrice = formatDecimal((product.AveragePrice * historicalWeight) +
			(newUnitPrice * newPurchaseWeight))
	}

	// Atualizar o preço médio do produto no banco de dados
	updateDTO := dto.UpdateProductDTO{
		AveragePrice: &newAveragePrice,
	}

	_, err = service.productService.UpdateProduct(productID, updateDTO)
	return err
}

// GetPurchaseByID retrieves a purchase by its ID
func (service *PurchaseService) GetPurchaseByID(purchaseID uint, userID uint, userRole string) (*models.Purchase, error) {
	// Get purchase
	purchase, err := service.purchaseRepository.GetPurchaseByID(purchaseID)
	if err != nil {
		return nil, errors.New("GetPurchaseByID: compra não encontrada")
	}

	// Check if user has permission to view this purchase
	if purchase.UserID != userID && userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetPurchaseByID: permissão negada: você não pode visualizar compras de outros usuários")
	}

	return purchase, nil
}

// GetPurchasesByUserID retrieves all purchases for a specific user
func (service *PurchaseService) GetPurchasesByUserID(userID uint) ([]*models.Purchase, error) {
	return service.purchaseRepository.GetPurchasesByUserID(userID)
}

// DeletePurchase deletes a purchase and its items
func (service *PurchaseService) DeletePurchase(purchaseID uint, userID uint, userRole string) error {
	// Get purchase with all its items and product information
	purchase, err := service.purchaseRepository.GetPurchaseByID(purchaseID)
	if err != nil {
		return errors.New("DeletePurchase: compra não encontrada")
	}

	// Check if user has permission to delete this purchase
	if purchase.UserID != userID && userRole != string(models.RoleAdmin) {
		return errors.New("DeletePurchase: permissão negada: você não pode excluir compras de outros usuários")
	}

	// Mapa para ajustar os preços médios após exclusão da compra
	productAdjustments := make(map[uint]struct {
		Quantity  float64
		TotalPaid float64
	})

	// Calcular o impacto da remoção de cada item
	for _, item := range purchase.Items {
		if adj, exists := productAdjustments[item.ProductID]; exists {
			adj.Quantity += item.Quantity
			adj.TotalPaid += item.TotalPrice
			productAdjustments[item.ProductID] = adj
		} else {
			productAdjustments[item.ProductID] = struct {
				Quantity  float64
				TotalPaid float64
			}{
				Quantity:  item.Quantity,
				TotalPaid: item.TotalPrice,
			}
		}
	}

	// Delete from database
	if err := service.purchaseRepository.DeletePurchase(purchaseID); err != nil {
		return err
	}

	// Ajustar os preços médios após a exclusão
	for productID, adj := range productAdjustments {
		// Buscar o produto
		product, err := service.productService.GetProductByID(productID)
		if err != nil {
			continue // Se não conseguir encontrar o produto, continua para o próximo
		}

		// Calcular o preço unitário deste item
		unitPrice := adj.TotalPaid / adj.Quantity

		// Reverter o efeito deste preço no preço médio
		// Na adição: avg = old * 0.8 + new * 0.2
		// Na remoção, calculamos:
		// 1. oldAvg = (currentAvg - (unitPrice * 0.2)) / 0.8
		var newAveragePrice float64

		// Se após a exclusão desta compra o preço médio ficaria negativo ou zero,
		// definimos para 0 (produto sem histórico de preço)
		if product.AveragePrice <= (unitPrice * 0.2 / 0.8) {
			newAveragePrice = 0 // Reset do preço médio
		} else {
			// Reverter o cálculo da média ponderada
			historicalWeight := 0.8
			newPurchaseWeight := 0.2

			// Fórmula inversa da que usamos para adicionar
			oldAverage := (product.AveragePrice - (unitPrice * newPurchaseWeight)) / historicalWeight

			// Formatar para 2 casas decimais
			newAveragePrice = formatDecimal(oldAverage)
		}

		// Atualizar o produto
		updateDTO := dto.UpdateProductDTO{
			AveragePrice: &newAveragePrice,
		}

		service.productService.UpdateProduct(productID, updateDTO)
	}

	return nil
}

// GetAllPurchases retrieves all purchases (admin only)
func (service *PurchaseService) GetAllPurchases(userRole string) ([]*models.Purchase, error) {
	// Check if user is admin
	if userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetAllPurchases: permissão negada: apenas administradores podem listar todas as compras")
	}

	return service.purchaseRepository.GetAllPurchases()
}

// formatDecimal formats a float64 to have exactly 2 decimal places
func formatDecimal(value float64) float64 {
	return float64(int64(value*100+0.5)) / 100
}

// ToPurchaseItemResponseDTO converts a PurchaseItem model to PurchaseItemResponseDTO
func (service *PurchaseService) ToPurchaseItemResponseDTO(item models.PurchaseItem) dto.PurchaseItemResponseDTO {
	return dto.PurchaseItemResponseDTO{
		ID:          item.ID,
		ProductID:   item.ProductID,
		ProductName: item.Product.Name,
		Quantity:    formatDecimal(item.Quantity),
		UnitPrice:   formatDecimal(item.UnitPrice),
		TotalPrice:  formatDecimal(item.TotalPrice),
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPurchaseResponseDTO converts a Purchase model to PurchaseResponseDTO
func (service *PurchaseService) ToPurchaseResponseDTO(purchase *models.Purchase) dto.PurchaseResponseDTO {
	itemDTOs := make([]dto.PurchaseItemResponseDTO, len(purchase.Items))
	for i, item := range purchase.Items {
		itemDTOs[i] = service.ToPurchaseItemResponseDTO(item)
	}

	return dto.PurchaseResponseDTO{
		ID:               purchase.ID,
		PurchaseDate:     purchase.PurchaseDate.Format(time.RFC3339),
		PurchaseLocation: purchase.PurchaseLocation,
		UserID:           purchase.UserID,
		Items:            itemDTOs,
		Total:            formatDecimal(purchase.Total),
		CreatedAt:        purchase.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        purchase.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPurchaseResponseDTOList converts a list of Purchase models to PurchaseResponseDTOs
func (service *PurchaseService) ToPurchaseResponseDTOList(purchases []*models.Purchase) []dto.PurchaseResponseDTO {
	dtos := make([]dto.PurchaseResponseDTO, len(purchases))
	for i, purchase := range purchases {
		dtos[i] = service.ToPurchaseResponseDTO(purchase)
	}
	return dtos
}
