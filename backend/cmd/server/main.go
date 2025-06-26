package main

import (
	"github.com/Parron01/AppMercado/backend/internal/handlers"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"github.com/Parron01/AppMercado/backend/internal/services"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-contrib/cors"
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
	priceHistoryRepository := repositories.NewPriceHistoryRepository(database)
	userCategoryProductRepository := repositories.NewUserCategoryProductRepository(database)

	// 4) Instancia serviços
	userService := services.NewUserService(userRepository)
	categoryService := services.NewCategoryService(categoryRepository)

	// Configura dependência circular entre UserService e CategoryService
	userService.SetCategoryService(categoryService)

	authService := services.NewAuthService(userService, appConfig)
	productService := services.NewProductService(productRepository)
	purchaseService := services.NewPurchaseService(purchaseRepository, productService)
	priceHistoryService := services.NewPriceHistoryService(priceHistoryRepository, productService, userService)
	userCategoryProductService := services.NewUserCategoryProductService(userCategoryProductRepository, categoryService, productService)

	// 5) Resolve circular dependencies
	purchaseService.SetPriceHistoryService(priceHistoryService)
	productService.SetPriceHistoryService(priceHistoryService)

	// 6) Cria Gin Engine e registra rotas/handlers
	router := gin.Default()

	// Configura CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Substitua pelo domínio do frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handlers.RegisterAuthRoutes(router, authService)
	handlers.RegisterUserRoutes(router, userService, appConfig)
	handlers.RegisterCategoryRoutes(router, categoryService, appConfig)
	handlers.RegisterProductRoutes(router, productService, appConfig)
	handlers.RegisterPurchaseRoutes(router, purchaseService, appConfig)
	handlers.RegisterPriceHistoryRoutes(router, priceHistoryService, appConfig)
	handlers.RegisterUserCategoryProductRoutes(router, userCategoryProductService, appConfig)

	// 7) Inicia servidor HTTP na porta configurada
	router.Run(":" + appConfig.ServerPort)
}
