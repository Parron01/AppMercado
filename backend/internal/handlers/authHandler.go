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

// formatValidationError formata erros de validação do gin e retorna mensagem e o tipo do erro
func formatValidationError(err error) (string, string) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		var errorType string

		for _, e := range validationErrors {
			errorType = e.Tag()

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

		return strings.Join(errorMessages, "; "), errorType
	}

	return err.Error(), ""
}

// RegisterAuthRoutes configura as rotas de autenticação
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", func(ginContext *gin.Context) {
			var createUserDTO dto.CreateUserDTO
			if bindError := ginContext.ShouldBindJSON(&createUserDTO); bindError != nil {
				errorMsg, errorType := formatValidationError(bindError)

				response := gin.H{
					"error": errorMsg,
				}

				// Adiciona informações específicas baseadas no tipo de erro
				if errorType == "oneof" {
					response["validRoles"] = []string{
						string(models.RoleAdmin),
						string(models.RoleStandard),
						string(models.RoleGuest),
					}
				} else if errorType == "email" {
					response["emailExample"] = "usuario@exemplo.com"
				} else if errorType == "min" {
					response["passwordInfo"] = "A senha deve ter no mínimo 6 caracteres"
				}

				ginContext.JSON(http.StatusBadRequest, response)
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
				errorMsg, errorType := formatValidationError(bindError)

				response := gin.H{
					"error": errorMsg,
				}

				if errorType == "email" {
					response["emailExample"] = "usuario@exemplo.com"
				}

				ginContext.JSON(http.StatusBadRequest, response)
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
