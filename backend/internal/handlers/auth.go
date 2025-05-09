package handlers

import (
	"net/http"
	"strings"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// formatValidationError formata erros de validação do gin
func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages,
					"O campo "+e.Field()+" é obrigatório")
			case "email":
				errorMessages = append(errorMessages,
					"O campo "+e.Field()+" deve ser um email válido")
			case "min":
				errorMessages = append(errorMessages,
					"O campo "+e.Field()+" deve ter no mínimo "+e.Param()+" caracteres")
			case "oneof":
				errorMessages = append(errorMessages,
					"O campo "+e.Field()+" deve ser um dos seguintes valores: "+e.Param())
			}
		}
		return strings.Join(errorMessages, "; ")
	}
	return err.Error()
}

// RegisterAuthRoutes configura as rotas de autenticação
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", func(ginContext *gin.Context) {
			var createUserDTO dto.CreateUserDTO
			if bindError := ginContext.ShouldBindJSON(&createUserDTO); bindError != nil {
				ginContext.JSON(http.StatusBadRequest, gin.H{
					"error": formatValidationError(bindError),
					"validRoles": []string{
						string(models.RoleAdmin),
						string(models.RoleStandard),
						string(models.RoleGuest),
					},
				})
				return
			}

			userResponseDTO, tokenString, registerError := authService.Register(createUserDTO)
			if registerError != nil {
				ginContext.JSON(http.StatusBadRequest, gin.H{"error": registerError.Error()})
				return
			}

			ginContext.JSON(http.StatusCreated, gin.H{
				"message": "Usuário criado com sucesso",
				"user":    userResponseDTO,
				"token":   tokenString,
			})
		})

		authGroup.POST("/login", func(ginContext *gin.Context) {
			var loginDTO dto.LoginUserDTO
			if bindError := ginContext.ShouldBindJSON(&loginDTO); bindError != nil {
				ginContext.JSON(http.StatusBadRequest, gin.H{"error": bindError.Error()})
				return
			}

			userResponseDTO, tokenString, loginError := authService.Login(loginDTO)

			if loginError != nil {
				ginContext.JSON(http.StatusUnauthorized, gin.H{"error": loginError.Error()})
				return
			}

			ginContext.JSON(http.StatusOK, gin.H{
				"user":  userResponseDTO,
				"token": tokenString,
			})
		})
	}
}
