package handlers

import (
	"ParserService/internal/models"
	"ParserService/internal/repository"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	activoRepo *repository.ActivoRepository
}

func NewAdminHandler(activoRepo *repository.ActivoRepository) *AdminHandler {
	return &AdminHandler{
		activoRepo: activoRepo,
	}
}

// CreateActivoAdmin crea o actualiza un activo desde gestion-service
func (h *AdminHandler) CreateActivoAdmin(c *gin.Context) {
	var req struct {
		IDActivo int    `json:"id_activo" binding:"required"`
		Nombre   string `json:"nombre" binding:"required"`
		Tipo     string `json:"tipo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si existe
	existing, _ := h.activoRepo.GetActivo(context.Background(), req.IDActivo)
	if existing != nil {
		// Si ya existe, solo retornar éxito (idempotencia)
		c.JSON(http.StatusOK, gin.H{
			"message": "Activo ya existe en ParserService",
			"activo":  existing,
		})
		return
	}

	// Crear activo nuevo
	activo := &models.Activo{
		ActivoID:   req.IDActivo,
		Estado:     "operativo",
		EdificioID: 0, // Se puede actualizar después si es necesario
	}

	_, err := h.activoRepo.CreateActivo(context.Background(), activo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear activo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Activo creado exitosamente en ParserService",
		"activo":  activo,
	})
}

// DeleteActivoAdmin elimina un activo desde gestion-service
func (h *AdminHandler) DeleteActivoAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar que existe
	activo, err := h.activoRepo.GetActivo(context.Background(), id)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Eliminar
	err = h.activoRepo.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar activo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activo eliminado exitosamente de ParserService"})
}

// CreateSensorAdmin crea un sensor desde gestion-service
func (h *AdminHandler) CreateSensorAdmin(c *gin.Context) {
	var req struct {
		IDActivo int    `json:"id_activo" binding:"required"`
		SensorID string `json:"sensor_id" binding:"required"` // ID generado por gestion-service
		Nombre   string `json:"nombre"`                       // Opcional - descripción del sensor
		Tipo     string `json:"tipo" binding:"required"`
		Unidad   string `json:"unidad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el activo existe
	activo, err := h.activoRepo.GetActivo(context.Background(), req.IDActivo)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado en ParserService"})
		return
	}

	// Crear sensor con el ID generado por gestion-service
	sensor := models.Sensor{
		SensorID: req.SensorID, // Usar el ID generado automáticamente
		Tipo:     req.Tipo,
		Unidad:   req.Unidad,
	}

	// Agregar sensor al activo
	err = h.activoRepo.AddSensor(context.Background(), req.IDActivo, sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear sensor", "details": err.Error()})
		return
	}

	// Obtener activo actualizado para obtener el ID del sensor
	activoActualizado, _ := h.activoRepo.GetActivo(context.Background(), req.IDActivo)

	response := gin.H{
		"message": "Sensor creado exitosamente en ParserService",
		"sensor":  sensor,
		"activo":  activoActualizado,
	}

	// Agregar nombre en la respuesta si existe
	if req.Nombre != "" {
		response["nombre"] = req.Nombre
	}

	c.JSON(http.StatusCreated, response)
}
