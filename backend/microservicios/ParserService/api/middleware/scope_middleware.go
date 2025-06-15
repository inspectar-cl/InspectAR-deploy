package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware para validar scope en JWT
func ScopeMiddleware(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener los claims del token
		scopes, exists := c.Get("scope")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Scope no encontrado"})
			c.Abort()
			return
		}

		// Convertir a string y validar si contiene el scope requerido
		scopeStr := scopes.(string)
		scopeList := strings.Split(scopeStr, " ")

		for _, scope := range scopeList {
			// Se busca el scope requerido en la lista de scopes del token
			if scope == requiredScope {
				c.Next() // Permitir acceso
				return
			}
		}

		// Si no tiene el scope necesario, rechazar la solicitud
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No tienes permiso para esta acción"})
		c.Abort()
	}
}
