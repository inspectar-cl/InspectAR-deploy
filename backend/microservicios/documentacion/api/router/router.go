package router

import (
	"documentacion/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(documentoHandler *handlers.DocumentoHandler, consultaHandler *handlers.ConsultaHandler) *gin.Engine {
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
		c.JSON(200, gin.H{"status": "ok", "service": "documentacion"})
	})

	// Test route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Microservicio de Documentación funcionando"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Rutas de documentos (HdU05 - Subida de documentación técnica)
		// NOTA: Las rutas específicas DEBEN ir antes que las parametrizadas
		v1.GET("/documentos", documentoHandler.ListarDocumentos)                // Listar todos los documentos
		v1.GET("/documentos/buscar", documentoHandler.BuscarDocumentos)         // Buscar con filtros
		v1.POST("/documentos", documentoHandler.SubirDocumento)                 // Subir documento
		v1.GET("/documentos/:id", documentoHandler.ObtenerDocumento)            // Obtener documento por ID
		v1.PUT("/documentos/:id", documentoHandler.ActualizarDocumento)         // Actualizar documento
		v1.GET("/documentos/:id/download", documentoHandler.DescargarDocumento) // Descargar archivo

		// Rutas por activo
		v1.GET("/documentos/activo/:activo_id", documentoHandler.ObtenerDocumentosPorActivo) // Documentos de un activo
		v1.GET("/documentos/fichas-tecnicas", documentoHandler.ObtenerFichaTecnica)          // Solo fichas técnicas

		// Rutas de análisis IA (HdU19 - Lectura y análisis con Gemini)
		v1.POST("/documentos/:id/analizar", documentoHandler.AnalizarDocumentoIA) // Solicitar análisis IA
		v1.GET("/documentos/:id/analisis", documentoHandler.ObtenerAnalisisIA)    // Obtener análisis IA

		// Rutas de consultas interactivas (NUEVA FUNCIONALIDAD)
		v1.POST("/documentos/:id/consultar", consultaHandler.ConsultarDocumento)        // Hacer pregunta específica
		v1.GET("/documentos/:id/consultas", consultaHandler.ObtenerHistorialConsultas)  // Historial de consultas
		v1.GET("/consultas/estadisticas", consultaHandler.ObtenerEstadisticasConsultas) // Estadísticas globales
	}

	return r
}
