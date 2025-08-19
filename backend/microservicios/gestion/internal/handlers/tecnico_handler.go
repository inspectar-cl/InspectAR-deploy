package handlers

import (
	"gestion/internal/models"
	"gestion/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TecnicoHandler struct {
	service *services.TecnicoService
}

func NewTecnicoHandler(service *services.TecnicoService) *TecnicoHandler {
	return &TecnicoHandler{service: service}
}

// POST /tecnicos
func (h *TecnicoHandler) CrearTecnico(c *gin.Context) {
	var req models.CreateTecnicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	tecnico, err := h.service.CrearTecnico(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el técnico", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tecnico)
}

// GET /tecnicos
func (h *TecnicoHandler) ListarTecnicos(c *gin.Context) {
	tecnicos, err := h.service.ListarTecnicos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los técnicos"})
		return
	}

	c.JSON(http.StatusOK, tecnicos)
}

// GET /tecnicos/:id
func (h *TecnicoHandler) ObtenerTecnico(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	tecnico, err := h.service.ObtenerTecnico(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Técnico no encontrado"})
		return
	}

	c.JSON(http.StatusOK, tecnico)
}

// PUT /tecnicos/:id/disponibilidad
func (h *TecnicoHandler) ActualizarDisponibilidad(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Disponible bool `json:"disponible"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.service.ActualizarDisponibilidad(id, req.Disponible); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la disponibilidad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Disponibilidad actualizada correctamente"})
}
