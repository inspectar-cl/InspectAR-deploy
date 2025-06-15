package router

import (
	"backend/microservicios/ParserService/internal/handler"
	"backend/microservicios/ParserService/internal/repository"
	"forms/api/middleware"
	"forms/internal/handler"
	"forms/internal/repository"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dataHandler *handler.DataHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5050"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/activo", dataHandler.CreateActivo)
	r.POST("/lectura", dataHandler.CreateLectura)
	r.GET("/activo/:activo_id", dataHandler.GetActivo)
	r.GET("/activo/:activo_id/datos", dataHandler.GetSensorByActivo)

	return r
}
