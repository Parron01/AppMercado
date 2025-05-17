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

// RegisterPurchaseRoutes configures purchase routes
func RegisterPurchaseRoutes(router *gin.Engine, purchaseService *services.PurchaseService, appConfig *config.Config) {
	authMiddleware := middleware.AuthMiddleware(appConfig)

	purchaseGroup := router.Group("/purchases")
	{
		// Create a new purchase
		purchaseGroup.POST("/create", authMiddleware, func(c *gin.Context) {
			var createDTO dto.CreatePurchaseDTO
			if err := c.ShouldBindJSON(&createDTO); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Get authenticated user ID
			userID := c.GetUint("userID")

			// Create purchase
			purchase, err := purchaseService.CreatePurchase(createDTO, userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTO
			purchaseResponse := purchaseService.ToPurchaseResponseDTO(purchase)

			c.JSON(http.StatusCreated, gin.H{
				"message":  "Purchase created successfully",
				"purchase": purchaseResponse,
			})
		})

		// Get a specific purchase
		purchaseGroup.GET("/:id", authMiddleware, func(c *gin.Context) {
			// Get purchase ID
			purchaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Get purchase
			purchase, err := purchaseService.GetPurchaseByID(uint(purchaseID), userID, userRole)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTO
			purchaseResponse := purchaseService.ToPurchaseResponseDTO(purchase)

			c.JSON(http.StatusOK, gin.H{
				"purchase": purchaseResponse,
			})
		})

		// Get all purchases for the authenticated user
		purchaseGroup.GET("/my", authMiddleware, func(c *gin.Context) {
			// Get authenticated user ID
			userID := c.GetUint("userID")

			// Get purchases
			purchases, err := purchaseService.GetPurchasesByUserID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			purchaseDTOs := purchaseService.ToPurchaseResponseDTOList(purchases)

			c.JSON(http.StatusOK, gin.H{
				"purchases": purchaseDTOs,
				"count":     len(purchaseDTOs),
			})
		})

		// Rota de update removida - compras são imutáveis após criação

		// Delete a purchase
		purchaseGroup.DELETE("/delete/:id", authMiddleware, func(c *gin.Context) {
			// Get purchase ID
			purchaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Delete purchase
			if err := purchaseService.DeletePurchase(uint(purchaseID), userID, userRole); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Purchase deleted successfully",
			})
		})

		// Get all purchases (admin only)
		purchaseGroup.GET("/all", authMiddleware, func(c *gin.Context) {
			// Get authenticated user role
			userRole := c.GetString("userRole")

			// Get all purchases
			purchases, err := purchaseService.GetAllPurchases(userRole)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			purchaseDTOs := purchaseService.ToPurchaseResponseDTOList(purchases)

			c.JSON(http.StatusOK, gin.H{
				"purchases": purchaseDTOs,
				"count":     len(purchaseDTOs),
			})
		})
	}
}
