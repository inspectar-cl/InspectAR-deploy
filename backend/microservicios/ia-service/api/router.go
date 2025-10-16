package api

import (
	"ia-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas del servicio
func SetupRouter(anomalyHandler *handler.AnomalyHandler) *gin.Engine {
	router := gin.Default()

	// Middleware para logging y CORS
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", anomalyHandler.HealthCheck)

	// Grupo de rutas para anomalías
	anomalies := router.Group("/anomalies")
	{
		// POST /anomalies/store - Guardar anomalía
		anomalies.POST("/store", anomalyHandler.StoreAnomaly)

		// GET /anomalies/sensor/:sensor_id - Obtener anomalías por sensor
		anomalies.GET("/sensor/:sensor_id", anomalyHandler.GetAnomaliesBySensor)
	}

	// Ruta de status
	router.GET("/status/", anomalyHandler.GetStatus)

	return router
}
