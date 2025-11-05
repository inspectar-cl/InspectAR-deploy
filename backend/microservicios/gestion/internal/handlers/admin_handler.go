package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
	"gestion/internal/services"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	activoRepo   *repository.ActivoRepository
	edificioRepo *repository.EdificioRepository
	usuarioRepo  *repository.UsuarioRepository
	logRepo      *repository.LogRepository
	qrService    *services.QRService
	parserURL    string // URL del microservicio ParserService
}

func NewAdminHandler(
	activoRepo *repository.ActivoRepository,
	edificioRepo *repository.EdificioRepository,
	usuarioRepo *repository.UsuarioRepository,
	logRepo *repository.LogRepository,
	qrService *services.QRService,
	parserURL string,
) *AdminHandler {
	return &AdminHandler{
		activoRepo:   activoRepo,
		edificioRepo: edificioRepo,
		usuarioRepo:  usuarioRepo,
		logRepo:      logRepo,
		qrService:    qrService,
		parserURL:    parserURL,
	}
}

// ==================== ACTIVOS ====================

// CrearActivo crea un activo en gestion y en ParserService
func (h *AdminHandler) CrearActivo(c *gin.Context) {
	var req models.CreateActivoAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar tipo de activo
	if !models.ValidarTipoActivo(req.Tipo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de activo inválido"})
		return
	}

	// Crear activo en gestion
	activo := models.Activo{
		Nombre:      req.Nombre,
		Tipo:        req.Tipo,
		Descripcion: req.Descripcion,
		Ubicacion:   req.Ubicacion,
		EdificioID:  req.EdificioID,
	}

	activoCreado, err := h.activoRepo.Crear(activo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear activo", "details": err.Error()})
		return
	}

	// Crear activo en ParserService
	parserReq := map[string]interface{}{
		"id_activo": activoCreado.ID,
		"nombre":    activoCreado.Nombre,
		"tipo":      activoCreado.Tipo,
	}

	jsonData, _ := json.Marshal(parserReq)
	resp, err := http.Post(h.parserURL+"/admin/activos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Rollback: eliminar de gestion si falla en parser
		h.activoRepo.Eliminar(activoCreado.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear activo en ParserService", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Rollback
		h.activoRepo.Eliminar(activoCreado.ID)
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ParserService rechazó la creación", "details": string(bodyBytes)})
		return
	}

	// Generar código QR para el activo
	if activoCreado.CodigoActivo != nil && *activoCreado.CodigoActivo != "" {
		qrBase64, urlQR, err := h.qrService.GenerarQRActivo(activoCreado.ID, *activoCreado.CodigoActivo)
		if err != nil {
			log.Printf("⚠️ Error generando QR para activo ID %d: %v", activoCreado.ID, err)
			// No falla la creación si QR falla, solo lo registra
		} else {
			// Guardar QR en base de datos
			if err := h.activoRepo.ActualizarQR(activoCreado.ID, qrBase64, urlQR); err != nil {
				log.Printf("⚠️ Error guardando QR en BD: %v", err)
			} else {
				activoCreado.CodigoQR = &qrBase64
				activoCreado.URLQR = &urlQR
				log.Printf("✓ QR generado exitosamente para activo ID %d", activoCreado.ID)
			}
		}
	}

	// Registrar en log de auditoría
	datosNuevos := map[string]interface{}{
		"id":          activoCreado.ID,
		"nombre":      activoCreado.Nombre,
		"tipo":        activoCreado.Tipo,
		"descripcion": activoCreado.Descripcion,
		"ubicacion":   activoCreado.Ubicacion,
		"edificio_id": activoCreado.EdificioID,
		"codigo":      activoCreado.CodigoActivo,
	}

	logAudit := models.LogAuditoria{
		UsuarioEmail: req.Email,
		Accion:       "crear",
		Entidad:      "activo",
		EntidadID:    activoCreado.ID,
		DatosNuevos:  datosNuevos,
		Descripcion:  fmt.Sprintf("Usuario %s creó el activo '%s' (ID: %d, Código: %s)", req.Email, activoCreado.Nombre, activoCreado.ID, *activoCreado.CodigoActivo),
		IPOrigen:     stringPtr(c.ClientIP()),
		UserAgent:    stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(logAudit)

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Activo creado exitosamente en ambos servicios",
		"activo":        activoCreado,
		"codigo":        activoCreado.CodigoActivo,
		"url_qr":        activoCreado.URLQR,
		"qr_disponible": activoCreado.CodigoQR != nil && *activoCreado.CodigoQR != "",
	})
}

