package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
)

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepo,
	}
}

func (service *CategoryService) CreateCategory(categoryDTO dto.CreateCategoryDTO, userID uint) (*models.Category, error) {
	// Validação básica
	if categoryDTO.Name == "" {
		return nil, errors.New("CreateCategory: nome é obrigatório")
	}

	// Criar categoria
	newCategory := &models.Category{
		Name:   categoryDTO.Name,
		UserID: userID,
	}

	// Salvar no banco
	if err := service.categoryRepository.CreateCategory(newCategory); err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (service *CategoryService) GetCategoryByID(categoryID uint) (*models.Category, error) {
	return service.categoryRepository.GetCategoryByID(categoryID)
}

func (service *CategoryService) GetCategoriesByUserID(userID uint) ([]*models.Category, error) {
	return service.categoryRepository.GetCategoriesByUserID(userID)
}

func (service *CategoryService) UpdateCategory(categoryID uint, categoryDTO dto.UpdateCategoryDTO, userID uint, userRole string) (*models.Category, error) {
	// Buscar categoria existente
	category, err := service.categoryRepository.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("UpdateCategory: categoria não encontrada")
	}

	// Verificar permissão: apenas o próprio usuário ou admin pode atualizar
	if category.UserID != userID && userRole != string(models.RoleAdmin) {
		return nil, errors.New("UpdateCategory: permissão negada: você não pode atualizar categorias de outros usuários")
	}

	// Atualizar dados
	category.Name = categoryDTO.Name

	// Salvar no banco
	if err := service.categoryRepository.UpdateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (service *CategoryService) DeleteCategory(categoryID uint, userID uint, userRole string) error {
	// Buscar categoria existente
	category, err := service.categoryRepository.GetCategoryByID(categoryID)
	if err != nil {
		return errors.New("DeleteCategory: categoria não encontrada")
	}

	// Verificar permissão: apenas o próprio usuário ou admin pode deletar
	if category.UserID != userID && userRole != string(models.RoleAdmin) {
		return errors.New("DeleteCategory: permissão negada: você não pode deletar categorias de outros usuários")
	}

	// Deletar do banco
	return service.categoryRepository.DeleteCategory(categoryID)
}

func (service *CategoryService) GetAllCategories(userRole string) ([]*models.Category, error) {
	// Verificar permissão: apenas admin pode listar todas as categorias
	if userRole != string(models.RoleAdmin) {
		return nil, errors.New("GetAllCategories: permissão negada: apenas administradores podem listar todas as categorias")
	}

	return service.categoryRepository.GetAllCategories()
}

func (service *CategoryService) ToCategoryResponseDTO(category *models.Category) dto.CategoryResponseDTO {
	return dto.CategoryResponseDTO{
		ID:        category.ID,
		Name:      category.Name,
		UserID:    category.UserID,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}
}

func (service *CategoryService) ToCategoryResponseDTOList(categories []*models.Category) []dto.CategoryResponseDTO {
	dtos := make([]dto.CategoryResponseDTO, len(categories))
	for i, category := range categories {
		dtos[i] = service.ToCategoryResponseDTO(category)
	}
	return dtos
}
