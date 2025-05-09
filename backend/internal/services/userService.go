package services

import (
	"errors"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService encapsula a lógica de negócio relacionada a usuários
type UserService struct {
	userRepository *repositories.UserRepository
}

// NewUserService cria uma nova instância do serviço de usuários
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

// CreateUser cria um novo usuário com senha criptografada
func (service *UserService) CreateUser(userDTO dto.CreateUserDTO) (*models.User, error) {
	// Validação básica
	if userDTO.Name == "" || userDTO.Email == "" || userDTO.Password == "" {
		return nil, errors.New("nome, email e senha são obrigatórios")
	}

	// Verificar se o email já existe
	existingUser, emailError := service.userRepository.GetUserByEmail(userDTO.Email)
	if emailError == nil && existingUser != nil {
		return nil, errors.New("email já cadastrado")
	}

	// Validação e atribuição do role
	roleToUse := string(models.DefaultRole())
	if userDTO.Role != "" {
		if !models.IsValidRole(userDTO.Role) {
			return nil, errors.New("papel inválido: deve ser Admin, Standard ou Guest")
		}
		roleToUse = userDTO.Role
	}

	// Hash da senha
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if hashError != nil {
		return nil, hashError
	}

	// Criar o usuário
	newUser := &models.User{
		Name:         userDTO.Name,
		Email:        userDTO.Email,
		PasswordHash: string(hashedPassword),
		Role:         roleToUse,
	}

	// Salvar no banco
	if saveError := service.userRepository.CreateUser(newUser); saveError != nil {
		return nil, saveError
	}

	return newUser, nil
}

// GetUserByEmail busca um usuário pelo email
func (service *UserService) GetUserByEmail(email string) (*models.User, error) {
	return service.userRepository.GetUserByEmail(email)
}

// GetUserByID busca um usuário pelo ID
func (service *UserService) GetUserByID(id uint) (*models.User, error) {
	return service.userRepository.GetUserByID(id)
}

// VerifyPassword verifica se uma senha corresponde ao hash armazenado
func (service *UserService) VerifyPassword(user *models.User, password string) bool {
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return passwordError == nil
}

// ToUserResponseDTO converte um User model para UserResponseDTO
func (service *UserService) ToUserResponseDTO(user *models.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