// ActualizarActivo actualiza un activo
func (h *AdminHandler) ActualizarActivo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateActivoAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener activo actual
	activoAnterior, err := h.activoRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Preparar datos actualizados
	activoActualizado := *activoAnterior
	if req.Nombre != nil {
		activoActualizado.Nombre = *req.Nombre
	}
	if req.Tipo != nil {
		if !models.ValidarTipoActivo(*req.Tipo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de activo inválido"})
			return
		}
		activoActualizado.Tipo = *req.Tipo
	}
	if req.Descripcion != nil {
		activoActualizado.Descripcion = req.Descripcion
	}
	if req.Ubicacion != nil {
		activoActualizado.Ubicacion = *req.Ubicacion
	}
	if req.EdificioID != nil {
		activoActualizado.EdificioID = *req.EdificioID
	}

	// Actualizar en gestion
	err = h.activoRepo.Actualizar(id, activoActualizado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar activo", "details": err.Error()})
		return
	}

	// Registrar en log
	datosAnteriores := map[string]interface{}{
		"id":          activoAnterior.ID,
		"nombre":      activoAnterior.Nombre,
		"tipo":        activoAnterior.Tipo,
		"descripcion": activoAnterior.Descripcion,
		"ubicacion":   activoAnterior.Ubicacion,
		"edificio_id": activoAnterior.EdificioID,
	}

	datosNuevos := map[string]interface{}{
		"id":          activoActualizado.ID,
		"nombre":      activoActualizado.Nombre,
		"tipo":        activoActualizado.Tipo,
		"descripcion": activoActualizado.Descripcion,
		"ubicacion":   activoActualizado.Ubicacion,
		"edificio_id": activoActualizado.EdificioID,
	}

	log := models.LogAuditoria{
		UsuarioEmail:    req.Email,
		Accion:          "modificar",
		Entidad:         "activo",
		EntidadID:       id,
		DatosAnteriores: datosAnteriores,
		DatosNuevos:     datosNuevos,
		Descripcion:     fmt.Sprintf("Usuario %s modificó el activo '%s' (ID: %d)", req.Email, activoActualizado.Nombre, id),
		IPOrigen:        stringPtr(c.ClientIP()),
		UserAgent:       stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusOK, gin.H{
		"message": "Activo actualizado exitosamente",
		"activo":  activoActualizado,
	})
}

// EliminarActivo elimina un activo
func (h *AdminHandler) EliminarActivo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email requerido"})
		return
	}

	// Obtener activo antes de eliminar
	activo, err := h.activoRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Eliminar de gestion
	err = h.activoRepo.Eliminar(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar activo", "details": err.Error()})
		return
	}

	// Eliminar de ParserService
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/admin/activos/%d", h.parserURL, id), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Log del error pero no revertir
		fmt.Printf("Error al eliminar activo de ParserService: %v\n", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	// Registrar en log
	datosAnteriores := map[string]interface{}{
		"id":          activo.ID,
		"nombre":      activo.Nombre,
		"tipo":        activo.Tipo,
		"descripcion": activo.Descripcion,
		"ubicacion":   activo.Ubicacion,
		"edificio_id": activo.EdificioID,
	}

	log := models.LogAuditoria{
		UsuarioEmail:    email,
		Accion:          "eliminar",
		Entidad:         "activo",
		EntidadID:       id,
		DatosAnteriores: datosAnteriores,
		Descripcion:     fmt.Sprintf("Usuario %s eliminó el activo '%s' (ID: %d)", email, activo.Nombre, id),
		IPOrigen:        stringPtr(c.ClientIP()),
		UserAgent:       stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusOK, gin.H{"message": "Activo eliminado exitosamente"})
}

// ==================== EDIFICIOS ====================

