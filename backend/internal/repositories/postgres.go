package repositories

import (
	"fmt"

	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConn(config *config.Config) *gorm.DB {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBName, config.DBPassword)

	database, error := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if error != nil {
		panic("failed to connect database: " + error.Error())
	}

	database.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{})
	return database
}
