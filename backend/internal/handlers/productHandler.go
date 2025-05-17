package handlers

import (
	"net/http"
	"strconv"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/middleware"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes configura as rotas de produto
func RegisterProductRoutes(router *gin.Engine, productService *services.ProductService, appConfig *config.Config) {
	authMw := middleware.AuthMiddleware(appConfig)

	productGroup := router.Group("/products")
	{
		// Rota para criar um novo produto (apenas Admin)
		productGroup.POST("/create", authMw, func(c *gin.Context) {
			userRole := c.GetString("userRole")
			if userRole != string(models.RoleAdmin) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Apenas administradores podem criar produtos"})
				return
			}

			var createDTO dto.CreateProductDTO
			if err := c.ShouldBindJSON(&createDTO); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			product, err := productService.CreateProduct(createDTO)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, productService.ToProductResponseDTO(product))
		})

		// Rota para buscar um produto específico (qualquer usuário autenticado)
		productGroup.GET("/:id", authMw, func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produto inválido"})
				return
			}

			product, err := productService.GetProductByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, productService.ToProductResponseDTO(product))
		})

		// Rota para listar todos os produtos (qualquer usuário autenticado)
		productGroup.GET("/all", authMw, func(c *gin.Context) {
			products, err := productService.GetAllProducts()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, productService.ToProductResponseDTOList(products))
		})

		// Rota para atualizar um produto (apenas Admin)
		productGroup.PUT("/update/:id", authMw, func(c *gin.Context) {
			userRole := c.GetString("userRole")
			if userRole != string(models.RoleAdmin) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Apenas administradores podem atualizar produtos"})
				return
			}

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produto inválido"})
				return
			}

			var updateDTO dto.UpdateProductDTO
			if err := c.ShouldBindJSON(&updateDTO); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			product, err := productService.UpdateProduct(uint(id), updateDTO)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Pode ser NotFound também
				return
			}
			c.JSON(http.StatusOK, productService.ToProductResponseDTO(product))
		})

		// Rota para deletar um produto (apenas Admin)
		productGroup.DELETE("/delete/:id", authMw, func(c *gin.Context) {
			userRole := c.GetString("userRole")
			if userRole != string(models.RoleAdmin) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Apenas administradores podem deletar produtos"})
				return
			}

			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produto inválido"})
				return
			}

			if err := productService.DeleteProduct(uint(id)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Pode ser NotFound também
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso"})
		})
	}
}
