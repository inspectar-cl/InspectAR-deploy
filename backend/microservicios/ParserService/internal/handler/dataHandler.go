package handlers

import (
	"net/http"
	"time"

	"ParserService/internal/models"
	"ParserService/internal/services"

	"github.com/gin-gonic/gin"
)

type DataHandler struct {
	activoService     *services.ActivoService
	sensorService     *services.SensorService
	monitoringService *services.SensorMonitoringService
}

func NewDataHandler(activoSvc *services.ActivoService, sensorSvc *services.SensorService, monitoringSvc *services.SensorMonitoringService) *DataHandler {
	return &DataHandler{
		activoService:     activoSvc,
		sensorService:     sensorSvc,
		monitoringService: monitoringSvc,
	}
}

// POST /activo
func (h *DataHandler) CreateActivo(c *gin.Context) {
	var activo models.Activo
	if err := c.ShouldBindJSON(&activo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	id, err := h.activoService.CrearActivo(c.Request.Context(), &activo)
	if err != nil {
		if err.Error() == "el activo ya existe" {
			c.JSON(http.StatusConflict, gin.H{"error": "El activo ya existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activo_id": id})
}

// POST /lectura
func (h *DataHandler) CreateLectura(c *gin.Context) {
	var data models.LecturaSensorRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de lectura inválido"})
		return
	}

	timestamp, err := time.Parse(time.RFC3339, data.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp inválido"})
		return
	}

	// Insertar lectura en InfluxDB
	err = h.sensorService.InsertarLectura(c.Request.Context(), data.SensorID, data.Valor, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar la lectura en InfluxDB"})
		return
	}

	// Actualizar estado del sensor en MongoDB
	err = h.monitoringService.UpdateSensorActivity(data.SensorID, timestamp)
	if err != nil {
		// Log error pero no fallar la request principal
		c.Header("X-Warning", "Error al actualizar estado del sensor")
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Lectura registrada en InfluxDB"})
}

// GET /activo/:activo_id
func (h *DataHandler) GetActivo(c *gin.Context) {
	activoID := c.Param("activo_id")

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, activo)
}

// GET /activo/:activo_id/datos
func (h *DataHandler) GetSensorByActivo(c *gin.Context) {
	activoID := c.Param("activo_id")

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	var allLecturas []models.SensorDataset

	for _, sensor := range activo.Sensores {
		datos, err := h.sensorService.GetDatosSensor(c.Request.Context(), sensor.SensorID, 30*time.Minute)
		if err != nil {
			continue
		}
		allLecturas = append(allLecturas, models.SensorDataset{
			SensorID: sensor.SensorID,
			Datos:    datos,
		})
	}

	respuesta := gin.H{
		"id":        activo.ID,
		"activo_id": activo.ActivoID,
		"nombre":    activo.Nombre,
		"ubicacion": activo.Ubicacion,
		"estado":    activo.Estado,
		"sensores":  allLecturas,
	}

	c.JSON(http.StatusOK, respuesta)
}

// GET /activo
func (h *DataHandler) GetAllActivos(c *gin.Context) {
	activos, err := h.activoService.GetAllActivos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los activos"})
		return
	}
	c.JSON(http.StatusOK, activos)
}

// GET /sensor/:sensor_id/last
func (h *DataHandler) GetSensorLastData(c *gin.Context) {
	activoID := c.Param("activo_id")
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	result := make(map[string]interface{})
	for _, sensor := range activo.Sensores {
		dato, err := h.sensorService.GetSensorLastData(c.Request.Context(), sensor.SensorID)
		if err != nil || dato == nil {
			result[sensor.SensorID] = nil
		} else {
			result[sensor.SensorID] = dato
		}
	}

	c.JSON(http.StatusOK, result)
}

// PUT /activo/:activo_id/estado
func (h *DataHandler) UpdActivoEstado(c *gin.Context) {
	activoID := c.Param("activo_id")

	// Leer el nuevo estado desde el cuerpo del request
	var body struct {
		Estado string `json:"estado"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Estado == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido o faltante"})
		return
	}

	// Actualizar en la base de datos
	err := h.activoService.ActualizarEstado(c.Request.Context(), activoID, body.Estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado actualizado correctamente"})
}

// GET /activo/:activo_id/sensores/estado - Obtiene un activo con el estado de sus sensores
func (h *DataHandler) GetActivoWithSensorStatus(c *gin.Context) {
	activoID := c.Param("activo_id")

	// Obtener el activo
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Crear la respuesta con información del activo y estado de sensores
	response := gin.H{
		"id":          activo.ID,
		"activo_id":   activo.ActivoID,
		"nombre":      activo.Nombre,
		"ubicacion":   activo.Ubicacion,
		"estado":      activo.Estado,
		"id_edificio": activo.Id_edificio,
		"sensores":    []gin.H{},
	}

	// Obtener estado de cada sensor
	var sensoresConEstado []gin.H
	for _, sensor := range activo.Sensores {
		sensorInfo := gin.H{
			"sensor_id":     sensor.SensorID,
			"tipo":          sensor.Tipo,
			"unidad":        sensor.Unidad,
			"estado":        "unknown", // Default
			"is_active":     false,     // Default
			"last_seen":     nil,       // Default
			"total_reports": 0,         // Default
		}

		// Obtener estado del sensor desde el servicio de monitoreo
		sensorStatus, err := h.monitoringService.GetSensorStatus(sensor.SensorID)
		if err == nil && sensorStatus != nil {
			var estado string
			if sensorStatus.IsActive {
				estado = "connected"
			} else {
				estado = "disconnected"
			}

			sensorInfo["estado"] = estado
			sensorInfo["is_active"] = sensorStatus.IsActive
			sensorInfo["last_seen"] = sensorStatus.LastSeen
			sensorInfo["first_seen"] = sensorStatus.FirstSeen
			sensorInfo["total_reports"] = sensorStatus.TotalReports
			sensorInfo["created_at"] = sensorStatus.CreatedAt
			sensorInfo["updated_at"] = sensorStatus.UpdatedAt
		} else {
			// Si no hay información de estado, el sensor nunca ha enviado datos
			sensorInfo["estado"] = "never_connected"
		}

		sensoresConEstado = append(sensoresConEstado, sensorInfo)
	}

	response["sensores"] = sensoresConEstado
	response["total_sensores"] = len(activo.Sensores)

	// Agregar estadísticas resumidas
	activeSensors := 0
	disconnectedSensors := 0
	neverConnectedSensors := 0

	for _, sensor := range sensoresConEstado {
		estado := sensor["estado"].(string)
		switch estado {
		case "connected":
			activeSensors++
		case "disconnected":
			disconnectedSensors++
		case "never_connected":
			neverConnectedSensors++
		}
	}

	response["resumen"] = gin.H{
		"sensores_activos":          activeSensors,
		"sensores_desconectados":    disconnectedSensors,
		"sensores_nunca_conectados": neverConnectedSensors,
	}

	c.JSON(http.StatusOK, response)
}

// POST /activo/:activo_id/sensores - Agregar sensor a un activo existente
func (h *DataHandler) AddSensorToActivo(c *gin.Context) {
	activoID := c.Param("activo_id")

	var sensor models.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del sensor inválidos"})
		return
	}

	// Validar campos requeridos
	if sensor.SensorID == "" || sensor.Tipo == "" || sensor.Unidad == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los campos sensor_id, tipo y unidad son requeridos"})
		return
	}

	// Validar que el activo existe
	_, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Agregar el sensor al activo
	err = h.activoService.AgregarSensor(c.Request.Context(), activoID, sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo agregar el sensor al activo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":   "Sensor agregado correctamente al activo",
		"activo_id": activoID,
		"sensor":    sensor,
	})
}
