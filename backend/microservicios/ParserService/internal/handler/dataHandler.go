package handlers

import (
	"log"
	"net/http"
	"time"

	"ParserService/internal/models"
	"ParserService/internal/services"

	"strconv"

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

// GET /activo/:activo_id - Obtener activo con información de sensores
func (h *DataHandler) GetActivo(c *gin.Context) {
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Obtener sensores del activo
	sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activoID)
	if err != nil {
		// Si no se pueden obtener sensores, retornar activo sin información de sensores
		c.JSON(http.StatusOK, gin.H{
			"id":             activo.ID,
			"activo_id":      activo.ActivoID,
			"estado":         activo.Estado,
			"id_edificio":    activo.EdificioID,
			"sensores":       []gin.H{},
			"total_sensores": 0,
		})
		return
	}

	// Obtener estado de cada sensor
	var sensoresConEstado []gin.H
	for _, sensor := range sensores {
		sensorInfo := gin.H{
			"sensor_id":     sensor.SensorID,
			"tipo":          sensor.Tipo,
			"unidad":        sensor.Unidad,
			"estado":        "unknown",
			"is_active":     false,
			"last_seen":     nil,
			"total_reports": 0,
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
			sensorInfo["estado"] = "never_connected"
		}

		sensoresConEstado = append(sensoresConEstado, sensorInfo)
	}

	// Retornar activo con información enriquecida de sensores
	c.JSON(http.StatusOK, gin.H{
		"id":             activo.ID,
		"activo_id":      activo.ActivoID,
		"estado":         activo.Estado,
		"id_edificio":    activo.EdificioID,
		"sensores":       sensoresConEstado,
		"total_sensores": len(sensores),
	})
}

// GET /activo/:activo_id/datos
func (h *DataHandler) GetSensorByActivo(c *gin.Context) {
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Obtener sensores desde el servicio de activos
	sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activoID)
	if err != nil {
		log.Printf("❌ ERROR: No se pudieron obtener los sensores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los sensores"})
		return
	}
	log.Printf("✅ Sensores obtenidos exitosamente: %d sensores", len(sensores))

	var allLecturas []models.SensorDataset

	log.Printf("🔄 Iniciando loop para %d sensores", len(sensores))

	for _, sensor := range sensores {
		log.Printf("📡 Procesando sensor: %s", sensor.SensorID)
		datos, err := h.sensorService.GetDatosSensor(c.Request.Context(), sensor.SensorID, 30*time.Minute)
		if err != nil {
			log.Printf("❌ Error obteniendo datos para sensor %s: %v", sensor.SensorID, err)
			continue
		}
		log.Printf("✅ Datos recibidos para sensor %s: %d registros", sensor.SensorID, len(datos))
		allLecturas = append(allLecturas, models.SensorDataset{
			SensorID: sensor.SensorID,
			Datos:    datos,
		})
	}

	respuesta := gin.H{
		"activo_id":   activo.ActivoID,
		"estado":      activo.Estado,
		"edificio_id": activo.EdificioID,
		"sensores":    allLecturas,
	}

	c.JSON(http.StatusOK, respuesta)
}

// GET /lectura/:activo_id/window - Obtiene datos paginados por sensor
// Parámetros query: page (default: 1), limit (default: 1000)
func (h *DataHandler) GetSensorByActivoWindow(c *gin.Context) {
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

	// Leer parámetros de paginación
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "1000")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page debe ser un entero mayor a 0"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit debe ser un entero entre 1 y 5000"})
		return
	}

	// Calcular offset: (page - 1) * limit
	offset := (page - 1) * limit

	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Obtener sensores desde el servicio de activos
	sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activoID)
	if err != nil {
		log.Printf("❌ ERROR: No se pudieron obtener los sensores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los sensores"})
		return
	}
	log.Printf("✅ Sensores obtenidos exitosamente para window: %d sensores (page: %d, limit: %d)", len(sensores), page, limit)

	var allLecturas []models.SensorDataset

	log.Printf("🔄 Iniciando loop window para %d sensores (página: %d, límite: %d, offset: %d)", len(sensores), page, limit, offset)

	for _, sensor := range sensores {
		log.Printf("📡 Procesando sensor window: %s", sensor.SensorID)
		datos, err := h.sensorService.GetDatosSensorWindow(c.Request.Context(), sensor.SensorID, limit, offset)
		if err != nil {
			log.Printf("❌ Error obteniendo datos window para sensor %s: %v", sensor.SensorID, err)
			continue
		}
		log.Printf("✅ Datos window recibidos para sensor %s: %d registros", sensor.SensorID, len(datos))
		allLecturas = append(allLecturas, models.SensorDataset{
			SensorID: sensor.SensorID,
			Datos:    datos,
		})
	}

	respuesta := gin.H{
		"activo_id":   activo.ActivoID,
		"estado":      activo.Estado,
		"edificio_id": activo.EdificioID,
		"page":        page,
		"limit":       limit,
		"offset":      offset,
		"sensores":    allLecturas,
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
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Obtener sensores del activo
	sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los sensores"})
		return
	}

	result := make(map[string]interface{})
	for _, sensor := range sensores {
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
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

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
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

	// Obtener el activo
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Obtener sensores del activo
	sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los sensores"})
		return
	}

	// Crear la respuesta con información del activo y estado de sensores
	response := gin.H{
		"activo_id":   activo.ActivoID,
		"estado":      activo.Estado,
		"edificio_id": activo.EdificioID,
		"sensores":    []gin.H{},
	}

	// Obtener estado de cada sensor
	var sensoresConEstado []gin.H
	for _, sensor := range sensores {
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
	response["total_sensores"] = len(sensores)

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
	idStr := c.Param("activo_id")
	activoID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activo_id debe ser entero"})
		return
	}

	var sensor models.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		// Log para diagnóstico
		c.Error(err)
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

