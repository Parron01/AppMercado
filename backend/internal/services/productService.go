package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"gorm.io/gorm"
)

// ProductService define a interface para a lógica de negócios de produtos
type ProductService struct {
	productRepo *repositories.ProductRepository
}

// NewProductService cria uma nova instância de ProductService
func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// CreateProduct cria um novo produto
func (s *ProductService) CreateProduct(createDTO dto.CreateProductDTO) (*models.Product, error) {
	var barcodeToSave *string

	if createDTO.Barcode != "" {
		// Verificar se já existe um produto com o mesmo código de barras (se não for vazio)
		_, err := s.productRepo.GetProductByBarcode(createDTO.Barcode)
		if err == nil { // Se err for nil, significa que um produto foi encontrado
			return nil, errors.New("produto com este código de barras já existe")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) { // Se o erro não for 'registro não encontrado', é outro erro
			return nil, err
		}
		// Se o barcode não é vazio, usamos o valor fornecido
		tempBarcode := createDTO.Barcode
		barcodeToSave = &tempBarcode
	} // Se createDTO.Barcode for "", barcodeToSave permanece nil, resultando em NULL no banco

	product := &models.Product{
		Name:         createDTO.Name,
		AveragePrice: 0, // Inicializa com zero - será atualizado com base nas compras futuras
		Barcode:      barcodeToSave,
	}

	if err := s.productRepo.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// GetProductByID busca um produto pelo ID
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}
	return product, nil
}

// GetAllProducts retorna todos os produtos
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAllProducts()
}

// UpdateProduct atualiza um produto existente
func (s *ProductService) UpdateProduct(id uint, updateDTO dto.UpdateProductDTO) (*models.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado para atualização")
		}
		return nil, err
	}

	if updateDTO.Name != nil {
		product.Name = *updateDTO.Name
	}
	if updateDTO.AveragePrice != nil {
		product.AveragePrice = *updateDTO.AveragePrice
	}

	// Lógica para atualizar o barcode
	if updateDTO.Barcode != nil { // Se o campo barcode foi fornecido na atualização (não é nil o ponteiro do DTO)
		newBarcodeValueFromDTO := *updateDTO.Barcode // Desreferencia o valor do DTO

		if newBarcodeValueFromDTO == "" { // Cliente quer definir o barcode como nulo/vazio
			product.Barcode = nil
		} else { // Cliente quer definir um barcode não vazio
			// Verificar se o novo barcode é diferente do atual (ou se o atual era nil)
			// e se é único, antes de atribuir.
			isCurrentBarcodeNil := product.Barcode == nil
			currentBarcodeValue := ""
			if !isCurrentBarcodeNil {
				currentBarcodeValue = *product.Barcode
			}

			// Só verificar unicidade e atualizar se o novo valor é realmente diferente
			// ou se o barcode atual era nil e o novo não é vazio.
			if newBarcodeValueFromDTO != currentBarcodeValue || (isCurrentBarcodeNil && newBarcodeValueFromDTO != "") {
				existingProduct, err := s.productRepo.GetProductByBarcode(newBarcodeValueFromDTO)
				if err == nil && existingProduct.ID != product.ID { // Encontrou outro produto com o mesmo barcode
					return nil, errors.New("novo código de barras já está em uso por outro produto")
				}
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Erro inesperado ao buscar barcode
					return nil, err
				}
			}
			product.Barcode = &newBarcodeValueFromDTO // Atribui o novo barcode não vazio
		}
	} // Se updateDTO.Barcode for nil (campo omitido no JSON), não fazemos nada com o barcode do produto

	if err := s.productRepo.UpdateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// DeleteProduct remove um produto
func (s *ProductService) DeleteProduct(id uint) error {
	_, err := s.productRepo.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produto não encontrado para exclusão")
		}
		return err
	}
	return s.productRepo.DeleteProduct(id)
}

// ToProductResponseDTO converte um modelo Product para ProductResponseDTO
func (s *ProductService) ToProductResponseDTO(product *models.Product) dto.ProductResponseDTO {
	return dto.ProductResponseDTO{
		ID:           product.ID,
		Name:         product.Name,
		AveragePrice: product.AveragePrice,
		Barcode:      product.Barcode, // product.Barcode é *string, dto.ProductResponseDTO.Barcode é *string
		CreatedAt:    product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    product.UpdatedAt.Format(time.RFC3339),
	}
}

// ToProductResponseDTOList converte uma lista de modelos Product para uma lista de ProductResponseDTO
func (s *ProductService) ToProductResponseDTOList(products []models.Product) []dto.ProductResponseDTO {
	responseDTOs := make([]dto.ProductResponseDTO, len(products))
	for i, p := range products {
		responseDTOs[i] = s.ToProductResponseDTO(&p)
	}
	return responseDTOs
}
