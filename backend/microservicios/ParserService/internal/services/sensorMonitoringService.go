package services

import (
	"ParserService/internal/models"
	"ParserService/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type SensorMonitoringService struct {
	sensorStatusRepo *repository.SensorStatusRepository
	influxRepo       *repository.InfluxRepository
	notificationURL  string
	timeoutMinutes   int
}

func NewSensorMonitoringService(
	sensorStatusRepo *repository.SensorStatusRepository,
	influxRepo *repository.InfluxRepository,
	notificationURL string,
	timeoutMinutes int,
) *SensorMonitoringService {
	if timeoutMinutes <= 0 {
		timeoutMinutes = 1 // Default: 5 minutos
	}

	return &SensorMonitoringService{
		sensorStatusRepo: sensorStatusRepo,
		influxRepo:       influxRepo,
		notificationURL:  notificationURL,
		timeoutMinutes:   timeoutMinutes,
	}
}

// UpdateSensorActivity registra la actividad de un sensor
func (s *SensorMonitoringService) UpdateSensorActivity(sensorID string, timestamp time.Time) error {
	return s.sensorStatusRepo.UpdateSensorActivity(sensorID, timestamp)
}

// CheckDisconnectedSensors verifica y marca sensores desconectados
func (s *SensorMonitoringService) CheckDisconnectedSensors() error {
	log.Printf("🔍 Verificando sensores desconectados (timeout: %d minutos)...", s.timeoutMinutes)

	// Obtener sensores inactivos
	inactiveSensors, err := s.sensorStatusRepo.GetInactiveSensors(s.timeoutMinutes)
	if err != nil {
		return fmt.Errorf("error al obtener sensores inactivos: %v", err)
	}

	if len(inactiveSensors) == 0 {
		log.Printf("✅ Todos los sensores están activos")
		return nil
	}

	log.Printf("⚠️  Encontrados %d sensores desconectados", len(inactiveSensors))

	// Procesar cada sensor desconectado
	for _, sensor := range inactiveSensors {
		err := s.processSensorDisconnection(sensor)
		if err != nil {
			log.Printf("❌ Error al procesar desconexión del sensor %s: %v", sensor.SensorID, err)
			continue
		}
	}

	return nil
}

// processSensorDisconnection maneja la desconexión de un sensor específico
func (s *SensorMonitoringService) processSensorDisconnection(sensor models.SensorStatus) error {
	// Marcar sensor como inactivo
	err := s.sensorStatusRepo.MarkSensorAsInactive(sensor.SensorID)
	if err != nil {
		return fmt.Errorf("error al marcar sensor como inactivo: %v", err)
	}

	// Crear datos de la desconexión
	disconnection := models.DisconnectedSensor{
		SensorID:       sensor.SensorID,
		LastSeen:       sensor.LastSeen,
		DisconnectedAt: time.Now(),
		TotalReports:   sensor.TotalReports,
	}

	// Enviar notificación
	err = s.sendDisconnectionNotification(disconnection)
	if err != nil {
		log.Printf("⚠️  Error al enviar notificación para sensor %s: %v", sensor.SensorID, err)
		// No retornamos error aquí porque la marcación como inactivo ya se hizo
	}

	timeSinceLastSeen := time.Since(sensor.LastSeen)
	log.Printf("🔴 Sensor %s marcado como desconectado (última actividad: %v ago)",
		sensor.SensorID, timeSinceLastSeen.Round(time.Minute))

	return nil
}

// sendDisconnectionNotification envía una notificación de desconexión
func (s *SensorMonitoringService) sendDisconnectionNotification(disconnection models.DisconnectedSensor) error {
	// Console log como antes
	log.Printf("📧 NOTIFICACIÓN: Sensor %s se desconectó", disconnection.SensorID)
	log.Printf("   Última actividad: %v", disconnection.LastSeen.Format("2006-01-02 15:04:05"))
	log.Printf("   Total reportes: %d", disconnection.TotalReports)
	log.Printf("   Desconectado en: %v", disconnection.DisconnectedAt.Format("2006-01-02 15:04:05"))

	// Enviar notificación al servicio de notificaciones
	return s.sendSensorAlertToNotificationService(disconnection)
}

