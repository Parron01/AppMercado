package models

// Role define os papéis disponíveis no sistema
type Role string

// Constantes para os valores válidos de Role
const (
	RoleAdmin    Role = "Admin"
	RoleStandard Role = "Standard"
	RoleGuest    Role = "Guest"
)

// IsValidRole verifica se o papel informado é válido
func IsValidRole(role string) bool {
	return role == string(RoleAdmin) ||
		role == string(RoleStandard) ||
		role == string(RoleGuest)
}

// DefaultRole retorna o papel padrão para novos usuários
func DefaultRole() Role {
	return RoleStandard
}
