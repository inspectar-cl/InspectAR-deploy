package handlers

import (
	"gestion/internal/models"
	"gestion/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReporteHandler struct {
	service *services.ReporteService
}

func NewReporteHandler(service *services.ReporteService) *ReporteHandler {
	return &ReporteHandler{service: service}
}

// POST /reportes - Crear nuevo reporte con observaciones
func (h *ReporteHandler) CrearReporte(c *gin.Context) {
	var req models.CreateReporteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	reporte, err := h.service.CrearReporte(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el reporte", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reporte creado exitosamente",
		"reporte": reporte,
	})
}

// GET /reportes/:id - Obtener reporte por ID
func (h *ReporteHandler) ObtenerReporte(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reporte inválido"})
		return
	}

	reporte, err := h.service.ObtenerReportePorID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reporte no encontrado", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reporte)
}

// PUT /reportes/:id/observaciones - Actualizar observaciones de un reporte
func (h *ReporteHandler) ActualizarObservaciones(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reporte inválido"})
		return
	}

	var req models.UpdateObservacionesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	err = h.service.ActualizarObservaciones(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron actualizar las observaciones", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Observaciones actualizadas exitosamente"})
}

// PUT /reportes/:id/revision - Actualizar estado de revisión
func (h *ReporteHandler) ActualizarEstadoRevision(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reporte inválido"})
		return
	}

	var req models.UpdateEstadoRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	err = h.service.ActualizarEstadoRevision(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado de revisión", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de revisión actualizado exitosamente"})
}

// GET /reportes - Obtener todos los reportes con observaciones
func (h *ReporteHandler) ObtenerTodosLosReportes(c *gin.Context) {
	reportes, err := h.service.ObtenerTodosLosReportes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los reportes", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reportes": reportes,
		"total":    len(reportes),
	})
}

// GET /reportes/activo/:activo_id/observaciones - Obtener reportes con observaciones por activo
func (h *ReporteHandler) ObtenerReportesConObservacionesPorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	reportes, err := h.service.ObtenerReportesConObservaciones(activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los reportes", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activo_id": activoID,
		"reportes":  reportes,
		"total":     len(reportes),
	})
}

// POST /reportes/activo/:activo_id - Generar reporte PDF por activo
func (h *ReporteHandler) GenerarReportePorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	// Leer body con campos solicitados
	var req models.GenerarReporteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Si no viene body, asumimos campos por defecto
		req = models.GenerarReporteRequest{Campos: []string{"ubicacion", "historico_mantenimientos", "ultima_acciones", "datos_sensores"}}
	}

	pdfBytes, filename, err := h.service.GenerarReportePDFPorActivo(activoID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el reporte", "details": err.Error()})
		return
	}

	// Configurar headers para descarga del PDF
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GET /reportes/activo/:activo_id - Obtener reportes de un activo
func (h *ReporteHandler) ObtenerReportesPorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	reportes, err := h.service.ObtenerReportesPorActivo(activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los reportes"})
		return
	}

	c.JSON(http.StatusOK, reportes)
}
