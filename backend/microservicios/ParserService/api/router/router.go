package router

import (
	handlers "ParserService/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(dataHandler *handlers.DataHandler, statusHandler *handlers.SensorStatusHandler) *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Rutas existentes de activos y lecturas
	r.POST("/activo", dataHandler.CreateActivo)
	r.GET("/activo", dataHandler.GetAllActivos)
	r.POST("/lectura", dataHandler.CreateLectura)
	r.GET("/activo/:activo_id", dataHandler.GetActivo)
	r.GET("/activo/:activo_id/sensores/estado", dataHandler.GetActivoWithSensorStatus)
	r.POST("/activo/:activo_id/sensores", dataHandler.AddSensorToActivo)
	r.GET("/lectura/:activo_id/datos", dataHandler.GetSensorByActivo)
	r.GET("/lectura/:activo_id/datos/ultimo", dataHandler.GetSensorLastData)
	r.PUT("/activo/:activo_id/estado", dataHandler.UpdActivoEstado)

	// Nuevas rutas de monitoreo de sensores
	api := r.Group("/api")
	{
		sensors := api.Group("/sensors")
		{
			sensors.GET("/status", statusHandler.GetAllSensorStatus)
			sensors.GET("/status/:sensor_id", statusHandler.GetSensorStatus)
			sensors.POST("/check-disconnected", statusHandler.CheckDisconnectedSensors)
			sensors.GET("/stats", statusHandler.GetSensorStats)
			sensors.GET("/health", statusHandler.HealthCheck)
		}
	}

	return r
}
