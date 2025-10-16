package handlers

import (
	"net/http"
	"strconv"
	"time"

	"documentacion/internal/models"
	"documentacion/internal/services"

	"github.com/gin-gonic/gin"
)

type ConsultaHandler struct {
	consultaService  *services.ConsultaService
	documentoService *services.DocumentoService
}

func NewConsultaHandler(consultaService *services.ConsultaService, documentoService *services.DocumentoService) *ConsultaHandler {
	return &ConsultaHandler{
		consultaService:  consultaService,
		documentoService: documentoService,
	}
}

// ConsultarDocumento maneja consultas específicas sobre documentos usando IA
func (h *ConsultaHandler) ConsultarDocumento(c *gin.Context) {
	// Obtener ID del documento
	documentoIDStr := c.Param("id")
	documentoID, err := strconv.Atoi(documentoIDStr)
	if err != nil || documentoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de documento inválido - debe ser un número positivo",
			"code":    "INVALID_DOCUMENT_ID",
			"details": err.Error(),
		})
		return
	}

	// Verificar que el documento existe antes de procesar la consulta
	_, err = h.documentoService.ObtenerDocumento(documentoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Documento no encontrado",
			"code":    "DOCUMENT_NOT_FOUND",
			"details": "El documento especificado no existe",
		})
		return
	}

	// Parsear request
	var request models.RequestConsulta
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de consulta inválidos",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	// Validar pregunta
	if err := request.ValidarPregunta(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Pregunta inválida",
			"code":    "INVALID_QUESTION",
			"details": err.Error(),
		})
		return
	}

	// Verificar que el documento existe
	documento, err := h.documentoService.ObtenerDocumento(documentoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Documento no encontrado",
			"code":    "DOCUMENT_NOT_FOUND",
			"details": err.Error(),
		})
		return
	}

	// Buscar consultas similares primero (optimización)
	consultasSimilares, err := h.consultaService.BuscarConsultasSimilares(documentoID, request.Pregunta)
	if err == nil && len(consultasSimilares) > 0 && consultasSimilares[0].Similaridad > 0.8 {
		// Si hay una consulta muy similar, devolverla
		consulta := consultasSimilares[0]
		response := models.ResponseConsulta{
			Pregunta:    request.Pregunta,
			Respuesta:   consulta.Respuesta,
			DocumentoID: documentoID,
			Confianza:   &consulta.Confianza,
			ProcesadoEn: time.Now(),
		}

		c.Header("X-Consulta-Cache", "similar")
		c.JSON(http.StatusOK, response)
		return
	}

	// Realizar nueva consulta con IA
	inicio := time.Now()
	respuesta, err := h.consultaService.ConsultarConIA(documentoID, request.Pregunta, documento)
	tiempoMs := int(time.Since(inicio).Milliseconds())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al procesar consulta con IA",
			"code":    "IA_PROCESSING_ERROR",
			"details": err.Error(),
		})
		return
	}

	// Configurar tiempo de respuesta
	respuesta.TiempoRespuestaMs = &tiempoMs
	respuesta.ProcesadoEn = time.Now()

	// Guardar consulta para futuras referencias
	consulta := &models.ConsultaIA{
		DocumentoID:       documentoID,
		Pregunta:          request.Pregunta,
		Respuesta:         respuesta.Respuesta,
		Confianza:         respuesta.Confianza,
		TiempoRespuestaMs: &tiempoMs,
		CreadoEn:          time.Now(),
	}

	if respuesta.Fuentes != nil {
		consulta.SetFuentesFromSlice(respuesta.Fuentes)
	}

	// Guardar en base de datos (no bloquear respuesta si falla)
	go func() {
		if err := h.consultaService.GuardarConsulta(consulta); err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}()

	c.JSON(http.StatusOK, respuesta)
}

// ObtenerHistorialConsultas obtiene el historial de consultas de un documento
func (h *ConsultaHandler) ObtenerHistorialConsultas(c *gin.Context) {
	// Obtener ID del documento
	documentoIDStr := c.Param("id")
	documentoID, err := strconv.Atoi(documentoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID de documento inválido",
			"code":    "INVALID_DOCUMENT_ID",
			"details": err.Error(),
		})
		return
	}

	// Parámetros de paginación
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Verificar que el documento existe
	_, err = h.documentoService.ObtenerDocumento(documentoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Documento no encontrado",
			"code":    "DOCUMENT_NOT_FOUND",
			"details": err.Error(),
		})
		return
	}

	// Obtener historial
	historial, err := h.consultaService.ObtenerHistorial(documentoID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener historial",
			"code":    "HISTORIAL_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, historial)
}

// ObtenerEstadisticasConsultas obtiene estadísticas globales de consultas
func (h *ConsultaHandler) ObtenerEstadisticasConsultas(c *gin.Context) {
	estadisticas, err := h.consultaService.ObtenerEstadisticas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estadísticas",
			"code":    "STATS_ERROR",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, estadisticas)
}
