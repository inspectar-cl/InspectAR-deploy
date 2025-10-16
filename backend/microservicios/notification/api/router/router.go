package router

import (
	handlers "notification/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(notificationHandler *handlers.NotificationHandler) *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Rutas para notificaciones
	r.POST("/notification", notificationHandler.CreateNotification)
	r.POST("/sensor/alert", notificationHandler.CreateSensorAlert) // Nueva ruta para alertas de sensor
	r.GET("/notification/:activo_id", notificationHandler.GetNotificationsByActivoID)
	r.PUT("/notification/:notification_id/send", notificationHandler.SendNotification)
	r.GET("/tipos-notificacion", notificationHandler.GetTiposNotificacion) // Nueva ruta para tipos de notificación

	// Rutas para comunicación con técnicos
	r.POST("/technician/contact", notificationHandler.SendTechnicianContact) // Nueva ruta para contacto técnico

	// Rutas para edificios
	r.GET("/edificios", notificationHandler.GetAllEdificios)
	r.GET("/edificio/:id", notificationHandler.GetEdificioByID)
	r.GET("/edificio/:id/usuarios", notificationHandler.GetUsuariosByEdificioID)
	r.GET("/edificio/:id/activos", notificationHandler.GetActivosByEdificioID)

	return r
}
