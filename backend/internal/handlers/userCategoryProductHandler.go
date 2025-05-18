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

// RegisterUserCategoryProductRoutes configures user-category-product routes
func RegisterUserCategoryProductRoutes(
	router *gin.Engine,
	ucpService *services.UserCategoryProductService,
	appConfig *config.Config) {

	authMiddleware := middleware.AuthMiddleware(appConfig)

	ucpGroup := router.Group("/user-category-products")
	{
		// Create a new user-category-product relationship
		ucpGroup.POST("/create", authMiddleware, func(c *gin.Context) {
			var createDTO dto.CreateUserCategoryProductDTO
			if err := c.ShouldBindJSON(&createDTO); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Get authenticated user ID
			userID := c.GetUint("userID")

			// Create relationship
			ucp, err := ucpService.CreateUserCategoryProduct(createDTO, userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTO
			ucpResponse := ucpService.ToUserCategoryProductResponseDTO(ucp)

			c.JSON(http.StatusCreated, gin.H{
				"message":             "Produto categorizado com sucesso",
				"userCategoryProduct": ucpResponse,
			})
		})

		// Get a specific user-category-product relationship
		ucpGroup.GET("/:id", authMiddleware, func(c *gin.Context) {
			// Get relationship ID
			ucpID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Get relationship
			ucp, err := ucpService.GetUserCategoryProductByID(uint(ucpID), userID, userRole)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTO
			ucpResponse := ucpService.ToUserCategoryProductResponseDTO(ucp)

			c.JSON(http.StatusOK, gin.H{
				"userCategoryProduct": ucpResponse,
			})
		})

		// Get all user-category-product relationships for the authenticated user
		ucpGroup.GET("/my", authMiddleware, func(c *gin.Context) {
			// Get authenticated user ID
			userID := c.GetUint("userID")

			// Get relationships
			ucps, err := ucpService.GetUserCategoryProductsByUserID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			ucpDTOs := ucpService.ToUserCategoryProductResponseDTOList(ucps)

			c.JSON(http.StatusOK, gin.H{
				"userCategoryProducts": ucpDTOs,
				"count":                len(ucpDTOs),
			})
		})

		// Get all user-category-product relationships for a specific category
		ucpGroup.GET("/category/:id", authMiddleware, func(c *gin.Context) {
			// Get category ID
			categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoria inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Get relationships
			ucps, err := ucpService.GetUserCategoryProductsByCategory(uint(categoryID), userID, userRole)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			ucpDTOs := ucpService.ToUserCategoryProductResponseDTOList(ucps)

			c.JSON(http.StatusOK, gin.H{
				"userCategoryProducts": ucpDTOs,
				"count":                len(ucpDTOs),
			})
		})

		// Delete a user-category-product relationship by ID
		ucpGroup.DELETE("/delete/:id", authMiddleware, func(c *gin.Context) {
			// Get relationship ID
			ucpID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Delete relationship
			if err := ucpService.DeleteUserCategoryProduct(uint(ucpID), userID, userRole); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Categorização de produto removida com sucesso",
			})
		})

		// Delete a user-category-product relationship by composite fields
		ucpGroup.DELETE("/delete", authMiddleware, func(c *gin.Context) {
			// Get query parameters
			categoryIDStr := c.Query("categoryId")
			productIDStr := c.Query("productId")

			if categoryIDStr == "" || productIDStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "categoryId e productId são obrigatórios como parâmetros de query"})
				return
			}

			categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoria inválido"})
				return
			}

			productID, err := strconv.ParseUint(productIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produto inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Delete relationship
			if err := ucpService.DeleteUserCategoryProductByFields(uint(categoryID), uint(productID), userID, userRole); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Categorização de produto removida com sucesso",
			})
		})

		// Get all user-category-product relationships (admin only)
		ucpGroup.GET("/all", authMiddleware, func(c *gin.Context) {
			// Get authenticated user role
			userRole := c.GetString("userRole")

			// Get all relationships
			ucps, err := ucpService.GetAllUserCategoryProducts(userRole)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			ucpDTOs := ucpService.ToUserCategoryProductResponseDTOList(ucps)

			c.JSON(http.StatusOK, gin.H{
				"userCategoryProducts": ucpDTOs,
				"count":                len(ucpDTOs),
			})
		})
	}
}