// CrearEdificio crea un edificio
func (h *AdminHandler) CrearEdificio(c *gin.Context) {
	var req models.CreateEdificioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crear edificio
	edificio := models.Edificio{
		Nombre:    req.Nombre,
		Direccion: req.Direccion,
		Latitud:   req.Latitud,
		Longitud:  req.Longitud,
	}

	edificioCreado, err := h.edificioRepo.Crear(edificio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear edificio", "details": err.Error()})
		return
	}

	// Registrar en log
	datosNuevos := map[string]interface{}{
		"id":        edificioCreado.ID,
		"nombre":    edificioCreado.Nombre,
		"direccion": edificioCreado.Direccion,
		"latitud":   edificioCreado.Latitud,
		"longitud":  edificioCreado.Longitud,
	}

	log := models.LogAuditoria{
		UsuarioEmail: req.Email,
		Accion:       "crear",
		Entidad:      "edificio",
		EntidadID:    edificioCreado.ID,
		DatosNuevos:  datosNuevos,
		Descripcion:  fmt.Sprintf("Usuario %s creó el edificio '%s' (ID: %d)", req.Email, edificioCreado.Nombre, edificioCreado.ID),
		IPOrigen:     stringPtr(c.ClientIP()),
		UserAgent:    stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Edificio creado exitosamente",
		"edificio": edificioCreado,
	})
}

// ActualizarEdificio actualiza un edificio
func (h *AdminHandler) ActualizarEdificio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdateEdificioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener edificio actual
	edificioAnterior, err := h.edificioRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Edificio no encontrado"})
		return
	}

	// Preparar datos actualizados
	edificioActualizado := *edificioAnterior
	if req.Nombre != nil {
		edificioActualizado.Nombre = *req.Nombre
	}
	if req.Direccion != nil {
		edificioActualizado.Direccion = *req.Direccion
	}
	if req.Latitud != nil {
		edificioActualizado.Latitud = req.Latitud
	}
	if req.Longitud != nil {
		edificioActualizado.Longitud = req.Longitud
	}

	// Actualizar
	err = h.edificioRepo.Actualizar(id, edificioActualizado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar edificio", "details": err.Error()})
		return
	}

	// Registrar en log
	datosAnteriores := map[string]interface{}{
		"id":        edificioAnterior.ID,
		"nombre":    edificioAnterior.Nombre,
		"direccion": edificioAnterior.Direccion,
		"latitud":   edificioAnterior.Latitud,
		"longitud":  edificioAnterior.Longitud,
	}

	datosNuevos := map[string]interface{}{
		"id":        edificioActualizado.ID,
		"nombre":    edificioActualizado.Nombre,
		"direccion": edificioActualizado.Direccion,
		"latitud":   edificioActualizado.Latitud,
		"longitud":  edificioActualizado.Longitud,
	}

	log := models.LogAuditoria{
		UsuarioEmail:    req.Email,
		Accion:          "modificar",
		Entidad:         "edificio",
		EntidadID:       id,
		DatosAnteriores: datosAnteriores,
		DatosNuevos:     datosNuevos,
		Descripcion:     fmt.Sprintf("Usuario %s modificó el edificio '%s' (ID: %d)", req.Email, edificioActualizado.Nombre, id),
		IPOrigen:        stringPtr(c.ClientIP()),
		UserAgent:       stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Edificio actualizado exitosamente",
		"edificio": edificioActualizado,
	})
}

// EliminarEdificio elimina un edificio
func (h *AdminHandler) EliminarEdificio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email requerido"})
		return
	}

	// Obtener edificio antes de eliminar
	edificio, err := h.edificioRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Edificio no encontrado"})
		return
	}

	// Eliminar
	err = h.edificioRepo.Eliminar(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar edificio", "details": err.Error()})
		return
	}

	// Registrar en log
	datosAnteriores := map[string]interface{}{
		"id":        edificio.ID,
		"nombre":    edificio.Nombre,
		"direccion": edificio.Direccion,
		"latitud":   edificio.Latitud,
		"longitud":  edificio.Longitud,
	}

	log := models.LogAuditoria{
		UsuarioEmail:    email,
		Accion:          "eliminar",
		Entidad:         "edificio",
		EntidadID:       id,
		DatosAnteriores: datosAnteriores,
		Descripcion:     fmt.Sprintf("Usuario %s eliminó el edificio '%s' (ID: %d)", email, edificio.Nombre, id),
		IPOrigen:        stringPtr(c.ClientIP()),
		UserAgent:       stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusOK, gin.H{"message": "Edificio eliminado exitosamente"})
}

