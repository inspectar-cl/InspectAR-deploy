package middleware

import (
    "os"
	"log"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

var GestionURL string

func init() {
	GestionURL = os.Getenv("GESTION_URL")
}

// ScopeMiddleware verifica que el usuario tenga al menos uno de los scopes requeridos
func ScopeMiddleware(requiredScopes []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obtener los claims del token
        scopes, exists := c.Get("scope")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Scope no encontrado"})
            c.Abort()
            return
        }

		// Imprimir el scope para depuración
		log.Printf("Scopes del usuario: %v", scopes)

        // Convertir a string y validar si contiene el scope requerido
        scopeStr := scopes.(string)
        scopeList := strings.Split(scopeStr, " ")

        // Verificar si el usuario tiene al menos uno de los scopes requeridos
        for _, userScope := range scopeList {
            for _, requiredScope := range requiredScopes {
                // Se busca el scope requerido en la lista de scopes del token
                if userScope == requiredScope {
                    c.Next() // Permitir acceso
                    return
                }
            }
        }

        // Si no tiene ningún scope necesario, rechazar la solicitud
        c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para esta acción"})
        c.Abort()
    }
}