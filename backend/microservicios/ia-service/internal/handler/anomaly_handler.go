package handler

import (
	"ia-service/internal/models"
	"ia-service/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnomalyHandler struct {
	service *services.AnomalyService
}

// NewAnomalyHandler crea una nueva instancia del handler
func NewAnomalyHandler(service *services.AnomalyService) *AnomalyHandler {
	return &AnomalyHandler{
		service: service,
	}
}

// StoreAnomaly maneja POST /anomalies/store
// @Summary Guardar una anomalía
// @Description Guarda los resultados de predicción de una anomalía en la base de datos
// @Tags anomalies
// @Accept json
// @Produce json
// @Param anomaly body models.AnomalyInput true "Datos de la anomalía"
// @Success 201 {object} models.Anomaly
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /anomalies/store [post]
func (h *AnomalyHandler) StoreAnomaly(c *gin.Context) {
	var input models.AnomalyInput

	// Validar JSON de entrada
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("❌ Error validando JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	log.Printf("📥 Recibiendo anomalía: Sensor=%s, Score=%.4f, Severidad=%s",
		input.SensorID, input.AnomalyScore, input.Severidad)

	// Guardar en la base de datos
	savedAnomaly, err := h.service.SaveAnomaly(&input)
	if err != nil {
		log.Printf("❌ Error guardando anomalía: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al guardar la anomalía",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Anomalía guardada exitosamente: ID=%d", savedAnomaly.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Anomalía guardada exitosamente",
		"data":    savedAnomaly,
	})
}

// GetStatus maneja GET /status/
// @Summary Obtener última anomalía
// @Description Retorna la anomalía más reciente detectada
// @Tags status
// @Produce json
// @Success 200 {object} models.Anomaly
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /status/ [get]
func (h *AnomalyHandler) GetStatus(c *gin.Context) {
	log.Println("🔍 Consultando última anomalía...")

	anomaly, err := h.service.GetLatestAnomaly()
	if err != nil {
		if err.Error() == "no se encontraron anomalías" {
			log.Println("ℹ️  No hay anomalías registradas")
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No se encontraron anomalías",
			})
			return
		}

		log.Printf("❌ Error obteniendo última anomalía: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener la última anomalía",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Última anomalía: ID=%d, Sensor=%s, Timestamp=%s",
		anomaly.ID, anomaly.SensorID, anomaly.Timestamp.Format("2006-01-02 15:04:05"))

	c.JSON(http.StatusOK, gin.H{
		"message": "Última anomalía obtenida exitosamente",
		"data":    anomaly,
	})
}

// GetAnomaliesBySensor maneja GET /anomalies/sensor/:sensor_id
// @Summary Obtener anomalías por sensor
// @Description Retorna las anomalías de un sensor específico
// @Tags anomalies
// @Produce json
// @Param sensor_id path string true "ID del sensor"
// @Param limit query int false "Límite de resultados" default(10)
// @Success 200 {array} models.Anomaly
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /anomalies/sensor/{sensor_id} [get]
func (h *AnomalyHandler) GetAnomaliesBySensor(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	limit := 10

	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := c.GetQuery("limit"); err {
			limit = 10
		}
	}

	log.Printf("🔍 Consultando anomalías para sensor: %s (limit=%d)", sensorID, limit)

	anomalies, err := h.service.GetAnomaliesBySensor(sensorID, limit)
	if err != nil {
		log.Printf("❌ Error obteniendo anomalías por sensor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener anomalías",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Encontradas %d anomalías para sensor %s", len(anomalies), sensorID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Anomalías obtenidas exitosamente",
		"count":     len(anomalies),
		"sensor_id": sensorID,
		"data":      anomalies,
	})
}

// HealthCheck maneja GET /health
// @Summary Health check
// @Description Verifica que el servicio esté funcionando
// @Tags health
// @Produce json
// @Success 200 {object} gin.H
// @Router /health [get]
func (h *AnomalyHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "ia-service",
		"message": "Servicio de IA funcionando correctamente",
	})
}
