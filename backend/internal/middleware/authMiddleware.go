package middleware

import (
	"net/http"
	"strings"

	"github.com/Parron01/AppMercado/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims é a estrutura dos dados contidos no token JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(appConfig *config.Config) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Obter o token do cabeçalho Authorization
		authorizationHeader := ginContext.GetHeader("Authorization")
		if authorizationHeader == "" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de autenticação não fornecido"})
			return
		}

		// O formato esperado é "Bearer {token}"
		tokenParts := strings.Split(authorizationHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}

		tokenString := tokenParts[1]

		// Validar e extrair as claims do token
		claims := &Claims{}
		token, tokenParseError := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(appConfig.JWTSecret), nil
		})

		if tokenParseError != nil || !token.Valid {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de autenticação inválido"})
			return
		}

		// Armazenar as informações do usuário no contexto para uso posterior
		ginContext.Set("userID", claims.UserID)
		ginContext.Set("userEmail", claims.Email)
		ginContext.Set("userRole", claims.Role)

		ginContext.Next()
	}
}
