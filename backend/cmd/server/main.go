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

	// 3) Instancia repositórios
	userRepository := repositories.NewUserRepository(database)
	categoryRepository := repositories.NewCategoryRepository(database)
	productRepository := repositories.NewProductRepository(database)
	purchaseRepository := repositories.NewPurchaseRepository(database)

	// 4) Instancia serviços
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userService, appConfig)
	categoryService := services.NewCategoryService(categoryRepository)
	productService := services.NewProductService(productRepository)
	purchaseService := services.NewPurchaseService(purchaseRepository, productService)

	// 5) Cria Gin Engine e registra rotas/handlers
	router := gin.Default()
	handlers.RegisterAuthRoutes(router, authService)
	handlers.RegisterUserRoutes(router, userService, appConfig)
	handlers.RegisterCategoryRoutes(router, categoryService, appConfig)
	handlers.RegisterProductRoutes(router, productService, appConfig)
	handlers.RegisterPurchaseRoutes(router, purchaseService, appConfig)

	// 6) Inicia servidor HTTP na porta configurada
	router.Run(":" + appConfig.ServerPort)
}
