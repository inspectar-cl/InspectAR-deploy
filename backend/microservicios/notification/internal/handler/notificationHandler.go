package handlers

import (
	"net/http"
	"notification/internal/models"
	"notification/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationSvc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationSvc,
	}
}

// POST /sensor/alert
func (h *NotificationHandler) CreateSensorAlert(c *gin.Context) {
	var request models.SensorAlertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	notificacion, err := h.notificationService.CreateSensorAlert(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear alerta de sensor: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Alerta de sensor creada y enviada exitosamente",
		"notificacion": notificacion,
	})
}

// POST /notification
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var request models.NotificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	notificacion, err := h.notificationService.CreateNotification(request.ActivoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear notificación: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Notificación creada exitosamente",
		"notificacion": notificacion,
	})
}

// GET /notification/:activo_id
func (h *NotificationHandler) GetNotificationsByActivoID(c *gin.Context) {
	activoIDStr := c.Param("activo_id")
	activoID, err := strconv.ParseUint(activoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	notifications, err := h.notificationService.GetNotificationsByActivoID(uint(activoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notificaciones": notifications,
	})
}

// PUT /notification/:notification_id/send
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	notificationIDStr := c.Param("notification_id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación inválido"})
		return
	}

	err = h.notificationService.SendNotification(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar notificación: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notificación enviada exitosamente",
	})
}

// GET /tipos-notificacion
func (h *NotificationHandler) GetTiposNotificacion(c *gin.Context) {
	tipos, err := h.notificationService.GetAllTiposNotificacion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tipos de notificación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tipos_notificacion": tipos,
	})
}

// GET /edificios
func (h *NotificationHandler) GetAllEdificios(c *gin.Context) {
	edificios, err := h.notificationService.GetAllEdificios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener edificios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edificios": edificios,
	})
}

// GET /edificio/:id
func (h *NotificationHandler) GetEdificioByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	edificio, err := h.notificationService.GetEdificioByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener edificio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edificio": edificio,
	})
}

// GET /edificio/:id/usuarios
func (h *NotificationHandler) GetUsuariosByEdificioID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	usuarios, err := h.notificationService.GetUsuariosByEdificioID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuarios": usuarios,
	})
}

// GET /edificio/:id/activos
func (h *NotificationHandler) GetActivosByEdificioID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	activos, err := h.notificationService.GetActivosByEdificioID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener activos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activos": activos,
	})
}

// POST /technician/contact
func (h *NotificationHandler) SendTechnicianContact(c *gin.Context) {
	var request models.TechnicianContactRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	err := h.notificationService.SendTechnicianContactRequest(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al enviar solicitud de contacto técnico: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Solicitud de contacto técnico enviada exitosamente",
		"technician_email":  request.TechnicianEmail,
		"user_email":        request.UserEmail,
		"activo_id":         request.ActivoID,
		"priority":          request.Priority,
	})
}
