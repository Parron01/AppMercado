package repositories

import (
	"errors"

	"github.com/Parron01/AppMercado/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository gerencia o acesso aos dados do usuário
type UserRepository struct {
	database *gorm.DB
}

// NewUserRepository cria uma nova instância de UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{database: db}
}

// CreateUser cria um novo usuário no banco de dados
func (repository *UserRepository) CreateUser(user *models.User) error {
	return repository.database.Create(user).Error
}

// GetUserByEmail busca um usuário pelo email
func (repository *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if databaseError := repository.database.Where("email = ?", email).First(&user).Error; databaseError != nil {
		if errors.Is(databaseError, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, databaseError
	}
	return &user, nil
}

func nada() {

}

// GetUserByID busca um usuário pelo ID
func (repository *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if databaseError := repository.database.First(&user, id).Error; databaseError != nil {
		return nil, databaseError
	}
	return &user, nil
}

// UpdateUser atualiza os dados de um usuário
func (repository *UserRepository) UpdateUser(user *models.User) error {
	return repository.database.Save(user).Error
}

// DeleteUser realiza soft delete do usuário (usando gorm.Model)
func (repository *UserRepository) DeleteUser(id uint) error {
	return repository.database.Delete(&models.User{}, id).Error
}
