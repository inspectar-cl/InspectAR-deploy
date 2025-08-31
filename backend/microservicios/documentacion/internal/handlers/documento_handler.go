package handlers

import (
	"documentacion/internal/models"
	"documentacion/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DocumentoHandler struct {
	documentoService *services.DocumentoService
	aiService        *services.AIService
}

func NewDocumentoHandler(documentoService *services.DocumentoService, aiService *services.AIService) *DocumentoHandler {
	return &DocumentoHandler{
		documentoService: documentoService,
		aiService:        aiService,
	}
}

// SubirDocumento maneja la subida de documentos (HdU05)
func (h *DocumentoHandler) SubirDocumento(c *gin.Context) {
	// Obtener archivo
	file, header, err := c.Request.FormFile("archivo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo obtener el archivo", "details": err.Error()})
		return
	}
	defer file.Close()

	// Parsear JSON del formulario
	var req models.CreateDocumentoRequest
	
	// Obtener campos del formulario
	activoIDStr := c.PostForm("activo_id")
	activoID, err := strconv.Atoi(activoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id inválido"})
		return
	}

	req.ActivoID = activoID
	req.Nombre = c.PostForm("nombre")
	req.Descripcion = c.PostForm("descripcion")
	req.Categoria = c.PostForm("categoria")
	req.SubidoPor = c.PostForm("subido_por")
	req.PalabrasClave = c.PostForm("palabras_clave")

	// Parsear fecha de emisión
	fechaEmisionStr := c.PostForm("fecha_emision")
	if fechaEmisionStr != "" {
		fechaEmision, err := time.Parse("2006-01-02", fechaEmisionStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
			return
		}
		req.FechaEmision = fechaEmision
	} else {
		req.FechaEmision = time.Now()
	}

	// Técnico ID (opcional)
	tecnicoIDStr := c.PostForm("tecnico_id")
	if tecnicoIDStr != "" {
		tecnicoID, err := strconv.Atoi(tecnicoIDStr)
		if err == nil {
			req.TecnicoID = &tecnicoID
		}
	}

	// Es ficha técnica
	esFichaTecnicaStr := c.PostForm("es_ficha_tecnica")
	req.EsFichaTecnica = esFichaTecnicaStr == "true"

	// Validaciones básicas
	if req.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre es requerido"})
		return
	}

	if req.Categoria == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Categoría es requerida"})
		return
	}

	// Subir documento
	documento, err := h.documentoService.SubirDocumento(&req, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error subiendo documento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, documento)
}

// ObtenerDocumento obtiene un documento por ID
func (h *DocumentoHandler) ObtenerDocumento(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	documento, err := h.documentoService.ObtenerDocumento(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, documento)
}

// ObtenerDocumentosPorActivo obtiene todos los documentos de un activo
func (h *DocumentoHandler) ObtenerDocumentosPorActivo(c *gin.Context) {
	activoIDStr := c.Param("activo_id")
	activoID, err := strconv.Atoi(activoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	documentos, err := h.documentoService.ObtenerDocumentosPorActivo(activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo documentos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activo_id":  activoID,
		"documentos": documentos,
		"total":      len(documentos),
	})
}

// ObtenerFichaTecnica obtiene la ficha técnica de un activo (HdU23)
func (h *DocumentoHandler) ObtenerFichaTecnica(c *gin.Context) {
	activoIDStr := c.Param("activo_id")
	activoID, err := strconv.Atoi(activoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	fichaTecnica, err := h.documentoService.ObtenerFichaTecnica(activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ficha técnica no encontrada para este activo"})
		return
	}

	c.JSON(http.StatusOK, fichaTecnica)
}

// DescargarDocumento descarga un archivo
func (h *DocumentoHandler) DescargarDocumento(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	file, filename, err := h.documentoService.DescargarArchivo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo no encontrado"})
		return
	}
	defer file.Close()

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	
	// Copiar archivo a la respuesta
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", file, nil)
}

// BuscarDocumentos busca documentos con filtros (HdU05)
func (h *DocumentoHandler) BuscarDocumentos(c *gin.Context) {
	var filtros models.DocumentoFiltros

	// Parsear query parameters
	if activoIDStr := c.Query("activo_id"); activoIDStr != "" {
		if activoID, err := strconv.Atoi(activoIDStr); err == nil {
			filtros.ActivoID = &activoID
		}
	}

	filtros.Categoria = c.Query("categoria")
	filtros.PalabraClave = c.Query("palabra_clave")

	// Fechas
	if fechaDesdeStr := c.Query("fecha_desde"); fechaDesdeStr != "" {
		if fechaDesde, err := time.Parse("2006-01-02", fechaDesdeStr); err == nil {
			filtros.FechaDesde = &fechaDesde
		}
	}

	if fechaHastaStr := c.Query("fecha_hasta"); fechaHastaStr != "" {
		if fechaHasta, err := time.Parse("2006-01-02", fechaHastaStr); err == nil {
			filtros.FechaHasta = &fechaHasta
		}
	}

	// Ficha técnica
	if esFichaTecnicaStr := c.Query("es_ficha_tecnica"); esFichaTecnicaStr != "" {
		esFichaTecnica := esFichaTecnicaStr == "true"
		filtros.EsFichaTecnica = &esFichaTecnica
	}

	// Paginación
	if limiteStr := c.Query("limite"); limiteStr != "" {
		if limite, err := strconv.Atoi(limiteStr); err == nil {
			filtros.Limite = limite
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filtros.Offset = offset
		}
	}

	documentos, err := h.documentoService.BuscarDocumentos(filtros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la búsqueda"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documentos": documentos,
		"total":      len(documentos),
		"filtros":    filtros,
	})
}

// ObtenerHistorialMantenimiento obtiene el historial completo de un activo
func (h *DocumentoHandler) ObtenerHistorialMantenimiento(c *gin.Context) {
	activoIDStr := c.Param("activo_id")
	activoID, err := strconv.Atoi(activoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	historial, err := h.documentoService.ObtenerHistorialMantenimiento(activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo historial"})
		return
	}

	// Obtener análisis de IA si está disponible
	if h.aiService != nil {
		analisisIA, err := h.aiService.ObtenerAnalisisPorActivo(activoID)
		if err == nil {
			historial.AnalisisIA = analisisIA
		}
	}

	c.JSON(http.StatusOK, historial)
}

// AnalizarDocumentoIA solicita análisis de IA de un documento (HdU19)
func (h *DocumentoHandler) AnalizarDocumentoIA(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if h.aiService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio de IA no disponible"})
		return
	}

	analisis, err := h.aiService.AnalizarDocumento(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iniciando análisis", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"mensaje":    "Análisis iniciado",
		"analisis":   analisis,
		"estado":     "procesando",
	})
}

// ObtenerAnalisisIA obtiene el análisis de IA de un documento
func (h *DocumentoHandler) ObtenerAnalisisIA(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if h.aiService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio de IA no disponible"})
		return
	}

	analisis, err := h.aiService.ObtenerAnalisis(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Análisis no encontrado"})
		return
	}

	c.JSON(http.StatusOK, analisis)
}
