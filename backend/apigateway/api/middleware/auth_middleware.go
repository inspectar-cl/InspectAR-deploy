package middleware

import (
	"os"
	"fmt"
	"net/http"
	"strings"
	"errors"
	"log"

	"apigateway/internal/models"
	"apigateway/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
)

var oauth2URL string

func init() {
	oauth2URL = os.Getenv("OAUTH2_URL")
}

type AuthService struct {
	Repo      *repository.UserRepository
	TokenRepo *repository.TokenRepository
	Config    *models.Config
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Extraer token del header
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Verificar formato del token Bearer
		tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
            c.Abort()
            return
        }

		// Validar token con OAuth2 service
		claims, err := ValidateToken(tokenParts[1])
		if err != nil {
			log.Printf("Error validating token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// c.Set("token", token)
		c.Set("exp", claims["exp"])
		c.Set("scope", claims["scope"])
		c.Set("device", claims["device"])
		c.Set("email", claims["email"])
		c.Set("empresa", claims["empresa"])
		c.Set("username", claims["username"])
		c.Next()

	}
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Cargar la clave pública desde el archivo
	publicKeyData, err := os.ReadFile("config/jwt_access_public.pem")
	if err != nil {
		return nil, fmt.Errorf("error leyendo clave pública: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave pública: %v", err)
	}

	// Parsear y validar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Asegurar que el método de firma sea RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extraer claims del token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// Método del AuthService
func (s *AuthService) ValidateTokenWithService(tokenString string) (jwt.MapClaims, error) {
    return ValidateToken(tokenString)
}