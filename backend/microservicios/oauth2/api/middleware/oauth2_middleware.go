package middleware

import (
	"net/http"
	"oauth2/internal/services"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(tokenString) // Usar la nueva función
		log.Printf("Claims: %v", err)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Set("exp", claims["exp"])
		c.Set("scope", claims["scope"])
		c.Set("device", claims["device"])
		c.Set("email", claims["email"])
		c.Set("empresa", claims["empresa"])
		c.Set("usermane", claims["usermane"])
		c.Next()
	}
}
