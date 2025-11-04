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
	email := c.PostForm("email")
	nombreArchivo := c.PostForm("nombre_archivo")
	esPredeterminadaStr := c.PostForm("es_predeterminada")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email requerido"})
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
	firma, err := h.service.CrearFirmaDesdeArchivo(email, nombreArchivo, esPredeterminada, file, header)
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

	firma, err := h.service.CrearFirmaDesdeSVG(req.Email, req.NombreArchivo, req.DatosSVG, req.EsPredeterminada)
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

// GET /firmas/usuario/:email - Obtener todas las firmas de un usuario por email
func (h *FirmaHandler) ObtenerFirmasUsuario(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email requerido"})
		return
	}

	firmas, err := h.service.ObtenerFirmasUsuarioPorEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firmas": firmas,
		"total":  len(firmas),
	})
}

// GET /firmas/usuario/:email/predeterminada - Obtener firma predeterminada de un usuario por email
func (h *FirmaHandler) ObtenerFirmaPredeterminada(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email requerido"})
		return
	}

	firma, err := h.service.ObtenerFirmaPredeterminadaPorEmail(email)
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

	// Obtener email del body para validar permisos
	type DeleteFirmaRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req DeleteFirmaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email requerido en el body"})
		return
	}

	// Validar que la firma pertenezca al usuario antes de eliminar
	err = h.service.EliminarFirmaConValidacion(id, req.Email)
	if err != nil {
		if err.Error() == "firma no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Firma no encontrada"})
			return
		}
		if err.Error() == "no tiene permisos para eliminar esta firma" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permisos para eliminar esta firma"})
			return
		}
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

	// Obtener email del body o query
	type SetDefaultRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req SetDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email requerido"})
		return
	}

	err = h.service.EstablecerComoPredeterminadaPorEmail(id, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Firma %d establecida como predeterminada", id),
	})
}
