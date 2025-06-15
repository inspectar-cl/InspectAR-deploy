package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"forms/internal/models"
	"forms/internal/services"
)

type DataHandler struct {
	activoService *services.ActivoService
	sensorService *services.SensorService
}

func NewDataHandler(activoSvc *services.ActivoService, sensorSvc *services.SensorService) *DataHandler {
	return &DataHandler{
		activoService: activoSvc,
		sensorService: sensorSvc,
	}
}

// POST /activo
func (h *DataHandler) CreateActivo(c *gin.Context) {
	var activo models.Activo
	if err := c.ShouldBindJSON(&activo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	id, err := h.activoService.CrearActivo(c.Request.Context(), &activo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activo_id": id})
}

// POST /lectura
func (h *DataHandler) CreateLectura(c *gin.Context) {
	var lectura models.LecturaRequest
	if err := c.ShouldBindJSON(&lectura); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de lectura inválido"})
		return
	}

	sensor := models.Sensor{
		SensorID: lectura.SensorID,
		Tipo:     lectura.Tipo,
		Unidad:   lectura.Unidad,
	}

	err := h.activoService.AgregarSensor(c.Request.Context(), lectura.ActivoID, sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo agregar el sensor al activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Sensor agregado al activo"})
}

// GET /activo/:activo_id
func (h *DataHandler) GetActivo(c *gin.Context) {
	activoID := c.Param("activo_id")

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, activo)
}

// GET /activo/:activo_id/datos
func (h *DataHandler) GetSensorByActivo(c *gin.Context) {
	activoID := c.Param("activo_id")

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	var allLecturas []models.SensorDataset

	for _, sensor := range activo.Sensores {
		datos, err := h.sensorService.GetDatosSensor(c.Request.Context(), sensor.SensorID, 30*time.Minute)
		if err != nil {
			continue
		}
		allLecturas = append(allLecturas, models.SensorDataset{
			SensorID: sensor.SensorID,
			Datos:    datos,
		})
	}

	c.JSON(http.StatusOK, allLecturas)
}
