package router

import (
	"gestion/internal/handlers"
	"strconv"
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

	// Test route for debugging
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test route works"})
	})

	// Rutas de técnicos (HdU16 - Lista de contactos de técnicos especializados)
	r.POST("/tecnicos", tecnicoHandler.CrearTecnico)
	r.GET("/tecnicos", tecnicoHandler.ListarTodosLosTecnicos)
	r.GET("/tecnicos/activo/:activo_id", tecnicoHandler.ListarTecnicosPorActivo)       // Técnicos relacionados con un activo
	r.GET("/tecnicos/edificio/:edificio_id", tecnicoHandler.ListarTecnicosPorEdificio) // Técnicos relacionados con un edificio

	// RUTA PRINCIPAL: Activos por técnico (usando /activos-de-tecnico/ para evitar conflicto)
	r.GET("/activos-de-tecnico/:tecnico_id", func(c *gin.Context) {
		tecnicoID, err := strconv.Atoi(c.Param("tecnico_id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "ID de técnico inválido"})
			return
		}

		// Llamar al método público del handler
		activos, err := tecnicoHandler.ObtenerActivosDirecto(tecnicoID)
		if err != nil {
			c.JSON(500, gin.H{"error": "No se pudieron obtener los activos del técnico", "details": err.Error()})
			return
		}

		c.JSON(200, activos)
	})

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

// Checking current router configuration for técnicos endpoints
