package router

import (
	handlers "ParserService/internal/handler"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dataHandler *handlers.DataHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/activo", dataHandler.CreateActivo)
	r.GET("/activo", dataHandler.GetAllActivos)
	r.POST("/lectura", dataHandler.CreateLectura)
	r.GET("/activo/:activo_id", dataHandler.GetActivo)
	r.GET("/lectura/:activo_id/datos", dataHandler.GetSensorByActivo)
	r.GET("/lectura/:activo_id/datos/ultimo", dataHandler.GetSensorLastData)
	r.PUT("/activo/:activo_id/estado", dataHandler.UpdActivoEstado)

	return r
}
