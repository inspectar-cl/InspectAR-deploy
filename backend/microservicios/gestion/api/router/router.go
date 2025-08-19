package router

import (
	"gestion/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(tecnicoHandler *handlers.TecnicoHandler, accionHandler *handlers.AccionMantenimientoHandler) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gestion"})
	})

	// Rutas de técnicos (HdU16 - Lista de contactos)
	r.POST("/tecnicos", tecnicoHandler.CrearTecnico)
	r.GET("/tecnicos", tecnicoHandler.ListarTecnicos)
	r.GET("/tecnicos/:id", tecnicoHandler.ObtenerTecnico)
	r.PUT("/tecnicos/:id/disponibilidad", tecnicoHandler.ActualizarDisponibilidad)

	// Rutas de acciones de mantenimiento (HdU13 - Acciones colaborativas)
	r.POST("/acciones", accionHandler.CrearAccion)
	r.GET("/activos/:id/acciones", accionHandler.ObtenerAccionesPorActivo)
	r.PUT("/acciones/:id/estado", accionHandler.ActualizarEstado)
	r.GET("/acciones/pendientes", accionHandler.ObtenerAccionesPendientes)

	// TODO: Rutas de reportes (HdU04 - Reportes automáticos)
	// r.GET("/reportes", reporteHandler.GenerarReporte)
	// r.GET("/reportes/activo/:id", reporteHandler.ReportePorActivo)

	return r
}
