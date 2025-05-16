package handlers

import (
	"net/http"
	"strconv"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/middleware"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.Engine, categoryService *services.CategoryService, appConfig *config.Config) {
	// Instancia o middleware de autenticação
	authMiddleware := middleware.AuthMiddleware(appConfig)

	categoryGroup := router.Group("/categories")
	{
		// Rota para criar uma nova categoria (autenticado)
		categoryGroup.POST("/create", authMiddleware, func(context *gin.Context) {
			var createCategoryDTO dto.CreateCategoryDTO
			if err := context.ShouldBindJSON(&createCategoryDTO); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Pegando ID do usuário autenticado do contexto
			userID := context.GetUint("userID")

			// Criando a categoria
			newCategory, err := categoryService.CreateCategory(createCategoryDTO, userID)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Converter para DTO
			categoryResponse := categoryService.ToCategoryResponseDTO(newCategory)

			context.JSON(http.StatusCreated, gin.H{
				"message":  "Categoria criada com sucesso",
				"category": categoryResponse,
			})
		})

		// Rota para listar categorias do usuário autenticado
		categoryGroup.GET("/my", authMiddleware, func(context *gin.Context) {
			// Pegando ID do usuário autenticado do contexto
			userID := context.GetUint("userID")

			// Buscando categorias do usuário
			categories, err := categoryService.GetCategoriesByUserID(userID)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Converter para DTOs
			categoryDTOs := categoryService.ToCategoryResponseDTOList(categories)

			context.JSON(http.StatusOK, gin.H{
				"categories": categoryDTOs,
				"count":      len(categoryDTOs),
			})
		})

		// Rota para buscar uma categoria específica
		categoryGroup.GET("/:id", authMiddleware, func(context *gin.Context) {
			// Obtendo ID da categoria
			categoryID, err := strconv.ParseUint(context.Param("id"), 10, 32)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			// Buscando a categoria
			category, err := categoryService.GetCategoryByID(uint(categoryID))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
				return
			}

			// Pegando ID e role do usuário autenticado do contexto
			userID := context.GetUint("userID")
			userRole := context.GetString("userRole")

			// Verificar permissão: apenas o próprio usuário ou admin pode ver detalhes
			if category.UserID != userID && userRole != "Admin" {
				context.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para ver esta categoria"})
				return
			}

			// Converter para DTO
			categoryResponse := categoryService.ToCategoryResponseDTO(category)

			context.JSON(http.StatusOK, gin.H{
				"category": categoryResponse,
			})
		})

		// Rota para atualizar uma categoria
		categoryGroup.PUT("/update/:id", authMiddleware, func(context *gin.Context) {
			// Obtendo ID da categoria
			categoryID, err := strconv.ParseUint(context.Param("id"), 10, 32)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			var updateCategoryDTO dto.UpdateCategoryDTO
			if err := context.ShouldBindJSON(&updateCategoryDTO); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Pegando ID e role do usuário autenticado do contexto
			userID := context.GetUint("userID")
			userRole := context.GetString("userRole")

			// Atualizando a categoria
			updatedCategory, err := categoryService.UpdateCategory(uint(categoryID), updateCategoryDTO, userID, userRole)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Converter para DTO
			categoryResponse := categoryService.ToCategoryResponseDTO(updatedCategory)

			context.JSON(http.StatusOK, gin.H{
				"message":  "Categoria atualizada com sucesso",
				"category": categoryResponse,
			})
		})

		// Rota para deletar uma categoria
		categoryGroup.DELETE("/delete/:id", authMiddleware, func(context *gin.Context) {
			// Obtendo ID da categoria
			categoryID, err := strconv.ParseUint(context.Param("id"), 10, 32)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			// Pegando ID e role do usuário autenticado do contexto
			userID := context.GetUint("userID")
			userRole := context.GetString("userRole")

			// Deletando a categoria
			if err := categoryService.DeleteCategory(uint(categoryID), userID, userRole); err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			context.JSON(http.StatusOK, gin.H{
				"message": "Categoria deletada com sucesso",
			})
		})

		// Rota para listar todas as categorias (apenas admin)
		categoryGroup.GET("/all", authMiddleware, func(context *gin.Context) {
			// Pegando o role do usuário autenticado do contexto
			userRole := context.GetString("userRole")

			categories, err := categoryService.GetAllCategories(userRole)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Converter para DTOs
			categoryDTOs := categoryService.ToCategoryResponseDTOList(categories)

			context.JSON(http.StatusOK, gin.H{
				"categories": categoryDTOs,
				"count":      len(categoryDTOs),
			})
		})
	}
}
