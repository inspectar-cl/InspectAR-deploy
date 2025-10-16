package handlers

import (
	"ParserService/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SensorStatusHandler struct {
	monitoringService *services.SensorMonitoringService
}

func NewSensorStatusHandler(monitoringService *services.SensorMonitoringService) *SensorStatusHandler {
	return &SensorStatusHandler{
		monitoringService: monitoringService,
	}
}

// GET /api/sensors/status - Obtiene el estado de todos los sensores
func (h *SensorStatusHandler) GetAllSensorStatus(c *gin.Context) {
	sensors, err := h.monitoringService.GetAllSensorStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estado de sensores",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sensors": sensors,
		"total":   len(sensors),
	})
}

// GET /api/sensors/status/:sensor_id - Obtiene el estado de un sensor específico
func (h *SensorStatusHandler) GetSensorStatus(c *gin.Context) {
	sensorID := c.Param("sensor_id")
	if sensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sensor_id es requerido",
		})
		return
	}

	sensor, err := h.monitoringService.GetSensorStatus(sensorID)
	if err != nil {
		if err.Error() == "sensor "+sensorID+" no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sensor no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estado del sensor",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sensor": sensor,
	})
}

// POST /api/sensors/check-disconnected - Fuerza verificación manual de sensores desconectados
func (h *SensorStatusHandler) CheckDisconnectedSensors(c *gin.Context) {
	err := h.monitoringService.CheckDisconnectedSensors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar sensores desconectados",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verificación de sensores desconectados completada",
	})
}

// GET /api/sensors/stats - Obtiene estadísticas generales de sensores
func (h *SensorStatusHandler) GetSensorStats(c *gin.Context) {
	stats, err := h.monitoringService.GetSensorStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estadísticas de sensores",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GET /api/sensors/health - Health check específico para el monitoreo de sensores
func (h *SensorStatusHandler) HealthCheck(c *gin.Context) {
	stats, err := h.monitoringService.GetSensorStats()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "Error al acceder al sistema de monitoreo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "healthy",
		"monitoring":       "active",
		"total_sensors":    stats["total_sensors"],
		"active_sensors":   stats["active_sensors"],
		"inactive_sensors": stats["inactive_sensors"],
		"last_check":       stats["last_check"],
	})
}