// ==================== SENSORES ====================

// CrearSensor crea un sensor en ParserService
func (h *AdminHandler) CrearSensor(c *gin.Context) {
	var req models.CreateSensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el activo existe
	_, err := h.activoRepo.GetByID(req.ActivoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activo no encontrado"})
		return
	}

	// Crear sensor en ParserService
	parserReq := map[string]interface{}{
		"id_activo": req.ActivoID,
		"nombre":    req.Nombre,
		"tipo":      req.Tipo,
		"unidad":    req.Unidad,
	}

	jsonData, _ := json.Marshal(parserReq)
	resp, err := http.Post(h.parserURL+"/admin/sensores", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear sensor en ParserService", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ParserService rechazó la creación", "details": string(bodyBytes)})
		return
	}

	// Leer respuesta de ParserService
	var parserResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&parserResp)

	// Registrar en log
	log := models.LogAuditoria{
		UsuarioEmail: req.Email,
		Accion:       "crear",
		Entidad:      "sensor",
		EntidadID:    req.ActivoID, // Usar activo_id como referencia
		DatosNuevos: map[string]interface{}{
			"activo_id": req.ActivoID,
			"nombre":    req.Nombre,
			"tipo":      req.Tipo,
			"unidad":    req.Unidad,
		},
		Descripcion: fmt.Sprintf("Usuario %s creó el sensor '%s' para activo ID %d", req.Email, req.Nombre, req.ActivoID),
		IPOrigen:    stringPtr(c.ClientIP()),
		UserAgent:   stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sensor creado exitosamente",
		"sensor":  parserResp,
	})
}

