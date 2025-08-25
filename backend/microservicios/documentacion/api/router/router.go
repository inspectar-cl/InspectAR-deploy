package router

import (
	"documentacion/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(documentoHandler *handlers.DocumentoHandler) *gin.Engine {
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

	// Rutas de documentos (HdU05 - Subida de documentación técnica)
	r.POST("/documentos", documentoHandler.SubirDocumento)                                    // Subir documento
	r.GET("/documentos/:id", documentoHandler.ObtenerDocumento)                              // Obtener documento por ID
	r.GET("/documentos/:id/descargar", documentoHandler.DescargarDocumento)                  // Descargar archivo
	r.GET("/documentos/buscar", documentoHandler.BuscarDocumentos)                           // Buscar con filtros
	
	// Rutas por activo
	r.GET("/activos/:activo_id/documentos", documentoHandler.ObtenerDocumentosPorActivo)     // Documentos de un activo
	r.GET("/activos/:activo_id/historial", documentoHandler.ObtenerHistorialMantenimiento)  // Historial completo
	
	// Rutas de fichas técnicas (HdU23 - Fichas técnicas de activos)
	r.GET("/activos/:activo_id/ficha-tecnica", documentoHandler.ObtenerFichaTecnica)         // Ficha técnica específica
	
	// Rutas de análisis IA (HdU19 - Lectura y análisis con Gemini)
	r.POST("/documentos/:id/analizar", documentoHandler.AnalizarDocumentoIA)                 // Solicitar análisis IA
	r.GET("/documentos/:id/analisis", documentoHandler.ObtenerAnalisisIA)                    // Obtener análisis IA

	return r
}
