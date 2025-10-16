package services

import (
	"database/sql"
	"fmt"
	"log"
	"notification/internal/models"
	"notification/internal/repository"
	"time"

	_ "github.com/lib/pq"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	emailService     *EmailService
}

func NewNotificationService(repo *repository.NotificationRepository, emailSvc *EmailService) *NotificationService {
	return &NotificationService{
		notificationRepo: repo,
		emailService:     emailSvc,
	}
}

func (s *NotificationService) CreateSensorAlert(request models.SensorAlertRequest) (*models.Notificacion, error) {
	// Validar prioridad por defecto
	if request.Prioridad == "" {
		request.Prioridad = "medium"
	}

	// Conectar directamente a PostgreSQL para insertar en la nueva estructura
	db, err := s.connectToPostgres()
	if err != nil {
		return nil, fmt.Errorf("error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Si tenemos asset_id, obtener información del activo y edificio
	var activoInfo *models.Activo
	var edificioInfo *models.Edificio
	var usuarios []models.Usuario

	if request.AssetID != 0 {
		// Obtener información del activo con su edificio
		query := `
			SELECT a.id, a.nombre, a.edificio_id, e.direccion 
			FROM activos a 
			JOIN edificios e ON a.edificio_id = e.id 
			WHERE a.id = $1`

		activoInfo = &models.Activo{}
		edificioInfo = &models.Edificio{}

		err = db.QueryRow(query, request.AssetID).Scan(
			&activoInfo.ID,
			&activoInfo.Nombre,
			&activoInfo.EdificioID,
			&edificioInfo.Direccion,
		)

		if err != nil {
			log.Printf("⚠️  Error al obtener información del activo %d: %v", request.AssetID, err)
		} else {
			edificioInfo.ID = activoInfo.EdificioID

			// Obtener usuarios del edificio
			usuariosQuery := `
				SELECT id, usuario, correo, numero, scope 
				FROM usuarios 
				WHERE edificio_id = $1`

			rows, err := db.Query(usuariosQuery, activoInfo.EdificioID)
			if err != nil {
				log.Printf("⚠️  Error al obtener usuarios del edificio %d: %v", activoInfo.EdificioID, err)
			} else {
				defer rows.Close()
				for rows.Next() {
					var usuario models.Usuario
					err := rows.Scan(&usuario.ID, &usuario.Usuario, &usuario.Correo, &usuario.Numero, &usuario.Scope)
					if err != nil {
						log.Printf("Error al escanear usuario: %v", err)
						continue
					}
					usuario.EdificioID = &activoInfo.EdificioID
					usuarios = append(usuarios, usuario)
				}
			}
		}
	}

	// Crear la notificación en la base de datos
	var notificationID int
	query := `
		INSERT INTO notificaciones (
			building_id, asset_id, sensor_id, message, alert_type, 
			tipo, prioridad, notification_mail, notification_sms, 
			status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id`

	var buildingIDParam, assetIDParam interface{}
	if request.BuildingID != 0 {
		buildingIDParam = request.BuildingID
	}
	if request.AssetID != 0 {
		assetIDParam = request.AssetID
	}

	err = db.QueryRow(
		query,
		buildingIDParam,
		assetIDParam,
		request.SensorID,
		request.Message,
		request.AlertType,
		request.Tipo,
		request.Prioridad,
		request.NotificationMail,
		request.NotificationSMS,
		"pending",
		time.Now(),
	).Scan(&notificationID)

	if err != nil {
		return nil, fmt.Errorf("error al crear notificación: %v", err)
	}

	// Crear objeto de respuesta
	notification := &models.Notificacion{
		ID:               notificationID,
		SensorID:         request.SensorID,
		Message:          request.Message,
		AlertType:        request.AlertType,
		Tipo:             request.Tipo,
		Prioridad:        request.Prioridad,
		NotificationMail: request.NotificationMail,
		NotificationSMS:  request.NotificationSMS,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	if buildingIDParam != nil {
		notification.BuildingID = &request.BuildingID
	}
	if assetIDParam != nil {
		notification.AssetID = &request.AssetID
	}

	// Enviar email personalizado si se solicita y tenemos información del activo
	if request.NotificationMail {
		if activoInfo != nil && edificioInfo != nil && len(usuarios) > 0 {
			err = s.sendPersonalizedEmailNotification(notification, activoInfo, edificioInfo, usuarios)
		} else {
			// Si no tenemos información del activo, usar email genérico
			err = s.sendGenericEmailNotification(notification)
		}

		if err != nil {
			log.Printf("Error al enviar email: %v", err)
		} else {
			// Actualizar estado
			db.Exec("UPDATE notificaciones SET status = 'sent' WHERE id = $1", notificationID)
		}
	}

	log.Printf("✅ Notificación creada exitosamente: ID=%d, Tipo=%s, Prioridad=%s",
		notificationID, request.Tipo, request.Prioridad)

	return notification, nil
}

func (s *NotificationService) connectToPostgres() (*sql.DB, error) {
	dsn := "host=notification-db port=5432 user=notification_user password=notification_pass dbname=notification_db sslmode=disable"
	return sql.Open("postgres", dsn)
}

func (s *NotificationService) sendPersonalizedEmailNotification(notification *models.Notificacion, activo *models.Activo, edificio *models.Edificio, usuarios []models.Usuario) error {
	// Construir el contenido del email personalizado con información completa
	prioridadEmoji := "⚠️"
	switch notification.Prioridad {
	case "critical":
		prioridadEmoji = "🚨"
	case "high":
		prioridadEmoji = "🔴"
	case "medium":
		prioridadEmoji = "🟡"
	case "low":
		prioridadEmoji = "🟢"
	}

	tipoEmoji := "📊"
	switch notification.Tipo {
	case "sensor":
		tipoEmoji = "📡"
	case "alerta":
		tipoEmoji = "⚠️"
	case "mantenimiento":
		tipoEmoji = "🔧"
	case "sistema":
		tipoEmoji = "💻"
	}

	asunto := fmt.Sprintf("%s ALERTA: Sensor %s - %s", prioridadEmoji, notification.SensorID, activo.Nombre)

	contenido := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; background-color: #f9f9f9; padding: 20px;">
			<div style="background-color: white; padding: 25px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
				<h2 style="color: #d32f2f; margin-top: 0;">%s Alerta de Sensor Detectada</h2>
				
				<div style="background-color: #ffebee; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #d32f2f;">
					<h3 style="margin: 0; color: #d32f2f;">%s Información del Sensor</h3>
					<p style="margin: 10px 0;"><strong>🆔 Sensor ID:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>📝 Descripción:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>🏷️ Tipo de Alerta:</strong> %s</p>
				</div>

				<div style="background-color: #e3f2fd; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #1976d2;">
					<h3 style="margin: 0; color: #1976d2;">🏢 Información del Activo</h3>
					<p style="margin: 10px 0;"><strong>⚙️ Activo:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>🏠 Edificio:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>🆔 ID del Activo:</strong> %d</p>
				</div>

				<div style="background-color: #f3e5f5; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #7b1fa2;">
					<h3 style="margin: 0; color: #7b1fa2;">⚡ Detalles de la Alerta</h3>
					<p style="margin: 10px 0;"><strong>%s Tipo:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>%s Prioridad:</strong> %s</p>
					<p style="margin: 10px 0;"><strong>🕒 Fecha y Hora:</strong> %s</p>
				</div>

				<div style="background-color: #fff3e0; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 4px solid #f57c00;">
					<p style="margin: 0; color: #e65100;"><strong>⚡ Acción Requerida:</strong> Se recomienda revisar inmediatamente el sensor y el activo asociado para determinar la causa de la desconexión.</p>
				</div>

				<hr style="border: none; border-top: 1px solid #e0e0e0; margin: 25px 0;">
				<p style="color: #666; font-size: 12px; text-align: center; margin: 0;">
					<em>Este es un mensaje automático del sistema InspectAR.<br>
					Sistema de Monitoreo Inteligente de Activos</em>
				</p>
			</div>
		</div>
	`,
		prioridadEmoji,
		tipoEmoji,
		notification.SensorID,
		notification.Message,
		notification.AlertType,
		activo.Nombre,
		edificio.Direccion,
		activo.ID,
		tipoEmoji,
		notification.Tipo,
		prioridadEmoji,
		notification.Prioridad,
		notification.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	// Enviar email a cada usuario del edificio
	emailsSent := 0
	for _, usuario := range usuarios {
		err := s.emailService.SendAlertEmail(usuario.Correo, asunto, contenido)
		if err != nil {
			log.Printf("❌ Error al enviar email a %s (%s): %v", usuario.Usuario, usuario.Correo, err)
			continue
		}
		emailsSent++
		log.Printf("📧 Email personalizado enviado exitosamente a %s (%s) - %s", usuario.Usuario, usuario.Correo, usuario.Scope)
	}

	if emailsSent > 0 {
		log.Printf("✅ %d emails personalizados enviados para alerta del sensor %s en activo '%s'", emailsSent, notification.SensorID, activo.Nombre)
	}

	return nil
}

func (s *NotificationService) sendGenericEmailNotification(notification *models.Notificacion) error {
	// Email genérico cuando no tenemos información del activo
	asunto := fmt.Sprintf("🚨 ALERTA: %s", notification.AlertType)

	contenido := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; background-color: #f9f9f9; padding: 20px;">
			<div style="background-color: white; padding: 25px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
				<h2 style="color: #d32f2f;">🚨 Alerta de Sistema</h2>
				<div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 10px 0;">
					<p><strong>🔍 Sensor ID:</strong> %s</p>
					<p><strong>📝 Mensaje:</strong> %s</p>
					<p><strong>🏷️ Tipo:</strong> %s</p>
					<p><strong>⚡ Prioridad:</strong> %s</p>
					<p><strong>🕒 Fecha y Hora:</strong> %s</p>
				</div>
				<p><em>Este es un mensaje automático del sistema InspectAR.</em></p>
			</div>
		</div>
	`,
		notification.SensorID,
		notification.Message,
		notification.Tipo,
		notification.Prioridad,
		notification.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	// Enviar a admin por defecto
	return s.emailService.SendAlertEmail("admin@inspectarar.com", asunto, contenido)
}

func (s *NotificationService) sendEmailNotification(notification *models.Notificacion) error {
	// Usar el servicio de email existente
	asunto := fmt.Sprintf("🚨 ALERTA: %s", notification.AlertType)

	contenido := fmt.Sprintf(`
		<h2>🚨 Alerta de Sistema</h2>
		<div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 10px 0;">
			<p><strong>🔍 Sensor ID:</strong> %s</p>
			<p><strong>📝 Mensaje:</strong> %s</p>
			<p><strong>🏷️ Tipo:</strong> %s</p>
			<p><strong>⚡ Prioridad:</strong> %s</p>
			<p><strong>🕒 Fecha y Hora:</strong> %s</p>
		</div>
		<p><em>Este es un mensaje automático del sistema InspectAR.</em></p>
	`,
		notification.SensorID,
		notification.Message,
		notification.Tipo,
		notification.Prioridad,
		notification.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	// Enviar a admin por defecto (debería obtener usuarios desde la BD)
	return s.emailService.SendAlertEmail("admin@inspectarar.com", asunto, contenido)
}

// CreateSensorAlertOld - DEPRECATED: Método comentado para la transición
// Este método usaba la estructura antigua con tipos_notificacion table
/*
func (s *NotificationService) CreateSensorAlertOld(request models.SensorAlertRequest) (*models.Notificacion, error) {
	// MÉTODO OBSOLETO - Usar CreateSensorAlert en su lugar
	return nil, fmt.Errorf("método obsoleto, usar CreateSensorAlert")
}
*/

func (s *NotificationService) enviarNotificacionEmail(notificacion *models.Notificacion, usuarios []models.Usuario) error {
	// Construir el contenido del email con información detallada
	asunto := fmt.Sprintf("🚨 ALERTA: %s", notificacion.AlertType)

	contenido := fmt.Sprintf(`
		<h2>🚨 Alerta de Sensor Detectada</h2>
		<p><strong>Sensor ID:</strong> %s</p>
		<p><strong>Mensaje:</strong> %s</p>
		<p><strong>Tipo:</strong> %s</p>
		<p><strong>Prioridad:</strong> %s</p>
		<p><strong>Fecha y Hora:</strong> %s</p>
	`, notificacion.SensorID, notificacion.Message, notificacion.Tipo, notificacion.Prioridad, notificacion.CreatedAt.Format("2006-01-02 15:04:05"))

	// Enviar email a cada usuario
	for _, usuario := range usuarios {
		err := s.emailService.SendAlertEmail(usuario.Correo, asunto, contenido)
		if err != nil {
			log.Printf("Error al enviar email a %s: %v", usuario.Correo, err)
			continue
		}
		log.Printf("Email de alerta enviado exitosamente a %s", usuario.Correo)
	}

	return nil
}

func (s *NotificationService) enviarNotificacionSMS(notificacion *models.Notificacion, usuarios []models.Usuario) error {
	// Construir mensaje SMS (limitado en caracteres)
	mensaje := fmt.Sprintf("🚨 ALERTA: Sensor %s - %s. Tiempo: %s",
		notificacion.SensorID,
		notificacion.AlertType,
		notificacion.CreatedAt.Format("15:04"))

	// Enviar SMS a cada usuario con número telefónico
	for _, usuario := range usuarios {
		if usuario.Numero == "" {
			log.Printf("Usuario %s no tiene número telefónico configurado", usuario.Usuario)
			continue
		}

		// Aquí iría la integración con servicio SMS (Twilio, AWS SNS, etc.)
		log.Printf("📱 SMS enviado a %s (%s): %s", usuario.Usuario, usuario.Numero, mensaje)
		// TODO: Implementar envío real de SMS
	}

	return nil
}

func (s *NotificationService) CreateNotification(activoID uint) (*models.Notificacion, error) {
	// Verificar que el activo existe y cargar las relaciones
	activo, err := s.notificationRepo.GetActivoByID(activoID)
	if err != nil {
		return nil, fmt.Errorf("activo no encontrado: %v", err)
	}

	// Conectar directamente a PostgreSQL para crear en la nueva estructura
	db, err := s.connectToPostgres()
	if err != nil {
		return nil, fmt.Errorf("error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Crear la notificación directamente en PostgreSQL
	var notificationID int
	query := `
		INSERT INTO notificaciones (
			asset_id, message, alert_type, tipo, prioridad, 
			notification_mail, notification_sms, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`

	err = db.QueryRow(
		query,
		activoID,
		fmt.Sprintf("Alerta detectada en %s", activo.Nombre),
		"offline",
		"alerta",
		"medium",
		true,
		false,
		"pending",
		time.Now(),
	).Scan(&notificationID)

	if err != nil {
		return nil, fmt.Errorf("error al crear notificación: %v", err)
	}

	// Crear objeto de respuesta
	assetIDInt := int(activoID)
	notification := &models.Notificacion{
		ID:               notificationID,
		AssetID:          &assetIDInt,
		Message:          fmt.Sprintf("Alerta detectada en %s", activo.Nombre),
		AlertType:        "offline",
		Tipo:             "alerta",
		Prioridad:        "medium",
		NotificationMail: true,
		NotificationSMS:  false,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	// Obtener todos los usuarios del mismo edificio que el activo
	usuarios, err := s.notificationRepo.GetUsuariosByEdificioID(activo.EdificioID)
	if err != nil {
		log.Printf("Error al obtener usuarios del edificio %d: %v", activo.EdificioID, err)
		// Continuar aunque no se puedan enviar correos
		return notification, nil
	}

	// Enviar correo electrónico a todos los usuarios del edificio
	if len(usuarios) > 0 {
		err = s.emailService.SendNotificationEmail(activo, usuarios, notification.Message)
		if err != nil {
			log.Printf("Error al enviar correo: %v", err)
			// No retornar error, ya que la notificación se creó correctamente
		} else {
			log.Printf("Correo enviado exitosamente a %d usuarios", len(usuarios))
			// Actualizar estado
			db.Exec("UPDATE notificaciones SET status = 'sent' WHERE id = $1", notificationID)
		}
	}

	return notification, nil
}

func (s *NotificationService) GetAllTiposNotificacion() ([]models.TipoNotificacion, error) {
	return s.notificationRepo.GetAllTiposNotificacion()
}

func (s *NotificationService) GetNotificationsByActivoID(activoID uint) ([]models.Notificacion, error) {
	return s.notificationRepo.GetNotificationsByActivoID(activoID)
}

func (s *NotificationService) SendNotification(notificationID uint) error {
	// Conectar directamente a PostgreSQL para actualizar
	db, err := s.connectToPostgres()
	if err != nil {
		return fmt.Errorf("error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Verificar que la notificación existe y actualizar estado
	_, err = db.Exec("UPDATE notificaciones SET status = 'sent' WHERE id = $1", notificationID)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %v", err)
	}

	return nil
}

func (s *NotificationService) GetAllEdificios() ([]models.Edificio, error) {
	return s.notificationRepo.GetAllEdificios()
}

func (s *NotificationService) GetEdificioByID(id uint) (*models.Edificio, error) {
	return s.notificationRepo.GetEdificioByID(id)
}

func (s *NotificationService) GetUsuariosByEdificioID(edificioID uint) ([]models.Usuario, error) {
	return s.notificationRepo.GetUsuariosByEdificioID(edificioID)
}

func (s *NotificationService) GetActivosByEdificioID(edificioID uint) ([]models.Activo, error) {
	return s.notificationRepo.GetActivosByEdificioID(edificioID)
}

// SendTechnicianContactRequest envía una solicitud de contacto a un técnico sin guardar en la base de datos
func (s *NotificationService) SendTechnicianContactRequest(request models.TechnicianContactRequest) error {
	// Validar prioridad por defecto
	if request.Priority == "" {
		request.Priority = "medium"
	}

	// Conectar a PostgreSQL para obtener información del activo
	db, err := s.connectToPostgres()
	if err != nil {
		return fmt.Errorf("error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Obtener información del activo y edificio
	var activoNombre, edificioDireccion string
	query := `
		SELECT a.nombre, e.direccion 
		FROM activos a 
		JOIN edificios e ON a.edificio_id = e.id 
		WHERE a.id = $1`

	err = db.QueryRow(query, request.ActivoID).Scan(&activoNombre, &edificioDireccion)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("activo con ID %d no encontrado", request.ActivoID)
		}
		return fmt.Errorf("error obteniendo información del activo: %v", err)
	}

	// Mapear prioridad a texto legible
	priorityText := map[string]string{
		"low":    "Baja",
		"medium": "Media",
		"high":   "Alta",
		"urgent": "Urgente",
	}
	priorityDisplay := priorityText[request.Priority]
	if priorityDisplay == "" {
		priorityDisplay = "Media"
	}

	// Crear el template del correo para el técnico
	subject := fmt.Sprintf("🔧 Solicitud de Contacto Técnico - Activo: %s", activoNombre)
	
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; border-radius: 10px 10px 0 0; text-align: center; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 10px 10px; }
        .priority-%s { border-left: 5px solid %s; padding-left: 15px; margin: 15px 0; }
        .info-box { background: white; padding: 15px; border-radius: 8px; margin: 15px 0; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        .contact-info { background: #e3f2fd; padding: 15px; border-radius: 8px; margin: 15px 0; }
        .btn { display: inline-block; background: #667eea; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; margin: 10px 0; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔧 Solicitud de Contacto Técnico</h1>
            <p>InspectAR - Sistema de Gestión de Activos</p>
        </div>
        <div class="content">
            <div class="priority-%s">
                <h3>Prioridad: %s</h3>
            </div>
            
            <div class="info-box">
                <h3>📍 Información del Activo</h3>
                <p><strong>Activo:</strong> %s</p>
                <p><strong>ID del Activo:</strong> %d</p>
                <p><strong>Edificio:</strong> %s</p>
            </div>
            
            <div class="contact-info">
                <h3>👤 Información del Usuario</h3>
                <p><strong>Nombre:</strong> %s</p>
                <p><strong>Email de contacto:</strong> %s</p>
            </div>
            
            <div class="info-box">
                <h3>💬 Mensaje del Usuario</h3>
                <p><em>"%s"</em></p>
            </div>
            
            <div style="text-align: center; margin: 20px 0;">
                <a href="mailto:%s?subject=Re: Solicitud Técnica - %s&body=Hola %s,%%0D%%0A%%0D%%0ARecibí tu solicitud técnica...%%0D%%0A%%0D%%0ASaludos,%%0D%%0ATécnico" class="btn">
                    📧 Responder al Usuario
                </a>
            </div>
            
            <div class="footer">
                <p>Este correo fue generado automáticamente por el sistema InspectAR</p>
                <p>Fecha: %s</p>
            </div>
        </div>
    </div>
</body>
</html>`,
		request.Priority, s.getPriorityColor(request.Priority),
		request.Priority, priorityDisplay,
		activoNombre, request.ActivoID, edificioDireccion,
		request.UserName, request.UserEmail,
		request.Message,
		request.UserEmail, activoNombre, request.UserName,
		time.Now().Format("2006-01-02 15:04:05"))

	// Enviar correo al técnico
	err = s.emailService.SendAlertEmail(request.TechnicianEmail, subject, htmlBody)
	if err != nil {
		return fmt.Errorf("error enviando correo al técnico: %v", err)
	}

	log.Printf("✅ Solicitud de contacto técnico enviada exitosamente a: %s", request.TechnicianEmail)
	log.Printf("📧 Usuario solicitante: %s (%s)", request.UserName, request.UserEmail)
	log.Printf("🏢 Activo: %s (ID: %d)", activoNombre, request.ActivoID)

	return nil
}

// getPriorityColor retorna el color asociado a cada prioridad
func (s *NotificationService) getPriorityColor(priority string) string {
	colors := map[string]string{
		"low":    "#4caf50", // Verde
		"medium": "#ff9800", // Naranja
		"high":   "#f44336", // Rojo
		"urgent": "#9c27b0", // Púrpura
	}
	if color, exists := colors[priority]; exists {
		return color
	}
	return "#ff9800" // Default naranja
}
