package handlers

import (
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

// POST /reportes/activo/:activo_id - Generar reporte PDF por activo
func (h *ReporteHandler) GenerarReportePorActivo(c *gin.Context) {
	activoID, err := strconv.Atoi(c.Param("activo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	pdfBytes, filename, err := h.service.GenerarReportePDFPorActivo(activoID)
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
