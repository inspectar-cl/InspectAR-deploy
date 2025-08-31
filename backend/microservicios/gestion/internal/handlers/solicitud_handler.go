package handlers

import (
	"gestion/internal/models"
	"gestion/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SolicitudHandler struct {
	service *services.SolicitudService
}

func NewSolicitudHandler(service *services.SolicitudService) *SolicitudHandler {
	return &SolicitudHandler{service: service}
}

// CreateSolicitud crea una nueva solicitud técnica
// @Summary Crear una nueva solicitud técnica
// @Description Crea una nueva solicitud de servicio técnico especializado
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param solicitud body models.CreateSolicitudRequest true "Datos de la solicitud"
// @Success 201 {object} models.SolicitudResponse "Solicitud creada exitosamente"
// @Failure 400 {object} gin.H "Error en los datos de entrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes [post]
func (h *SolicitudHandler) CreateSolicitud(c *gin.Context) {
	var req models.CreateSolicitudRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	solicitud, err := h.service.CreateSolicitud(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al crear la solicitud",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Solicitud creada exitosamente",
		"solicitud": solicitud,
	})
}

// GetSolicitudes obtiene solicitudes con filtros opcionales
// @Summary Obtener solicitudes con filtros
// @Description Lista solicitudes con filtros opcionales como técnico, estado, prioridad, etc.
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param tecnico_id query int false "ID del técnico"
// @Param residente_id query int false "ID del residente"
// @Param activo_id query int false "ID del activo"
// @Param edificio_id query int false "ID del edificio"
// @Param estado query string false "Estado de la solicitud"
// @Param tipo query string false "Tipo de solicitud"
// @Param prioridad query string false "Prioridad de la solicitud"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Límite de resultados por página" default(10)
// @Success 200 {object} models.SolicitudListResponse "Lista de solicitudes"
// @Failure 400 {object} gin.H "Error en los parámetros"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes [get]
func (h *SolicitudHandler) GetSolicitudes(c *gin.Context) {
	// Extraer filtros de los query parameters
	filter := models.SolicitudFilter{}

	if tecnicoID := c.Query("tecnico_id"); tecnicoID != "" {
		if id, err := strconv.Atoi(tecnicoID); err == nil {
			filter.TecnicoID = &id
		}
	}

	if residenteID := c.Query("residente_id"); residenteID != "" {
		if id, err := strconv.Atoi(residenteID); err == nil {
			filter.ResidenteID = &id
		}
	}

	if activoID := c.Query("activo_id"); activoID != "" {
		if id, err := strconv.Atoi(activoID); err == nil {
			filter.ActivoID = &id
		}
	}

	if edificioID := c.Query("edificio_id"); edificioID != "" {
		if id, err := strconv.Atoi(edificioID); err == nil {
			filter.EdificioID = &id
		}
	}

	if estado := c.Query("estado"); estado != "" {
		estadoEnum := models.EstadoSolicitud(estado)
		filter.Estado = &estadoEnum
	}

	if tipo := c.Query("tipo"); tipo != "" {
		tipoEnum := models.TipoSolicitud(tipo)
		filter.Tipo = &tipoEnum
	}

	if prioridad := c.Query("prioridad"); prioridad != "" {
		prioridadEnum := models.PrioridadSolicitud(prioridad)
		filter.Prioridad = &prioridadEnum
	}

	if busqueda := c.Query("busqueda"); busqueda != "" {
		filter.Busqueda = &busqueda
	}

	// Parámetros de paginación
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Obtener solicitudes
	response, err := h.solicitudService.CreateSolicitud(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener solicitudes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSolicitudByID obtiene una solicitud específica por ID
// @Summary Obtener solicitud por ID
// @Description Obtiene los detalles completos de una solicitud específica
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param id path int true "ID de la solicitud"
// @Success 200 {object} models.SolicitudResponse "Detalles de la solicitud"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Solicitud no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes/{id} [get]
func (h *SolicitudHandler) GetSolicitudByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de solicitud inválido",
		})
		return
	}

	solicitud, err := h.service.GetSolicitudByID(id)
	if err != nil {
		if err.Error() == "solicitud no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Solicitud no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener la solicitud",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Solicitud obtenida exitosamente",
		"solicitud": solicitud,
	})
}

// EnviarSolicitud envía una solicitud al técnico seleccionado
// @Summary Enviar solicitud al técnico
// @Description Envía la solicitud al técnico especializado y actualiza su estado
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param id path int true "ID de la solicitud"
// @Success 200 {object} gin.H "Solicitud enviada exitosamente"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Solicitud no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes/{id}/enviar [post]
func (h *SolicitudHandler) EnviarSolicitud(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de solicitud inválido",
		})
		return
	}

	err = h.service.EnviarSolicitud(id)
	if err != nil {
		if err.Error() == "solicitud no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Solicitud no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al enviar la solicitud",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Solicitud enviada exitosamente",
	})
}

// ActualizarEstadoSolicitud actualiza el estado de una solicitud
// @Summary Actualizar estado de solicitud
// @Description Permite al técnico actualizar el estado y respuesta de una solicitud
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param id path int true "ID de la solicitud"
// @Param update body models.UpdateSolicitudRequest true "Datos de actualización"
// @Success 200 {object} gin.H "Estado actualizado exitosamente"
// @Failure 400 {object} gin.H "Datos inválidos"
// @Failure 404 {object} gin.H "Solicitud no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes/{id}/estado [put]
func (h *SolicitudHandler) ActualizarEstadoSolicitud(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de solicitud inválido",
		})
		return
	}

	var req models.UpdateSolicitudRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar que al menos el estado esté presente
	if req.Estado == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El estado es requerido",
		})
		return
	}

	// Obtener respuesta del técnico si está presente
	respuesta := ""
	if req.RespuestaTecnico != nil {
		respuesta = *req.RespuestaTecnico
	}

	err = h.service.ActualizarEstadoSolicitud(id, *req.Estado, respuesta)
	if err != nil {
		if err.Error() == "solicitud no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Solicitud no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar el estado",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado actualizado exitosamente",
	})
}

// GetEstadisticasSolicitudes obtiene estadísticas de solicitudes
// @Summary Obtener estadísticas de solicitudes
// @Description Obtiene estadísticas agregadas de solicitudes por edificio o técnico
// @Tags Solicitudes
// @Accept json
// @Produce json
// @Param edificio_id query int false "ID del edificio para filtrar estadísticas"
// @Param tecnico_id query int false "ID del técnico para filtrar estadísticas"
// @Success 200 {object} models.EstadisticasSolicitudes "Estadísticas de solicitudes"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /api/v1/solicitudes/estadisticas [get]
func (h *SolicitudHandler) GetEstadisticasSolicitudes(c *gin.Context) {
	var edificioID *int
	var tecnicoID *int

	if eID := c.Query("edificio_id"); eID != "" {
		if id, err := strconv.Atoi(eID); err == nil {
			edificioID = &id
		}
	}

	if tID := c.Query("tecnico_id"); tID != "" {
		if id, err := strconv.Atoi(tID); err == nil {
			tecnicoID = &id
		}
	}

	estadisticas, err := h.service.GetEstadisticasSolicitudes(edificioID, tecnicoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener estadísticas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, estadisticas)
}
