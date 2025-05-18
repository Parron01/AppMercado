package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/middleware"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

// RegisterPriceHistoryRoutes configures price history routes
func RegisterPriceHistoryRoutes(router *gin.Engine, priceHistoryService *services.PriceHistoryService, appConfig *config.Config) {
	authMw := middleware.AuthMiddleware(appConfig)

	priceHistoryGroup := router.Group("/price-history")
	{
		// Get a specific price history entry
		priceHistoryGroup.GET("/:id", authMw, func(c *gin.Context) {
			// Get price history ID
			priceHistoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de histórico de preço inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Get price history
			priceHistory, err := priceHistoryService.GetPriceHistoryByID(uint(priceHistoryID), userID, userRole)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTO
			priceHistoryResponse := priceHistoryService.ToPriceHistoryResponseDTO(priceHistory)

			c.JSON(http.StatusOK, gin.H{
				"priceHistory": priceHistoryResponse,
			})
		})

		// Get price history for a specific product
		priceHistoryGroup.GET("/product/:id", authMw, func(c *gin.Context) {
			// Get product ID
			productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produto inválido"})
				return
			}

			// Parse optional date range parameters
			startDateStr := c.Query("startDate")
			endDateStr := c.Query("endDate")

			var priceHistories []*models.PriceHistory
			var err2 error

			if startDateStr != "" && endDateStr != "" {
				// Parse dates
				startDate, err := time.Parse(time.RFC3339, startDateStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inicial inválido. Use RFC3339."})
					return
				}

				endDate, err := time.Parse(time.RFC3339, endDateStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data final inválido. Use RFC3339."})
					return
				}

				// Get price history for product in date range
				priceHistories, err2 = priceHistoryService.GetPriceHistoryByProductAndDateRange(
					uint(productID), startDate, endDate)
			} else {
				// Get all price history for product
				priceHistories, err2 = priceHistoryService.GetPriceHistoryByProductID(uint(productID))
			}

			if err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
				return
			}

			// Convert to DTOs
			priceHistoryDTOs := priceHistoryService.ToPriceHistoryResponseDTOList(priceHistories)

			c.JSON(http.StatusOK, gin.H{
				"priceHistories": priceHistoryDTOs,
				"count":          len(priceHistoryDTOs),
			})
		})

		// Remover endpoint de estatísticas, agora está no ProductHandler

		// Delete a price history entry
		priceHistoryGroup.DELETE("/delete/:id", authMw, func(c *gin.Context) {
			// Get price history ID
			priceHistoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de histórico de preço inválido"})
				return
			}

			// Get authenticated user information
			userID := c.GetUint("userID")
			userRole := c.GetString("userRole")

			// Delete price history
			if err := priceHistoryService.DeletePriceHistory(uint(priceHistoryID), userID, userRole); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Histórico de preço removido com sucesso",
			})
		})

		// Get all price history entries (admin only)
		priceHistoryGroup.GET("/all", authMw, func(c *gin.Context) {
			// Get authenticated user role
			userRole := c.GetString("userRole")

			// Get all price history entries
			priceHistories, err := priceHistoryService.GetAllPriceHistory(userRole)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Convert to DTOs
			priceHistoryDTOs := priceHistoryService.ToPriceHistoryResponseDTOList(priceHistories)

			c.JSON(http.StatusOK, gin.H{
				"priceHistories": priceHistoryDTOs,
				"count":          len(priceHistoryDTOs),
			})
		})
	}
}
