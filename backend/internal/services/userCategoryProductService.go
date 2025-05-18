package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"gorm.io/gorm"
)

// UserCategoryProductService handles business logic for user-category-product relationships
type UserCategoryProductService struct {
	ucpRepository   *repositories.UserCategoryProductRepository
	categoryService *CategoryService
	productService  *ProductService
}

// NewUserCategoryProductService creates a new instance of UserCategoryProductService
func NewUserCategoryProductService(
	ucpRepo *repositories.UserCategoryProductRepository,
	categoryService *CategoryService,
	productService *ProductService) *UserCategoryProductService {
	return &UserCategoryProductService{
		ucpRepository:   ucpRepo,
		categoryService: categoryService,
		productService:  productService,
	}
}

// CreateUserCategoryProduct creates a new user-category-product relationship
func (service *UserCategoryProductService) CreateUserCategoryProduct(
	createDTO dto.CreateUserCategoryProductDTO,
	userID uint) (*models.UserCategoryProduct, error) {

	// Verify if category exists and belongs to the user
	category, err := service.categoryService.GetCategoryByID(createDTO.CategoryID)
	if err != nil {
		return nil, errors.New("CreateUserCategoryProduct: categoria não encontrada")
	}
	if category.UserID != userID {
		return nil, errors.New("CreateUserCategoryProduct: esta categoria não pertence ao usuário")
	}

	// Verify if product exists
	_, err = service.productService.GetProductByID(createDTO.ProductID)
	if err != nil {
		return nil, errors.New("CreateUserCategoryProduct: produto não encontrado")
	}

	// Check if the relationship already exists
	existingUCP, err := service.ucpRepository.GetUserCategoryProduct(userID, createDTO.CategoryID, createDTO.ProductID)
	if err == nil && existingUCP != nil {
		return nil, errors.New("CreateUserCategoryProduct: este produto já está associado a esta categoria para este usuário")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the relationship
	ucp := &models.UserCategoryProduct{
		UserID:     userID,
		CategoryID: createDTO.CategoryID,
		ProductID:  createDTO.ProductID,
	}

	if err := service.ucpRepository.CreateUserCategoryProduct(ucp); err != nil {
		return nil, err
	}

	// Recarregar o objeto com todos os relacionamentos para ter os nomes disponíveis
	reloadedUCP, err := service.ucpRepository.ReloadUserCategoryProductWithRelations(ucp.ID)
	if err != nil {
		return nil, err
	}

	return reloadedUCP, nil
}

// GetUserCategoryProductByID retrieves a user-category-product relationship by its ID
func (service *UserCategoryProductService) GetUserCategoryProductByID(
	ucpID uint, userID uint, userRole string) (*models.UserCategoryProduct, error) {

	// Get the relationship
	ucp, err := service.ucpRepository.GetUserCategoryProductByID(ucpID)
	if err != nil {
		return nil, errors.New("GetUserCategoryProductByID: relação não encontrada")
	}

	// Check if user has permission to view this relationship
	if ucp.UserID != userID && userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetUserCategoryProductByID: permissão negada: você não pode visualizar categorias de produtos de outros usuários")
	}

	return ucp, nil
}

// GetUserCategoryProductsByUserID retrieves all user-category-product relationships for a specific user
func (service *UserCategoryProductService) GetUserCategoryProductsByUserID(userID uint) ([]*models.UserCategoryProduct, error) {
	return service.ucpRepository.GetUserCategoryProductsByUserID(userID)
}

// GetUserCategoryProductsByCategory retrieves all user-category-product relationships for a specific category
func (service *UserCategoryProductService) GetUserCategoryProductsByCategory(
	categoryID uint, userID uint, userRole string) ([]*models.UserCategoryProduct, error) {

	// Verify if category exists and belongs to the user
	category, err := service.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("GetUserCategoryProductsByCategory: categoria não encontrada")
	}

	// Non-admin users can only view their own category relationships
	if category.UserID != userID && userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetUserCategoryProductsByCategory: permissão negada: você não pode visualizar categorias de outros usuários")
	}

	return service.ucpRepository.GetUserCategoryProductsByCategoryID(categoryID)
}

// DeleteUserCategoryProduct deletes a user-category-product relationship
func (service *UserCategoryProductService) DeleteUserCategoryProduct(
	ucpID uint, userID uint, userRole string) error {

	// Get the relationship
	ucp, err := service.ucpRepository.GetUserCategoryProductByID(ucpID)
	if err != nil {
		return errors.New("DeleteUserCategoryProduct: relação não encontrada")
	}

	// Check if user has permission to delete this relationship
	if ucp.UserID != userID && userRole != string(models.RoleAdmin) {
		return errors.New("DeleteUserCategoryProduct: permissão negada: você não pode excluir categorias de produtos de outros usuários")
	}

	return service.ucpRepository.DeleteUserCategoryProduct(ucpID)
}

// DeleteUserCategoryProductByFields deletes a user-category-product relationship by its composite fields
func (service *UserCategoryProductService) DeleteUserCategoryProductByFields(
	categoryID, productID, requestingUserID uint, userRole string) error {

	// Verify if category exists and belongs to the user
	category, err := service.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		return errors.New("DeleteUserCategoryProductByFields: categoria não encontrada")
	}

	// Non-admin users can only delete their own category relationships
	if category.UserID != requestingUserID && userRole != string(models.RoleAdmin) {
		return errors.New("DeleteUserCategoryProductByFields: permissão negada: você não pode manipular categorias de outros usuários")
	}

	return service.ucpRepository.DeleteUserCategoryProductByFields(requestingUserID, categoryID, productID)
}

// GetAllUserCategoryProducts retrieves all user-category-product relationships (admin only)
func (service *UserCategoryProductService) GetAllUserCategoryProducts(userRole string) ([]*models.UserCategoryProduct, error) {
	// Check if user is admin
	if userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetAllUserCategoryProducts: permissão negada: apenas administradores podem listar todas as relações produto-categoria")
	}

	return service.ucpRepository.GetAllUserCategoryProducts()
}

// ToUserCategoryProductResponseDTO converts a UserCategoryProduct model to UserCategoryProductResponseDTO
func (service *UserCategoryProductService) ToUserCategoryProductResponseDTO(ucp *models.UserCategoryProduct) dto.UserCategoryProductResponseDTO {
	return dto.UserCategoryProductResponseDTO{
		ID:           ucp.ID,
		UserID:       ucp.UserID,
		UserName:     ucp.User.Name,
		CategoryID:   ucp.CategoryID,
		CategoryName: ucp.Category.Name,
		ProductID:    ucp.ProductID,
		ProductName:  ucp.Product.Name,
		CreatedAt:    ucp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    ucp.UpdatedAt.Format(time.RFC3339),
	}
}

// ToUserCategoryProductResponseDTOList converts a list of UserCategoryProduct models to UserCategoryProductResponseDTOs
func (service *UserCategoryProductService) ToUserCategoryProductResponseDTOList(
	ucps []*models.UserCategoryProduct) []dto.UserCategoryProductResponseDTO {

	dtos := make([]dto.UserCategoryProductResponseDTO, len(ucps))
	for i, ucp := range ucps {
		dtos[i] = service.ToUserCategoryProductResponseDTO(ucp)
	}
	return dtos
}
