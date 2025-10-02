package handlers

import (
	"fmt"
	"gestion/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	repo *repository.UsuarioRepository
}

func NewUsuarioHandler(repo *repository.UsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

// GetEdificiosByEmail obtiene todos los edificios asociados a un usuario por su email
// @Summary Obtener edificios de un usuario por email
// @Description Retorna una lista con todos los datos de los edificios asociados a un usuario basándose en su email
// @Tags usuarios
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Success 200 {object} map[string]interface{} "Lista de edificios del usuario"
// @Failure 400 {object} map[string]interface{} "Email inválido"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /usuarios/edificios/{email} [get]
func (h *UsuarioHandler) GetEdificiosByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email es requerido"})
		return
	}

	// Verificar que el usuario existe
	usuario, err := h.repo.GetUsuarioByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Usuario no encontrado",
			"details": err.Error(),
		})
		return
	}

	// Obtener edificios asociados
	edificios, err := h.repo.GetEdificiosByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener edificios del usuario",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario":   usuario,
		"edificios": edificios,
		"total":     len(edificios),
	})
}

// VerificarAccesoEdificio verifica si un usuario tiene acceso a un edificio específico
// @Summary Verificar acceso de usuario a edificio
// @Description Valida si un usuario tiene relación con un edificio específico. Retorna 200 si tiene acceso, 403 si no tiene acceso.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Param edificio_id path int true "ID del edificio"
// @Success 200 {object} map[string]interface{} "Usuario tiene acceso al edificio"
// @Failure 403 {object} map[string]interface{} "Usuario no tiene acceso al edificio"
// @Failure 400 {object} map[string]interface{} "Parámetros inválidos"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /usuarios/{email}/edificio/{edificio_id}/acceso [get]
func (h *UsuarioHandler) VerificarAccesoEdificio(c *gin.Context) {
	email := c.Param("email")
	edificioIDStr := c.Param("edificio_id")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email es requerido"})
		return
	}

	edificioID := 0
	if _, err := fmt.Sscanf(edificioIDStr, "%d", &edificioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	tieneAcceso, err := h.repo.VerificarAccesoEdificio(email, edificioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar acceso",
			"details": err.Error(),
		})
		return
	}

	if !tieneAcceso {
		c.JSON(http.StatusForbidden, gin.H{
			"error":        "Acceso denegado",
			"mensaje":      "El usuario no tiene acceso a este edificio",
			"email":        email,
			"edificio_id":  edificioID,
			"tiene_acceso": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":      "Usuario tiene acceso al edificio",
		"email":        email,
		"edificio_id":  edificioID,
		"tiene_acceso": true,
	})
}

// VerificarAccesoActivo verifica si un usuario tiene acceso a un activo específico
// @Summary Verificar acceso de usuario a activo
// @Description Valida si un usuario tiene relación con un activo específico a través del edificio. Retorna 200 si tiene acceso, 403 si no tiene acceso.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Param activo_id path int true "ID del activo"
// @Success 200 {object} map[string]interface{} "Usuario tiene acceso al activo"
// @Failure 403 {object} map[string]interface{} "Usuario no tiene acceso al activo"
// @Failure 400 {object} map[string]interface{} "Parámetros inválidos"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /usuarios/{email}/activo/{activo_id}/acceso [get]
func (h *UsuarioHandler) VerificarAccesoActivo(c *gin.Context) {
	email := c.Param("email")
	activoIDStr := c.Param("activo_id")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email es requerido"})
		return
	}

	activoID := 0
	if _, err := fmt.Sscanf(activoIDStr, "%d", &activoID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	tieneAcceso, err := h.repo.VerificarAccesoActivo(email, activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar acceso",
			"details": err.Error(),
		})
		return
	}

	if !tieneAcceso {
		c.JSON(http.StatusForbidden, gin.H{
			"error":        "Acceso denegado",
			"mensaje":      "El usuario no tiene acceso a este activo",
			"email":        email,
			"activo_id":    activoID,
			"tiene_acceso": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":      "Usuario tiene acceso al activo",
		"email":        email,
		"activo_id":    activoID,
		"tiene_acceso": true,
	})
}
