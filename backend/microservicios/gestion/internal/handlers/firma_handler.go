package handlers

import (
	"fmt"
	"gestion/internal/models"
	"gestion/internal/services"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FirmaHandler struct {
	service *services.FirmaService
}

func NewFirmaHandler(service *services.FirmaService) *FirmaHandler {
	return &FirmaHandler{service: service}
}

// POST /firmas/upload - Subir una firma como archivo (imagen)
func (h *FirmaHandler) SubirFirma(c *gin.Context) {
	// Obtener parámetros del form
	usuarioIDStr := c.PostForm("usuario_id")
	nombreArchivo := c.PostForm("nombre_archivo")
	esPredeterminadaStr := c.PostForm("es_predeterminada")

	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario_id inválido"})
		return
	}

	esPredeterminada := esPredeterminadaStr == "true"

	// Obtener archivo
	file, header, err := c.Request.FormFile("archivo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo requerido"})
		return
	}
	defer file.Close()

	// Crear firma
	firma, err := h.service.CrearFirmaDesdeArchivo(usuarioID, nombreArchivo, esPredeterminada, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Firma subida exitosamente",
		"firma":   firma,
	})
}

// POST /firmas/svg - Crear una firma desde datos SVG (pizarra)
func (h *FirmaHandler) CrearFirmaSVG(c *gin.Context) {
	var req models.CreateFirmaSVGRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firma, err := h.service.CrearFirmaDesdeSVG(req.UsuarioID, req.NombreArchivo, req.DatosSVG, req.EsPredeterminada)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Firma creada exitosamente",
		"firma":   firma,
	})
}

// GET /firmas/:id - Obtener una firma por ID
func (h *FirmaHandler) ObtenerFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	firma, err := h.service.ObtenerFirma(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Firma no encontrada"})
		return
	}

	c.JSON(http.StatusOK, firma)
}

// GET /firmas/:id/imagen - Obtener la imagen de la firma
func (h *FirmaHandler) ObtenerImagenFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	firma, err := h.service.ObtenerFirma(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Firma no encontrada"})
		return
	}

	// Verificar que el archivo existe
	if _, err := os.Stat(firma.RutaArchivo); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo de firma no encontrado"})
		return
	}

	// Servir el archivo con el tipo MIME correcto
	c.Header("Content-Type", firma.TipoMime)
	c.File(firma.RutaArchivo)
}

// GET /firmas/usuario/:usuario_id - Obtener todas las firmas de un usuario
func (h *FirmaHandler) ObtenerFirmasUsuario(c *gin.Context) {
	usuarioIDStr := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario_id inválido"})
		return
	}

	firmas, err := h.service.ObtenerFirmasUsuario(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firmas": firmas,
		"total":  len(firmas),
	})
}

// GET /firmas/usuario/:usuario_id/predeterminada - Obtener firma predeterminada de un usuario
func (h *FirmaHandler) ObtenerFirmaPredeterminada(c *gin.Context) {
	usuarioIDStr := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario_id inválido"})
		return
	}

	firma, err := h.service.ObtenerFirmaPredeterminada(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró firma predeterminada"})
		return
	}

	c.JSON(http.StatusOK, firma)
}

// PUT /firmas/:id - Actualizar una firma
func (h *FirmaHandler) ActualizarFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateFirmaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firma, err := h.service.ActualizarFirma(id, req.NombreArchivo, req.EsPredeterminada)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Firma actualizada exitosamente",
		"firma":   firma,
	})
}

// DELETE /firmas/:id - Eliminar una firma
func (h *FirmaHandler) EliminarFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.EliminarFirma(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Firma eliminada exitosamente",
	})
}

// POST /firmas/:id/predeterminada - Establecer una firma como predeterminada
func (h *FirmaHandler) EstablecerComoPredeterminada(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener usuario_id del body o query
	type SetDefaultRequest struct {
		UsuarioID int `json:"usuario_id" binding:"required"`
	}

	var req SetDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario_id requerido"})
		return
	}

	err = h.service.EstablecerComoPredeterminada(id, req.UsuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Firma %d establecida como predeterminada", id),
	})
}
