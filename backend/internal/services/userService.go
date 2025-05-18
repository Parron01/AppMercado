package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService encapsula a lógica de negócio relacionada a usuários
type UserService struct {
	userRepository  *repositories.UserRepository
	categoryService *CategoryService // Adicionado para criar categorias padrão
}

// NewUserService cria uma nova instância do serviço de usuários
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

// SetCategoryService configura o serviço de categoria para permitir criação de categorias padrão
func (service *UserService) SetCategoryService(categoryService *CategoryService) {
	service.categoryService = categoryService
}

// CreateUser cria um novo usuário com senha criptografada
func (service *UserService) CreateUser(userDTO dto.CreateUserDTO) (*models.User, error) {
	// Validação básica
	if userDTO.Name == "" || userDTO.Email == "" || userDTO.Password == "" {
		return nil, errors.New("CreateUser: nome, email e senha são obrigatórios")
	}

	// Verificar se o email já existe
	existingUser, emailError := service.userRepository.GetUserByEmail(userDTO.Email)
	if emailError == nil && existingUser != nil {
		return nil, errors.New("CreateUser: email já cadastrado")
	}

	// Validação e atribuição do role
	roleToUse := string(models.DefaultRole())
	if userDTO.Role != "" {
		if !models.IsValidRole(userDTO.Role) {
			return nil, errors.New("CreateUser: papel inválido: deve ser Admin, Standard ou Guest")
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

	// Criar categorias padrão para o novo usuário
	if service.categoryService != nil {
		go service.createDefaultCategories(newUser.ID) // Executa assincronamente para não bloquear resposta
	}

	return newUser, nil
}

// createDefaultCategories cria categorias padrão para um novo usuário
func (service *UserService) createDefaultCategories(userID uint) {
	// Lista de categorias padrão comuns em supermercados
	defaultCategories := []string{
		"Hortifrúti",
		"Laticínios",
		"Carnes",
		"Padaria",
		"Bebidas",
		"Higiene",
		"Limpeza",
		"Congelados",
		"Mercearia",
		"Cereais",
	}

	for _, categoryName := range defaultCategories {
		service.categoryService.CreateCategory(dto.CreateCategoryDTO{
			Name: categoryName,
		}, userID)
		// Ignora erros para não interromper a criação das outras categorias
		// Se uma falhar, as outras ainda são criadas
	}
}

// DeleteUser remove um usuário pelo ID (apenas o próprio usuário ou admin)
func (service *UserService) DeleteUser(userID uint, requestingUserID uint, requestingUserRole string) error {
	// Se não é o próprio usuário e não é um admin, negue a operação
	if userID != requestingUserID && requestingUserRole != string(models.RoleAdmin) {
		return errors.New("DeleteUser: permissão negada: você não pode deletar outro usuário")
	}

	// Verificar se o usuário existe
	_, err := service.userRepository.GetUserByID(userID)
	if err != nil {
		return errors.New("DeleteUser: usuário não encontrado")
	}

	// Deletar o usuário
	return service.userRepository.DeleteUser(userID)
}

// GetAllUsers retorna todos os usuários (apenas para admin)
func (service *UserService) GetAllUsers(requestingUserRole string) ([]*models.User, error) {
	// Verificar se o usuário é admin
	if requestingUserRole != string(models.RoleAdmin) {
		return nil, errors.New("GetAllUsers: permissão negada: apenas administradores podem listar todos os usuários")
	}

	// Buscar todos os usuários
	return service.userRepository.GetAllUsers()
}

// ToUserResponseDTOList converte uma lista de User model para lista de UserResponseDTO
func (service *UserService) ToUserResponseDTOList(users []*models.User) []dto.UserResponseDTO {
	dtos := make([]dto.UserResponseDTO, len(users))
	for i, user := range users {
		dtos[i] = service.ToUserResponseDTO(user)
	}
	return dtos
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
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
