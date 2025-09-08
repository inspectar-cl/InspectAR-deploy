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

// POST /tecnicos - Crear técnico con especialidades
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

// GET /tecnicos - Listar todos los técnicos
func (h *TecnicoHandler) ListarTodosLosTecnicos(c *gin.Context) {
	tecnicos, err := h.service.ListarTodosLosTecnicos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los técnicos"})
		return
	}

	c.JSON(http.StatusOK, tecnicos)
}

// GET /tecnicos/activo/:activo_id - Listar técnicos relacionados con un activo específico
func (h *TecnicoHandler) ListarTecnicosPorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil || activoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	soloAutorizados := c.Query("solo_autorizados") == "true"

	tecnicos, err := h.service.ListarTecnicosPorActivo(activoID, soloAutorizados)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los técnicos del activo"})
		return
	}

	c.JSON(http.StatusOK, tecnicos)
}

// GET /tecnicos/edificio/:edificio_id - Listar técnicos relacionados con un edificio (indirectamente a través de activos)
func (h *TecnicoHandler) ListarTecnicosPorEdificio(c *gin.Context) {
	edificioID, err := strconv.Atoi(c.Param("edificio_id"))
	if err != nil || edificioID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	soloAutorizados := c.Query("solo_autorizados") == "true"

	tecnicos, err := h.service.ListarTecnicosPorEdificio(edificioID, soloAutorizados)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los técnicos del edificio"})
		return
	}

	c.JSON(http.StatusOK, tecnicos)
}

// GET /tecnicos/:id - Consultar información de un técnico específico
func (h *TecnicoHandler) ObtenerTecnico(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
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

// PUT /tecnicos/:id/autorizado - Actualizar estado autorizado de técnicos
func (h *TecnicoHandler) ActualizarAutorizado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateAutorizadoTecnicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.service.ActualizarAutorizado(id, req.Autorizado); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado autorizado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Estado autorizado actualizado correctamente"})
}

// POST /activos/:activo_id/tecnicos - Asignar técnico a activo (relación muchos a muchos)
func (h *TecnicoHandler) AsignarTecnicoAActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	var req models.AsignarTecnicoActivoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.service.AsignarTecnicoAActivo(activoID, req.TecnicoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo asignar el técnico al activo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensaje": "Técnico asignado al activo correctamente"})
}

// GET /tecnicos/:tecnico_id/activos - Obtener activos asociados a un técnico (RUTA PRINCIPAL)
func (h *TecnicoHandler) ObtenerActivosPorTecnico(c *gin.Context) {
	tecnicoID, err := strconv.Atoi(c.Param("tecnico_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de técnico inválido"})
		return
	}

	activos, err := h.service.ObtenerActivosPorTecnico(tecnicoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los activos del técnico", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activos)
}

// ObtenerActivosDirecto - Método público para acceder directamente desde el router
func (h *TecnicoHandler) ObtenerActivosDirecto(tecnicoID int) (interface{}, error) {
	return h.service.ObtenerActivosPorTecnico(tecnicoID)
}
