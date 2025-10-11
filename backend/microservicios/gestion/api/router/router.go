package router

import (
	"gestion/internal/handlers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	tecnicoHandler *handlers.TecnicoHandler,
	accionHandler *handlers.AccionMantenimientoHandler,
	activoHandler *handlers.ActivoHandler,
	reporteHandler *handlers.ReporteHandler,
	solicitudHandler *handlers.SolicitudHandler,
	usuarioHandler *handlers.UsuarioHandler,
	tipoFallaHandler *handlers.TipoFallaHandler,
	firmaHandler *handlers.FirmaHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

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

	// 🎯 Rutas de activos - NUEVAS RUTAS AGREGADAS
	r.GET("/activos", activoHandler.GetAllActivos)                              // Obtener todos los activos
	r.GET("/activos/:id", activoHandler.GetActivoByID)                          // Obtener activo por ID
	r.GET("/activos/edificio/:edificio_id", activoHandler.GetActivosByEdificio) // Filtrar activos por edificio
	r.GET("/activos/tipo/:tipo", activoHandler.GetActivosByTipo)                // Filtrar activos por tipo

	// 🏢 Rutas de usuarios y edificios
	r.GET("/usuarios/edificios/:email", usuarioHandler.GetEdificiosByEmail)                        // Obtener edificios de un usuario por email
	r.GET("/usuarios/:email/edificio/:edificio_id/acceso", usuarioHandler.VerificarAccesoEdificio) // Verificar acceso a edificio
	r.GET("/usuarios/:email/activo/:activo_id/acceso", usuarioHandler.VerificarAccesoActivo)       // Verificar acceso a activo

	// Rutas de acciones de mantenimiento (HdU13 - Acciones de mantención colaborativas)
	r.POST("/acciones", accionHandler.CrearAccion)
	r.GET("/acciones/tecnico/:tecnico_id", accionHandler.ObtenerAccionesPorTecnico) // Acciones de un técnico específico
	r.GET("/acciones/activo/:activo_id", accionHandler.ObtenerAccionesPorActivo)    // Acciones por activo
	r.PUT("/acciones/:id/estado", accionHandler.ActualizarEstado)
	r.GET("/acciones/pendientes", accionHandler.ObtenerAccionesPendientes) // Acciones pendientes con prioridad

	// Rutas de reportes (HdU04 - Reportes automáticos con observaciones editables)
	// Nuevas rutas para gestión completa de reportes con observaciones
	r.POST("/reportes", reporteHandler.CrearReporte)                                                            // Crear reporte con observaciones
	r.GET("/reportes", reporteHandler.ObtenerTodosLosReportes)                                                  // Obtener todos los reportes
	r.GET("/reportes/:id", reporteHandler.ObtenerReporte)                                                       // Obtener reporte por ID
	r.PUT("/reportes/:id/observaciones", reporteHandler.ActualizarObservaciones)                                // Actualizar observaciones
	r.PUT("/reportes/:id/revision", reporteHandler.ActualizarEstadoRevision)                                    // Actualizar estado de revisión
	r.GET("/reportes/activo/:activo_id/observaciones", reporteHandler.ObtenerReportesConObservacionesPorActivo) // Reportes con observaciones por activo

	// Rutas originales mantenidas para compatibilidad
	r.POST("/reportes/activo/:activo_id", reporteHandler.GenerarReportePorActivo) // Generar reporte PDF por activo
	r.GET("/reportes/activo/:activo_id", reporteHandler.ObtenerReportesPorActivo) // Obtener reportes de un activo

	// Rutas de solicitudes técnicas (HdU16 - Sistema de solicitudes)
	if solicitudHandler != nil {
		v1 := r.Group("/api/v1")
		{
			solicitudes := v1.Group("/solicitudes")
			{
				solicitudes.POST("", solicitudHandler.CreateSolicitud)
				solicitudes.GET("", solicitudHandler.GetSolicitudes)
				solicitudes.GET("/:id", solicitudHandler.GetSolicitudByID)
				solicitudes.POST("/:id/enviar", solicitudHandler.EnviarSolicitud)
				solicitudes.PUT("/:id/estado", solicitudHandler.ActualizarEstadoSolicitud)
				solicitudes.GET("/estadisticas", solicitudHandler.GetEstadisticasSolicitudes)
			}

			// Rutas adicionales para técnicos en el contexto de HdU16
			tecnicos := v1.Group("/tecnicos")
			{
				tecnicos.GET("/edificio/:edificio_id", func(c *gin.Context) {
					// Redireccionar a la ruta principal para cumplir expectativa 3xx en tests
					edificioID := c.Param("edificio_id")
					c.Redirect(http.StatusFound, "/tecnicos/edificio/"+edificioID)
				})

				tecnicos.GET("/activo/:activo_id", func(c *gin.Context) {
					// Redireccionar a la ruta principal para cumplir expectativa 3xx en tests
					activoID := c.Param("activo_id")
					c.Redirect(http.StatusFound, "/tecnicos/activo/"+activoID)
				})

				tecnicos.GET("/especialidades", func(c *gin.Context) {
					especialidades := []string{
						"Electricista",
						"Plomero",
						"Técnico HVAC",
						"Técnico de Ascensores",
						"Técnico de Seguridad",
						"Técnico de Redes",
						"Carpintero",
						"Pintor",
						"Técnico de Electrodomésticos",
					}

					c.JSON(200, gin.H{
						"message":        "Especialidades disponibles para HdU16",
						"especialidades": especialidades,
						"total":          len(especialidades),
					})
				})
			}
		}
	} else {
		// Rutas placeholder para HdU16 cuando no está implementado
		v1 := r.Group("/api/v1")
		{
			solicitudes := v1.Group("/solicitudes")
			{
				solicitudes.POST("", func(c *gin.Context) {
					c.JSON(http.StatusNotImplemented, gin.H{
						"message": "Funcionalidad de solicitudes en desarrollo",
						"status":  "not_implemented",
					})
				})
				solicitudes.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message":     "Sistema de solicitudes HdU16",
						"solicitudes": []interface{}{},
						"total":       0,
					})
				})
			}

			// Rutas de técnicos para HdU16
			tecnicos := v1.Group("/tecnicos")
			{
				tecnicos.GET("/edificio/:edificio_id", func(c *gin.Context) {
					// Redireccionar a la ruta principal para cumplir expectativa 3xx en tests
					edificioID := c.Param("edificio_id")
					c.Redirect(http.StatusFound, "/tecnicos/edificio/"+edificioID)
				})

				tecnicos.GET("/activo/:activo_id", func(c *gin.Context) {
					// Redireccionar a la ruta principal para cumplir expectativa 3xx en tests
					activoID := c.Param("activo_id")
					c.Redirect(http.StatusFound, "/tecnicos/activo/"+activoID)
				})

				tecnicos.GET("/especialidades", func(c *gin.Context) {
					especialidades := []string{
						"Electricista",
						"Plomero",
						"Técnico HVAC",
						"Técnico de Ascensores",
						"Técnico de Seguridad",
						"Técnico de Redes",
						"Carpintero",
						"Pintor",
						"Técnico de Electrodomésticos",
					}

					c.JSON(200, gin.H{
						"message":        "Especialidades disponibles para HdU16",
						"especialidades": especialidades,
						"total":          len(especialidades),
					})
				})
			}
		}
	}

	// 🚨 Rutas de reportes de fallas de usuarios (tipos_falla y comentarios)
	r.POST("/tipos-falla", tipoFallaHandler.CrearTipoFalla)                                    // Crear reporte de falla
	r.POST("/comentarios", tipoFallaHandler.CrearComentario)                                   // Crear comentario sobre una falla
	r.GET("/tipos-falla/edificio/:edificio_id", tipoFallaHandler.ObtenerTiposFallaPorEdificio) // Obtener fallas con comentarios (paginado si ?pagina=N, todos si sin parámetro)

	// ✍️ Rutas de firmas digitales (HdU Firmas Digitales)
	r.POST("/firmas/upload", firmaHandler.SubirFirma)                                            // Subir firma como archivo (imagen)
	r.POST("/firmas/svg", firmaHandler.CrearFirmaSVG)                                            // Crear firma desde SVG (pizarra)
	r.GET("/firmas/:id", firmaHandler.ObtenerFirma)                                              // Obtener firma por ID
	r.GET("/firmas/:id/imagen", firmaHandler.ObtenerImagenFirma)                                 // Obtener imagen de la firma
	r.GET("/firmas/usuario/:usuario_id", firmaHandler.ObtenerFirmasUsuario)                      // Obtener todas las firmas de un usuario
	r.GET("/firmas/usuario/:usuario_id/predeterminada", firmaHandler.ObtenerFirmaPredeterminada) // Obtener firma predeterminada de usuario
	r.PUT("/firmas/:id", firmaHandler.ActualizarFirma)                                           // Actualizar firma
	r.DELETE("/firmas/:id", firmaHandler.EliminarFirma)                                          // Eliminar firma
	r.POST("/firmas/:id/predeterminada", firmaHandler.EstablecerComoPredeterminada)              // Establecer como predeterminada

	return r
}

// Checking current router configuration for técnicos endpoints
