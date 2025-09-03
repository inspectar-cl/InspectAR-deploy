package handlers

import (
	"net/http"
	"notification/internal/models"
	"notification/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SensorNotificationHandler struct {
	service *services.SensorNotificationService
}

func NewSensorNotificationHandler(service *services.SensorNotificationService) *SensorNotificationHandler {
	return &SensorNotificationHandler{
		service: service,
	}
}

// CreateSensorAlert crea una nueva alerta de sensor
// POST /api/notifications/sensor-alert
func (h *SensorNotificationHandler) CreateSensorAlert(c *gin.Context) {
	var request models.SensorAlertRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar que el tipo esté dentro de los valores permitidos
	validTipos := map[string]bool{
		"sensor":        true,
		"alerta":        true,
		"mantenimiento": true,
		"sistema":       true,
	}

	if !validTipos[request.Tipo] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de notificación inválido. Valores permitidos: sensor, alerta, mantenimiento, sistema",
		})
		return
	}

	// Validar prioridad si se proporciona
	if request.Prioridad != "" {
		validPrioridades := map[string]bool{
			"low":      true,
			"medium":   true,
			"high":     true,
			"critical": true,
		}

		if !validPrioridades[request.Prioridad] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Prioridad inválida. Valores permitidos: low, medium, high, critical",
			})
			return
		}
	}

	notification, err := h.service.CreateSensorAlert(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al crear la notificación",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Notificación creada exitosamente",
		"notification": notification,
	})
}

// GetNotificationsByTipo obtiene notificaciones filtradas por tipo
// GET /api/notifications/by-tipo/:tipo?limit=50
func (h *SensorNotificationHandler) GetNotificationsByTipo(c *gin.Context) {
	tipo := c.Param("tipo")

	// Validar tipo
	validTipos := map[string]bool{
		"sensor":        true,
		"alerta":        true,
		"mantenimiento": true,
		"sistema":       true,
	}

	if !validTipos[tipo] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de notificación inválido. Valores permitidos: sensor, alerta, mantenimiento, sistema",
		})
		return
	}

	// Obtener límite de la query
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	notifications, err := h.service.GetNotificationsByTipo(tipo, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener notificaciones",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tipo":          tipo,
		"count":         len(notifications),
		"notifications": notifications,
	})
}

// GetNotificationsByPrioridad obtiene notificaciones filtradas por prioridad
// GET /api/notifications/by-prioridad/:prioridad?limit=50
func (h *SensorNotificationHandler) GetNotificationsByPrioridad(c *gin.Context) {
	prioridad := c.Param("prioridad")

	// Validar prioridad
	validPrioridades := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}

	if !validPrioridades[prioridad] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prioridad inválida. Valores permitidos: low, medium, high, critical",
		})
		return
	}

	// Obtener límite de la query
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	notifications, err := h.service.GetNotificationsByPrioridad(prioridad, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener notificaciones",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prioridad":     prioridad,
		"count":         len(notifications),
		"notifications": notifications,
	})
}

// GetAllNotifications obtiene todas las notificaciones
// GET /api/notifications?limit=100
func (h *SensorNotificationHandler) GetAllNotifications(c *gin.Context) {
	// Obtener límite de la query
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	notifications, err := h.service.GetAllNotifications(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener notificaciones",
			"details": err.Error(),
		})
		return
	}

	// Agrupar por tipo para estadísticas
	tipoCount := make(map[string]int)
	prioridadCount := make(map[string]int)

	for _, notification := range notifications {
		tipoCount[notification.Tipo]++
		prioridadCount[notification.Prioridad]++
	}

	c.JSON(http.StatusOK, gin.H{
		"total_count":   len(notifications),
		"notifications": notifications,
		"estadisticas": gin.H{
			"por_tipo":      tipoCount,
			"por_prioridad": prioridadCount,
		},
	})
}

// GetNotificationByID obtiene una notificación específica
// GET /api/notifications/:id
func (h *SensorNotificationHandler) GetNotificationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de notificación inválido",
		})
		return
	}

	notification, err := h.service.GetNotificationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Notificación no encontrada",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification": notification,
	})
}

// GetTiposDisponibles retorna los tipos de notificación disponibles
// GET /api/notifications/tipos
func (h *SensorNotificationHandler) GetTiposDisponibles(c *gin.Context) {
	tipos := []gin.H{
		{"codigo": "sensor", "nombre": "Sensor", "descripcion": "Alertas relacionadas con sensores", "icono": "📡"},
		{"codigo": "alerta", "nombre": "Alerta", "descripcion": "Alertas generales del sistema", "icono": "🚨"},
		{"codigo": "mantenimiento", "nombre": "Mantenimiento", "descripcion": "Notificaciones de mantenimiento", "icono": "🔧"},
		{"codigo": "sistema", "nombre": "Sistema", "descripcion": "Notificaciones del sistema", "icono": "⚙️"},
	}

	prioridades := []gin.H{
		{"codigo": "low", "nombre": "Baja", "valor": 1, "color": "#28a745", "icono": "🟢"},
		{"codigo": "medium", "nombre": "Media", "valor": 2, "color": "#ffc107", "icono": "🟡"},
		{"codigo": "high", "nombre": "Alta", "valor": 3, "color": "#fd7e14", "icono": "🟠"},
		{"codigo": "critical", "nombre": "Crítica", "valor": 4, "color": "#dc3545", "icono": "🔴"},
	}

	c.JSON(http.StatusOK, gin.H{
		"tipos":       tipos,
		"prioridades": prioridades,
	})
}