// sendSensorAlertToNotificationService envía una alerta de sensor al servicio de notificaciones
func (s *SensorMonitoringService) sendSensorAlertToNotificationService(disconnection models.DisconnectedSensor) error {
	// Crear la estructura de datos para la nueva API de notificaciones
	alertData := map[string]interface{}{
		"sensor_id":   disconnection.SensorID,
		"sensor_data": fmt.Sprintf("Última actividad: %s, Total reportes: %d", disconnection.LastSeen.Format("2006-01-02 15:04:05"), disconnection.TotalReports),
		"building_id": 1, // Default building ID, podría obtenerse del contexto
		"asset_id":    1, // Default asset ID, podría obtenerse del contexto
		"message": fmt.Sprintf("El sensor %s se ha desconectado. Última actividad: %s, tiempo sin reportar: %s",
			disconnection.SensorID,
			disconnection.LastSeen.Format("2006-01-02 15:04:05"),
			time.Since(disconnection.LastSeen).Round(time.Minute).String()),
		"alert_type":        "sensor_disconnected",
		"tipo":              "sensor", // Tipo de notificación
		"prioridad":         "high",   // Prioridad como string
		"notification_mail": true,
		"notification_sms":  false,
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(alertData)
	if err != nil {
		return fmt.Errorf("error al serializar datos de alerta: %v", err)
	}

	// URL del servicio de notificaciones
	notificationURL := "http://notification-service:8091/sensor/alert"

	// Si hay URL de notificación configurada en config, usarla
	if s.notificationURL != "" {
		notificationURL = s.notificationURL
	}

	// Crear request HTTP
	req, err := http.NewRequest("POST", notificationURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear request HTTP: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar notificación HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en respuesta del servicio de notificaciones: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ Alerta de sensor enviada exitosamente al servicio de notificaciones")
	return nil
}

// sendHTTPNotification envía una notificación HTTP al servicio de notificaciones
func (s *SensorMonitoringService) sendHTTPNotification(disconnection models.DisconnectedSensor) error {
	// Preparar payload para el servicio de notificaciones
	payload := map[string]interface{}{
		"type":      "sensor_disconnection",
		"sensor_id": disconnection.SensorID,
		"message": fmt.Sprintf("El sensor %s se ha desconectado. Última actividad: %v",
			disconnection.SensorID, disconnection.LastSeen.Format("2006-01-02 15:04:05")),
		"last_seen":       disconnection.LastSeen,
		"disconnected_at": disconnection.DisconnectedAt,
		"total_reports":   disconnection.TotalReports,
		"severity":        "warning",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar notificación: %v", err)
	}

	// Enviar request HTTP
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(s.notificationURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al enviar notificación HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ Notificación HTTP enviada exitosamente para sensor %s", disconnection.SensorID)
	} else {
		return fmt.Errorf("notificación HTTP falló con status %d", resp.StatusCode)
	}

	return nil
}

// GetAllSensorStatus obtiene el estado de todos los sensores
func (s *SensorMonitoringService) GetAllSensorStatus() ([]models.SensorStatus, error) {
	return s.sensorStatusRepo.GetAllSensors()
}

// GetSensorStatus obtiene el estado de un sensor específico
func (s *SensorMonitoringService) GetSensorStatus(sensorID string) (*models.SensorStatus, error) {
	return s.sensorStatusRepo.GetSensorStatus(sensorID)
}

// StartMonitoring inicia el monitoreo automático de sensores
func (s *SensorMonitoringService) StartMonitoring(checkInterval time.Duration) {
	log.Printf("🚀 Iniciando monitoreo automático de sensores (intervalo: %v, timeout: %d min)",
		checkInterval, s.timeoutMinutes)

	ticker := time.NewTicker(checkInterval)
	go func() {
		for range ticker.C {
			err := s.CheckDisconnectedSensors()
			if err != nil {
				log.Printf("❌ Error en verificación de sensores: %v", err)
			}
		}
	}()
}

// GetSensorStats obtiene estadísticas generales de sensores
func (s *SensorMonitoringService) GetSensorStats() (map[string]interface{}, error) {
	sensors, err := s.sensorStatusRepo.GetAllSensors()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_sensors":    len(sensors),
		"active_sensors":   0,
		"inactive_sensors": 0,
		"last_check":       time.Now(),
	}

	for _, sensor := range sensors {
		if sensor.IsActive {
			stats["active_sensors"] = stats["active_sensors"].(int) + 1
		} else {
			stats["inactive_sensors"] = stats["inactive_sensors"].(int) + 1
		}
	}

	return stats, nil
}
