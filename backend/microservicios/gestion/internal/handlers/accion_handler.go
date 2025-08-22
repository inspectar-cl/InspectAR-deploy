package handlers

import (
	"gestion/internal/models"
	"gestion/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccionMantenimientoHandler struct {
	service *services.AccionMantenimientoService
}

func NewAccionMantenimientoHandler(service *services.AccionMantenimientoService) *AccionMantenimientoHandler {
	return &AccionMantenimientoHandler{service: service}
}

// POST /acciones - Crear acciones de mantenimiento (preventivo, correctivo, emergencia)
func (h *AccionMantenimientoHandler) CrearAccion(c *gin.Context) {
	var req models.CreateAccionMantenimientoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	accion, err := h.service.CrearAccion(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la acción de mantenimiento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, accion)
}

// GET /acciones/tecnico/:tecnico_id - Acciones donde el técnico puede buscar sus acciones
func (h *AccionMantenimientoHandler) ObtenerAccionesPorTecnico(c *gin.Context) {
	tecnicoID, err := strconv.Atoi(c.Param("tecnico_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de técnico inválido"})
		return
	}

	acciones, err := h.service.ObtenerAccionesPorTecnico(tecnicoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las acciones del técnico"})
		return
	}

	c.JSON(http.StatusOK, acciones)
}

// GET /acciones/activo/:activo_id - Consultar acciones por activo
func (h *AccionMantenimientoHandler) ObtenerAccionesPorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	acciones, err := h.service.ObtenerAccionesPorActivo(activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las acciones"})
		return
	}

	c.JSON(http.StatusOK, acciones)
}

// PUT /acciones/:id/estado - Actualizar estado de acciones (pendiente, en_progreso, completado)
func (h *AccionMantenimientoHandler) ActualizarEstado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateEstadoAccionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.service.ActualizarEstado(id, req.Estado); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Estado actualizado correctamente"})
}

// GET /acciones/pendientes - Listar acciones pendientes con prioridad
func (h *AccionMantenimientoHandler) ObtenerAccionesPendientes(c *gin.Context) {
	acciones, err := h.service.ObtenerAccionesPendientes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las acciones pendientes"})
		return
	}

	c.JSON(http.StatusOK, acciones)
}
