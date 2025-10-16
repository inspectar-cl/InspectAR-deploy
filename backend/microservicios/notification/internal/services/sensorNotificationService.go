package services

import (
	"database/sql"
	"fmt"
	"log"
	"notification/internal/models"
	"time"
)

type SensorNotificationService struct {
	db           *sql.DB
	emailService *EmailService
}

func NewSensorNotificationService(db *sql.DB, emailSvc *EmailService) *SensorNotificationService {
	return &SensorNotificationService{
		db:           db,
		emailService: emailSvc,
	}
}

func (s *SensorNotificationService) CreateSensorAlert(request models.SensorAlertRequest) (*models.Notificacion, error) {
	// Validar prioridad por defecto
	if request.Prioridad == "" {
		request.Prioridad = "medium"
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

	err := s.db.QueryRow(
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

	// Obtener la notificación creada
	notification, err := s.GetNotificationByID(notificationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener notificación creada: %v", err)
	}

	// Enviar notificaciones si se solicita
	if request.NotificationMail {
		err = s.sendEmailNotification(notification)
		if err != nil {
			log.Printf("Error al enviar email: %v", err)
		} else {
			// Actualizar estado
			s.updateNotificationStatus(notificationID, "sent")
		}
	}

	log.Printf("Notificación creada exitosamente: ID=%d, Tipo=%s, Prioridad=%s",
		notificationID, request.Tipo, request.Prioridad)

	return notification, nil
}

func (s *SensorNotificationService) GetNotificationByID(id int) (*models.Notificacion, error) {
	query := `
		SELECT id, building_id, asset_id, sensor_id, message, alert_type, 
			   tipo, prioridad, notification_mail, notification_sms, 
			   status, created_at
		FROM notificaciones 
		WHERE id = $1`

	notification := &models.Notificacion{}
	row := s.db.QueryRow(query, id)

	var buildingID, assetID sql.NullInt32

	err := row.Scan(
		&notification.ID,
		&buildingID,
		&assetID,
		&notification.SensorID,
		&notification.Message,
		&notification.AlertType,
		&notification.Tipo,
		&notification.Prioridad,
		&notification.NotificationMail,
		&notification.NotificationSMS,
		&notification.Status,
		&notification.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if buildingID.Valid {
		notification.BuildingID = &[]int{int(buildingID.Int32)}[0]
	}
	if assetID.Valid {
		notification.AssetID = &[]int{int(assetID.Int32)}[0]
	}

	return notification, nil
}

func (s *SensorNotificationService) GetNotificationsByTipo(tipo string, limit int) ([]models.Notificacion, error) {
	query := `
		SELECT id, building_id, asset_id, sensor_id, message, alert_type, 
			   tipo, prioridad, notification_mail, notification_sms, 
			   status, created_at
		FROM notificaciones 
		WHERE tipo = $1 
		ORDER BY created_at DESC 
		LIMIT $2`

	rows, err := s.db.Query(query, tipo, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notificacion
	for rows.Next() {
		var notification models.Notificacion
		var buildingID, assetID sql.NullInt32

		err := rows.Scan(
			&notification.ID,
			&buildingID,
			&assetID,
			&notification.SensorID,
			&notification.Message,
			&notification.AlertType,
			&notification.Tipo,
			&notification.Prioridad,
			&notification.NotificationMail,
			&notification.NotificationSMS,
			&notification.Status,
			&notification.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if buildingID.Valid {
			notification.BuildingID = &[]int{int(buildingID.Int32)}[0]
		}
		if assetID.Valid {
			notification.AssetID = &[]int{int(assetID.Int32)}[0]
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *SensorNotificationService) GetNotificationsByPrioridad(prioridad string, limit int) ([]models.Notificacion, error) {
	query := `
		SELECT id, building_id, asset_id, sensor_id, message, alert_type, 
			   tipo, prioridad, notification_mail, notification_sms, 
			   status, created_at
		FROM notificaciones 
		WHERE prioridad = $1 
		ORDER BY created_at DESC 
		LIMIT $2`

	rows, err := s.db.Query(query, prioridad, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notificacion
	for rows.Next() {
		var notification models.Notificacion
		var buildingID, assetID sql.NullInt32

		err := rows.Scan(
			&notification.ID,
			&buildingID,
			&assetID,
			&notification.SensorID,
			&notification.Message,
			&notification.AlertType,
			&notification.Tipo,
			&notification.Prioridad,
			&notification.NotificationMail,
			&notification.NotificationSMS,
			&notification.Status,
			&notification.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if buildingID.Valid {
			notification.BuildingID = &[]int{int(buildingID.Int32)}[0]
		}
		if assetID.Valid {
			notification.AssetID = &[]int{int(assetID.Int32)}[0]
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *SensorNotificationService) GetAllNotifications(limit int) ([]models.Notificacion, error) {
	query := `
		SELECT id, building_id, asset_id, sensor_id, message, alert_type, 
			   tipo, prioridad, notification_mail, notification_sms, 
			   status, created_at
		FROM notificaciones 
		ORDER BY created_at DESC 
		LIMIT $1`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notificacion
	for rows.Next() {
		var notification models.Notificacion
		var buildingID, assetID sql.NullInt32

		err := rows.Scan(
			&notification.ID,
			&buildingID,
			&assetID,
			&notification.SensorID,
			&notification.Message,
			&notification.AlertType,
			&notification.Tipo,
			&notification.Prioridad,
			&notification.NotificationMail,
			&notification.NotificationSMS,
			&notification.Status,
			&notification.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if buildingID.Valid {
			notification.BuildingID = &[]int{int(buildingID.Int32)}[0]
		}
		if assetID.Valid {
			notification.AssetID = &[]int{int(assetID.Int32)}[0]
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *SensorNotificationService) sendEmailNotification(notification *models.Notificacion) error {
	// Obtener usuarios para envío
	var usuarios []models.Usuario

	if notification.BuildingID != nil {
		var err error
		usuarios, err = s.getUsersByBuildingID(*notification.BuildingID)
		if err != nil {
			log.Printf("Error al obtener usuarios del edificio %d: %v", *notification.BuildingID, err)
		}
	}

	// Si no hay usuarios específicos del edificio, obtener administradores
	if len(usuarios) == 0 {
		var err error
		usuarios, err = s.getAdminUsers()
		if err != nil {
			return fmt.Errorf("error al obtener usuarios admin: %v", err)
		}
	}

	// Construir el contenido del email
	prioridadIcon := s.getPrioridadIcon(notification.Prioridad)
	tipoIcon := s.getTipoIcon(notification.Tipo)

	asunto := fmt.Sprintf("%s %s ALERTA: %s", prioridadIcon, tipoIcon, notification.AlertType)

	contenido := fmt.Sprintf(`
		<h2>%s Alerta de Sistema</h2>
		<div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 10px 0;">
			<p><strong>🔍 Sensor ID:</strong> %s</p>
			<p><strong>📝 Mensaje:</strong> %s</p>
			<p><strong>🏷️ Tipo:</strong> %s</p>
			<p><strong>⚡ Prioridad:</strong> %s</p>
			<p><strong>🕒 Fecha y Hora:</strong> %s</p>
		</div>
		<p><em>Este es un mensaje automático del sistema InspectAR.</em></p>
	`,
		prioridadIcon,
		notification.SensorID,
		notification.Message,
		notification.Tipo,
		notification.Prioridad,
		notification.CreatedAt.Format("2006-01-02 15:04:05"),
	)

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

func (s *SensorNotificationService) getUsersByBuildingID(buildingID int) ([]models.Usuario, error) {
	query := `SELECT id, scope, usuario, correo, numero, edificio_id 
			  FROM usuarios WHERE edificio_id = $1`

	rows, err := s.db.Query(query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var usuario models.Usuario
		var edificioID sql.NullInt32

		err := rows.Scan(
			&usuario.ID,
			&usuario.Scope,
			&usuario.Usuario,
			&usuario.Correo,
			&usuario.Numero,
			&edificioID,
		)
		if err != nil {
			return nil, err
		}

		if edificioID.Valid {
			usuario.EdificioID = &[]uint{uint(edificioID.Int32)}[0]
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (s *SensorNotificationService) getAdminUsers() ([]models.Usuario, error) {
	query := `SELECT id, scope, usuario, correo, numero, edificio_id 
			  FROM usuarios WHERE scope = 'admin'`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var usuario models.Usuario
		var edificioID sql.NullInt32

		err := rows.Scan(
			&usuario.ID,
			&usuario.Scope,
			&usuario.Usuario,
			&usuario.Correo,
			&usuario.Numero,
			&edificioID,
		)
		if err != nil {
			return nil, err
		}

		if edificioID.Valid {
			usuario.EdificioID = &[]uint{uint(edificioID.Int32)}[0]
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (s *SensorNotificationService) updateNotificationStatus(id int, status string) error {
	query := `UPDATE notificaciones SET status = $1 WHERE id = $2`
	_, err := s.db.Exec(query, status, id)
	return err
}

func (s *SensorNotificationService) getPrioridadIcon(prioridad string) string {
	switch prioridad {
	case "low":
		return "🟢"
	case "medium":
		return "🟡"
	case "high":
		return "🟠"
	case "critical":
		return "🔴"
	default:
		return "⚪"
	}
}

func (s *SensorNotificationService) getTipoIcon(tipo string) string {
	switch tipo {
	case "sensor":
		return "📡"
	case "alerta":
		return "🚨"
	case "mantenimiento":
		return "🔧"
	case "sistema":
		return "⚙️"
	default:
		return "📋"
	}
}
