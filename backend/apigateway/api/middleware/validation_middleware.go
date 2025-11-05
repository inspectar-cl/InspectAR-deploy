package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidationMiddleware valida el acceso del usuario a edificios o activos específicos
func ValidationMiddleware(validationType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Si no hay validación requerida, continuar
		if validationType == "" {
			c.Next()
			return
		}

		// Si el usuario es Root, omitir validación (tiene acceso a todo)
		scopeInterface, exists := c.Get("scope")
		if exists {
			if scope, ok := scopeInterface.(string); ok {
				if strings.Contains(scope, "user-type:Root") {
					fmt.Println("Usuario Root detectado, omitiendo validación de acceso")
					c.Next()
					return
				}
			}
		}

		gestionURL := os.Getenv("GESTION_URL")
		if gestionURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GESTION_URL no configurado"})
			c.Abort()
			return
		}

		// Extraer email del contexto (guardado por AuthMiddleware)
		emailInterface, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email no encontrado en token"})
			c.Abort()
			return
		}

		email, ok := emailInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email inválido"})
			c.Abort()
			return
		}

		// Determinar el tipo de validación y extraer el ID correspondiente
		var resourceID string
		var validationURL string

		switch strings.ToLower(validationType) {
		case "edificio":
			// Buscar :id_edificio en los parámetros de la ruta
			resourceID = c.Param("id_edificio")
			if resourceID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio no encontrado en la ruta"})
				c.Abort()
				return
			}
			validationURL = fmt.Sprintf("%s/usuarios/%s/edificio/%s/acceso", gestionURL, email, resourceID)

		case "activo":
			// Buscar :id_activo en los parámetros de la ruta
			resourceID = c.Param("id_activo")
			if resourceID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo no encontrado en la ruta"})
				c.Abort()
				return
			}
			validationURL = fmt.Sprintf("%s/usuarios/%s/activo/%s/acceso", gestionURL, email, resourceID)

		default:
			// Tipo de validación no soportado
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipo de validación no soportado"})
			c.Abort()
			return
		}

		// Hacer la petición GET al microservicio de gestión
		resp, err := http.Get(validationURL)
		if err != nil {
			fmt.Printf("Error validando acceso: %v\n", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio de validación no disponible"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Decodificar la respuesta
		var validationResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
			fmt.Printf("Error decodificando respuesta de validación: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando validación"})
			c.Abort()
			return
		}

		// Verificar el campo "tiene_acceso"
		tieneAcceso, exists := validationResponse["tiene_acceso"]
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Respuesta de validación inválida"})
			c.Abort()
			return
		}

		// Convertir a bool
		tieneAccesoBool, ok := tieneAcceso.(bool)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Formato de respuesta inválido"})
			c.Abort()
			return
		}

		// Si no tiene acceso, denegar
		if !tieneAccesoBool {
			mensaje := "Acceso denegado al recurso"
			if msg, exists := validationResponse["mensaje"]; exists {
				if msgStr, ok := msg.(string); ok {
					mensaje = msgStr
				}
			}

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Acceso denegado",
				"mensaje": mensaje,
			})
			c.Abort()
			return
		}

		// Si tiene acceso, continuar con el siguiente middleware/handler
		fmt.Printf("Usuario %s tiene acceso al %s %s\n", email, validationType, resourceID)
		c.Next()
	}
}
