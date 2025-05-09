package dto

// CreateUserDTO representa os dados necessários para criar um usuário
type CreateUserDTO struct {
	Name     string `json:"name" binding:"required" example:"João Silva"`
	Email    string `json:"email" binding:"required,email" example:"joao@email.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Role     string `json:"role" binding:"omitempty,oneof=Admin Standard Guest" example:"Standard"`
}

// LoginUserDTO representa os dados necessários para autenticar um usuário
type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponseDTO representa os dados de um usuário para resposta HTTP
type UserResponseDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
