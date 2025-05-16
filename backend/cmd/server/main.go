package main

import (
	"github.com/Parron01/AppMercado/backend/internal/handlers"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Carrega configurações (porta, banco, JWT, etc.)
	appConfig := config.Load()

	// 2) Inicializa conexão com DB (PostgreSQL via GORM)
	database := repositories.NewPostgresConn(appConfig)

	// 3) Instancia repositórios e serviços
	userRepository := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userService, appConfig)

	// 4) Cria Gin Engine e registra rotas/handlers
	router := gin.Default()
	handlers.RegisterAuthRoutes(router, authService)
	handlers.RegisterUserRoutes(router, userService, appConfig)

	// 5) Inicia servidor HTTP na porta configurada
	router.Run(":" + appConfig.ServerPort)
}