// GET /activo/edificio/:edificio_id - Obtener todos los activos de un edificio con información de sensores
func (h *DataHandler) GetActivosByEdificio(c *gin.Context) {
	idStr := c.Param("edificio_id")
	edificioID, errConv := strconv.Atoi(idStr)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "edificio_id debe ser entero"})
		return
	}

	activos, err := h.activoService.ObtenerActivosPorEdificio(c.Request.Context(), edificioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los activos del edificio"})
		return
	}

	// Enriquecer cada activo con el estado de sus sensores
	var activosEnriquecidos []gin.H
	for _, activo := range activos {
		// Obtener sensores del activo
		sensores, err := h.activoService.ObtenerSensoresPorActivo(c.Request.Context(), activo.ActivoID)
		if err != nil {
			// Si no se pueden obtener sensores, incluir el activo sin información de sensores
			activosEnriquecidos = append(activosEnriquecidos, gin.H{
				"id":             activo.ID,
				"activo_id":      activo.ActivoID,
				"estado":         activo.Estado,
				"id_edificio":    activo.EdificioID,
				"sensores":       []gin.H{},
				"total_sensores": 0,
			})
			continue
		}

		// Obtener estado de cada sensor
		var sensoresConEstado []gin.H
		for _, sensor := range sensores {
			sensorInfo := gin.H{
				"sensor_id":     sensor.SensorID,
				"tipo":          sensor.Tipo,
				"unidad":        sensor.Unidad,
				"estado":        "unknown",
				"is_active":     false,
				"last_seen":     nil,
				"total_reports": 0,
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
				sensorInfo["estado"] = "never_connected"
			}

			sensoresConEstado = append(sensoresConEstado, sensorInfo)
		}

		activosEnriquecidos = append(activosEnriquecidos, gin.H{
			"id":             activo.ID,
			"activo_id":      activo.ActivoID,
			"estado":         activo.Estado,
			"id_edificio":    activo.EdificioID,
			"sensores":       sensoresConEstado,
			"total_sensores": len(sensores),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"edificio_id": edificioID,
		"activos":     activosEnriquecidos,
		"total":       len(activosEnriquecidos),
	})
}

// ==================== MÉTODOS DE ADMINISTRACIÓN ====================

// DELETE /admin/activos/:activo_id - Eliminar activo
func (h *DataHandler) DeleteActivo(c *gin.Context) {
	activoIDStr := c.Param("activo_id")
	activoID, err := strconv.Atoi(activoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de activo inválido"})
		return
	}

	// Verificar que el activo existe
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), activoID)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Eliminar el activo
	err = h.activoService.EliminarActivo(c.Request.Context(), activoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el activo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Activo eliminado exitosamente",
		"activo_id": activoID,
	})
}

// POST /admin/sensores - Crear sensor (desde gestion-service)
func (h *DataHandler) AddSensorToActivoAdmin(c *gin.Context) {
	var req struct {
		IDActivo int    `json:"id_activo" binding:"required"`
		Nombre   string `json:"nombre" binding:"required"`
		Tipo     string `json:"tipo" binding:"required"`
		Unidad   string `json:"unidad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Verificar que el activo existe
	activo, err := h.activoService.ObtenerActivo(c.Request.Context(), req.IDActivo)
	if err != nil || activo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado en ParserService"})
		return
	}

	// Crear el sensor
	sensor := models.Sensor{
		SensorID: req.Nombre, // Usar nombre como ID inicial
		Tipo:     req.Tipo,
		Unidad:   req.Unidad,
	}

	err = h.activoService.AgregarSensor(c.Request.Context(), req.IDActivo, sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el sensor", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Sensor creado exitosamente",
		"sensor_id": sensor.SensorID,
		"activo_id": req.IDActivo,
		"sensor":    sensor,
	})
}

// PUT /admin/sensores/:sensor_id - Actualizar sensor
func (h *DataHandler) UpdateSensor(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	if sensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sensor requerido"})
		return
	}

	var req struct {
		Nombre string `json:"nombre"`
		Tipo   string `json:"tipo"`
		Unidad string `json:"unidad"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Actualizar el sensor
	err := h.sensorService.ActualizarSensor(c.Request.Context(), sensorID, req.Nombre, req.Tipo, req.Unidad)
	if err != nil {
		if err.Error() == "sensor no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sensor no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el sensor", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sensor actualizado exitosamente",
		"sensor_id": sensorID,
	})
}

// DELETE /admin/sensores/:sensor_id - Eliminar sensor
func (h *DataHandler) DeleteSensor(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	if sensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sensor requerido"})
		return
	}

	// Eliminar el sensor del activo
	err := h.sensorService.EliminarSensor(c.Request.Context(), sensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el sensor", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sensor eliminado exitosamente",
		"sensor_id": sensorID,
	})
}
