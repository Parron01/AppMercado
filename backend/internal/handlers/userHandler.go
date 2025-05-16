package handlers

import (
	"net/http"
	"strconv"

	"github.com/Parron01/AppMercado/backend/internal/middleware"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes configura as rotas de usuário
func RegisterUserRoutes(router *gin.Engine, userService *services.UserService, appConfig *config.Config) {
	// Instancia o middleware de autenticação
	authMiddleware := middleware.AuthMiddleware(appConfig)

	userGroup := router.Group("/users")
	{
		// Rota para listar todos os usuários (apenas admin)
		userGroup.GET("/all", authMiddleware, func(context *gin.Context) {
			// Pegando o role do usuário autenticado do contexto
			userRole := context.GetString("userRole")

			users, err := userService.GetAllUsers(userRole)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			// Converter para DTOs
			userDTOs := userService.ToUserResponseDTOList(users)

			context.JSON(http.StatusOK, gin.H{
				"users": userDTOs,
				"count": len(userDTOs),
			})
		})

		// Rota para deletar um usuário
		userGroup.DELETE("/delete/:id", authMiddleware, func(context *gin.Context) {
			// Obtendo ID do usuário a ser deletado
			userID, err := strconv.ParseUint(context.Param("id"), 10, 32)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			// Pegando ID e role do usuário autenticado do contexto
			requestingUserID := context.GetUint("userID")
			requestingUserRole := context.GetString("userRole")

			// Tentando deletar o usuário
			err = userService.DeleteUser(uint(userID), requestingUserID, requestingUserRole)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			context.JSON(http.StatusOK, gin.H{
				"message": "Usuário deletado com sucesso",
			})
		})
	}
}
