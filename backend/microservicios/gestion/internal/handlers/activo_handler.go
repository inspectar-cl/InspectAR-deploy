package handlers

import (
	"net/http"
	"strconv"

	"gestion/internal/services"

	"github.com/gin-gonic/gin"
)

type ActivoHandler struct {
	activoService *services.ActivoService
}

func NewActivoHandler(activoService *services.ActivoService) *ActivoHandler {
	return &ActivoHandler{
		activoService: activoService,
	}
}

// GetActivosByEdificio obtiene todos los activos de un edificio específico
// @Summary Obtener activos por edificio
// @Description Obtiene todos los activos que pertenecen a un edificio específico
// @Tags activos
// @Accept json
// @Produce json
// @Param edificio_id path int true "ID del edificio"
// @Success 200 {array} models.Activo
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /activos/edificio/{edificio_id} [get]
func (h *ActivoHandler) GetActivosByEdificio(c *gin.Context) {
	edificioIDStr := c.Param("edificio_id")
	edificioID, err := strconv.Atoi(edificioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de edificio inválido",
		})
		return
	}

	activos, err := h.activoService.GetActivosByEdificio(edificioID)
	if err != nil {
		if err.Error() == "edificio no encontrado" || err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Edificio no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener activos del edificio",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edificio_id": edificioID,
		"activos":     activos,
		"total":       len(activos),
	})
}

// GetActivosByTipo obtiene todos los activos de un tipo específico
// @Summary Obtener activos por tipo
// @Description Obtiene todos los activos que son de un tipo específico (caldera, bomba de agua, ascensor, transformador)
// @Tags activos
// @Accept json
// @Produce json
// @Param tipo path string true "Tipo de activo (caldera, bomba de agua, ascensor, transformador)"
// @Success 200 {array} models.Activo
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /activos/tipo/{tipo} [get]
func (h *ActivoHandler) GetActivosByTipo(c *gin.Context) {
	tipo := c.Param("tipo")
	if tipo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de activo requerido",
		})
		return
	}

	activos, err := h.activoService.GetActivosByTipo(tipo)
	if err != nil {
		if err.Error() != "" && err.Error()[:25] == "tipo de activo inválido:" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "Tipo de activo inválido",
				"tipo_recibido": tipo,
				"tipos_validos": []string{"caldera", "bomba de agua", "ascensor", "transformador"},
				"details":       err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener activos por tipo",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tipo":    tipo,
		"activos": activos,
		"total":   len(activos),
	})
}

// GetAllActivos obtiene todos los activos
// @Summary Obtener todos los activos
// @Description Obtiene una lista completa de todos los activos en el sistema
// @Tags activos
// @Accept json
// @Produce json
// @Success 200 {array} models.Activo
// @Failure 500 {object} gin.H
// @Router /activos [get]
func (h *ActivoHandler) GetAllActivos(c *gin.Context) {
	activos, err := h.activoService.GetAllActivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener activos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activos": activos,
		"total":   len(activos),
	})
}

// GetActivoByID obtiene un activo específico por su ID
// @Summary Obtener activo por ID
// @Description Obtiene la información completa de un activo específico
// @Tags activos
// @Accept json
// @Produce json
// @Param id path int true "ID del activo"
// @Success 200 {object} models.Activo
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /activos/{id} [get]
func (h *ActivoHandler) GetActivoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de activo inválido",
		})
		return
	}

	activo, err := h.activoService.GetActivoByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Activo no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener activo",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, activo)
}
