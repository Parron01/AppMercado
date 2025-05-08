package handlers

import (
	"net/http"

	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", func(ginContext *gin.Context) {

			ginContext.JSON(http.StatusCreated, gin.H{"message": "usuário criado (endpoint placeholder)"})
		})

		auth.POST("/login", func(ginContext *gin.Context) {

			ginContext.JSON(http.StatusOK, gin.H{"token": "jwt-token (endpoint placeholder)"})
		})
	}
}