// ActualizarSensor actualiza un sensor en ParserService
func (h *AdminHandler) ActualizarSensor(c *gin.Context) {
	sensorID := c.Param("id")

	var req struct {
		Email  string `json:"email" binding:"required,email"`
		Nombre string `json:"nombre"`
		Tipo   string `json:"tipo"`
		Unidad string `json:"unidad"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preparar datos anteriores (obtener del ParserService si es necesario)
	datosAnteriores := map[string]interface{}{
		"sensor_id": sensorID,
	}

	// Actualizar sensor en ParserService
	parserReq := map[string]interface{}{
		"nombre": req.Nombre,
		"tipo":   req.Tipo,
		"unidad": req.Unidad,
	}

	jsonData, _ := json.Marshal(parserReq)
	client := &http.Client{}
	httpReq, _ := http.NewRequest("PUT", h.parserURL+"/admin/sensores/"+sensorID, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar sensor en ParserService", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ParserService rechazó la actualización", "details": string(bodyBytes)})
		return
	}

	// Registrar en log
	log := models.LogAuditoria{
		UsuarioEmail:    req.Email,
		Accion:          "modificar",
		Entidad:         "sensor",
		EntidadID:       0, // No hay ID numérico para sensores
		DatosAnteriores: datosAnteriores,
		DatosNuevos: map[string]interface{}{
			"sensor_id": sensorID,
			"nombre":    req.Nombre,
			"tipo":      req.Tipo,
			"unidad":    req.Unidad,
		},
		Descripcion: fmt.Sprintf("Usuario %s actualizó el sensor '%s'", req.Email, sensorID),
		IPOrigen:    stringPtr(c.ClientIP()),
		UserAgent:   stringPtr(c.Request.UserAgent()),
	}

	h.logRepo.CrearLog(log)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sensor actualizado exitosamente",
		"sensor_id": sensorID,
	})
}

// ==================== LOGS ====================

// ObtenerLogs obtiene logs de auditoría con filtros
func (h *AdminHandler) ObtenerLogs(c *gin.Context) {
	// Parámetros de paginación
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Filtros opcionales
	usuarioID := c.Query("usuario_id")
	entidad := c.Query("entidad")
	entidadID := c.Query("entidad_id")

	var logs []models.LogAuditoria
	var err error

	if usuarioID != "" {
		uid, _ := strconv.Atoi(usuarioID)
		logs, err = h.logRepo.ObtenerLogsPorUsuario(uid, limit, offset)
	} else if entidad != "" && entidadID != "" {
		eid, _ := strconv.Atoi(entidadID)
		logs, err = h.logRepo.ObtenerLogsPorEntidad(entidad, eid)
	} else {
		logs, err = h.logRepo.ObtenerTodosLosLogs(limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener logs", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  len(logs),
		"limit":  limit,
		"offset": offset,
	})
}

// ==================== RELACIONES USUARIOS-EDIFICIOS ====================

// AsignarUsuarioAEdificio asigna un usuario a un edificio
func (h *AdminHandler) AsignarUsuarioAEdificio(c *gin.Context) {
	var req struct {
		Email      string `json:"email" binding:"required"`
		EdificioID int    `json:"edificio_id" binding:"required"`
		AdminEmail string `json:"admin_email" binding:"required"` // Usuario que realiza la asignación
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el usuario existe
	usuario, err := h.usuarioRepo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar que el edificio existe
	edificio, err := h.edificioRepo.GetByID(req.EdificioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Edificio no encontrado"})
		return
	}

	// Asignar usuario al edificio
	err = h.usuarioRepo.AsignarEdificio(usuario.ID, req.EdificioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar usuario al edificio", "details": err.Error()})
		return
	}

	// Registrar log de auditoría
	adminUsuario, _ := h.usuarioRepo.GetByEmail(req.AdminEmail)
	if adminUsuario != nil {
		logData := map[string]interface{}{
			"usuario_id":      usuario.ID,
			"usuario_email":   req.Email,
			"edificio_id":     req.EdificioID,
			"edificio_nombre": edificio.Nombre,
		}
		h.logRepo.CrearLog(models.LogAuditoria{
			UsuarioID:    &adminUsuario.ID,
			UsuarioEmail: req.AdminEmail,
			Accion:       "ASIGNAR_USUARIO_EDIFICIO",
			Entidad:      "usuario_edificio",
			EntidadID:    usuario.ID,
			DatosNuevos:  logData,
			Descripcion:  fmt.Sprintf("Usuario %s asignado al edificio %s", req.Email, edificio.Nombre),
			IPOrigen:     stringPtr(c.ClientIP()),
			UserAgent:    stringPtr(c.Request.UserAgent()),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario asignado al edificio exitosamente",
		"usuario": gin.H{
			"id":    usuario.ID,
			"email": usuario.Email,
		},
		"edificio": gin.H{
			"id":     edificio.ID,
			"nombre": edificio.Nombre,
		},
	})
}

// RemoverUsuarioDeEdificio remueve la asignación de un usuario a un edificio
func (h *AdminHandler) RemoverUsuarioDeEdificio(c *gin.Context) {
	var req struct {
		Email      string `json:"email" binding:"required"`
		EdificioID int    `json:"edificio_id" binding:"required"`
		AdminEmail string `json:"admin_email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el usuario existe
	usuario, err := h.usuarioRepo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Remover asignación
	err = h.usuarioRepo.RemoverEdificio(usuario.ID, req.EdificioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al remover asignación", "details": err.Error()})
		return
	}

	// Registrar log de auditoría
	adminUsuario, _ := h.usuarioRepo.GetByEmail(req.AdminEmail)
	if adminUsuario != nil {
		logData := map[string]interface{}{
			"usuario_id":    usuario.ID,
			"usuario_email": req.Email,
			"edificio_id":   req.EdificioID,
		}
		h.logRepo.CrearLog(models.LogAuditoria{
			UsuarioID:       &adminUsuario.ID,
			UsuarioEmail:    req.AdminEmail,
			Accion:          "REMOVER_USUARIO_EDIFICIO",
			Entidad:         "usuario_edificio",
			EntidadID:       usuario.ID,
			DatosAnteriores: logData,
			Descripcion:     fmt.Sprintf("Usuario %s removido del edificio %d", req.Email, req.EdificioID),
			IPOrigen:        stringPtr(c.ClientIP()),
			UserAgent:       stringPtr(c.Request.UserAgent()),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario removido del edificio exitosamente",
	})
}

// ObtenerUsuariosDeEdificio obtiene todos los usuarios asignados a un edificio
func (h *AdminHandler) ObtenerUsuariosDeEdificio(c *gin.Context) {
	edificioID, err := strconv.Atoi(c.Param("edificio_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de edificio inválido"})
		return
	}

	usuarios, err := h.usuarioRepo.GetUsuariosPorEdificio(edificioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"edificio_id": edificioID,
		"usuarios":    usuarios,
		"total":       len(usuarios),
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
