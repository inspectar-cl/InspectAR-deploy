package router

import (
	"gestion/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(tecnicoHandler *handlers.TecnicoHandler, accionHandler *handlers.AccionMantenimientoHandler, reporteHandler *handlers.ReporteHandler) *gin.Engine {
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

	// Rutas de técnicos (HdU16 - Lista de contactos de técnicos especializados)
	r.POST("/tecnicos", tecnicoHandler.CrearTecnico)
	r.GET("/tecnicos", tecnicoHandler.ListarTodosLosTecnicos)
	r.GET("/tecnicos/activo/:activo_id", tecnicoHandler.ListarTecnicosPorActivo)       // Técnicos relacionados con un activo
	r.GET("/tecnicos/edificio/:edificio_id", tecnicoHandler.ListarTecnicosPorEdificio) // Técnicos relacionados con un edificio
	r.GET("/tecnicos/:id", tecnicoHandler.ObtenerTecnico)
	r.PUT("/tecnicos/:id/autorizado", tecnicoHandler.ActualizarAutorizado)
	r.POST("/activos/:activo_id/tecnicos", tecnicoHandler.AsignarTecnicoAActivo) // Asignar técnico a activo

	// Rutas de acciones de mantenimiento (HdU13 - Acciones de mantención colaborativas)
	r.POST("/acciones", accionHandler.CrearAccion)
	r.GET("/acciones/tecnico/:tecnico_id", accionHandler.ObtenerAccionesPorTecnico) // Acciones de un técnico específico
	r.GET("/acciones/activo/:activo_id", accionHandler.ObtenerAccionesPorActivo)    // Acciones por activo
	r.PUT("/acciones/:id/estado", accionHandler.ActualizarEstado)
	r.GET("/acciones/pendientes", accionHandler.ObtenerAccionesPendientes) // Acciones pendientes con prioridad

	// Rutas de reportes (HdU04 - Reportes automáticos)
	r.POST("/reportes/activo/:activo_id", reporteHandler.GenerarReportePorActivo) // Generar reporte PDF por activo
	r.GET("/reportes/activo/:activo_id", reporteHandler.ObtenerReportesPorActivo) // Obtener reportes de un activo

	return r
}
