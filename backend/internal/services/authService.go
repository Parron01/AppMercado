package services

import (
	"errors"
	"time"

	"github.com/Parron01/AppMercado/backend/internal/dto"
	"github.com/Parron01/AppMercado/backend/internal/models"
	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/golang-jwt/jwt/v4"
)

// AuthService lida com a lógica de negócios de autenticação
type AuthService struct {
	userService *UserService
	appConfig   *config.Config
}

// Estrutura de claims para o JWT (similar a payload no JWT)
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewAuthService cria uma nova instância de AuthService
func NewAuthService(userService *UserService, cfg *config.Config) *AuthService {
	return &AuthService{
		userService: userService,
		appConfig:   cfg,
	}
}

// Register registra um novo usuário e retorna um token JWT
func (authService *AuthService) Register(userDTO dto.CreateUserDTO) (*dto.UserResponseDTO, string, error) {
	// Criar usuário usando userService
	newUser, createError := authService.userService.CreateUser(userDTO)
	if createError != nil {
		return nil, "", createError
	}

	// Gerar JWT
	tokenString, tokenError := authService.GenerateToken(newUser)
	if tokenError != nil {
		return nil, "", tokenError
	}

	// Converter para DTO
	userResponse := authService.userService.ToUserResponseDTO(newUser)

	return &userResponse, tokenString, nil
}

// Login autentica um usuário e retorna um token JWT
func (authService *AuthService) Login(loginDTO dto.LoginUserDTO) (*dto.UserResponseDTO, string, error) {
	// Buscar usuário pelo email usando userService
	user, findError := authService.userService.GetUserByEmail(loginDTO.Email)
	if findError != nil {
		return nil, "", errors.New("credenciais inválidas")
	}

	// Verificar senha usando userService
	if !authService.userService.VerifyPassword(user, loginDTO.Password) {
		return nil, "", errors.New("credenciais inválidas")
	}

	// Gerar JWT
	tokenString, tokenError := authService.GenerateToken(user)
	if tokenError != nil {
		return nil, "", tokenError
	}

	// Converter para DTO
	userResponse := authService.userService.ToUserResponseDTO(user)

	return &userResponse, tokenString, nil
}

// GenerateToken gera um token JWT para um usuário
func (authService *AuthService) GenerateToken(user *models.User) (string, error) {
	// Definir período de expiração
	expirationTime := time.Now().Add(time.Hour * time.Duration(authService.appConfig.JWTExpirationHours))

	// Criar claims (payload do JWT)
	tokenClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Criar token com claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Assinar o token com a chave secreta
	tokenString, signError := jwtToken.SignedString([]byte(authService.appConfig.JWTSecret))
	if signError != nil {
		return "", signError
	}

	return tokenString, nil
}
