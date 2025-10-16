package services

import (
	"encoding/json"
	"fmt"
	"ia-service/internal/models"
	"ia-service/internal/repository"
	"io"
	"log"
	"net/http"
	"time"
)

type AnomalyService struct {
	repo                 *repository.AnomalyRepository
	client               *http.Client
	parserServiceURL     string
	fetchIntervalMinutes int
	defaultActivoID      int
}

// NewAnomalyService crea una nueva instancia del servicio
func NewAnomalyService(repo *repository.AnomalyRepository, parserServiceURL string, fetchIntervalMinutes int, defaultActivoID int) *AnomalyService {
	return &AnomalyService{
		repo:                 repo,
		client:               &http.Client{Timeout: 30 * time.Second},
		parserServiceURL:     parserServiceURL,
		fetchIntervalMinutes: fetchIntervalMinutes,
		defaultActivoID:      defaultActivoID,
	}
}

// StartScheduler inicia el proceso periódico de obtención de datos del ParserService
func (s *AnomalyService) StartScheduler() {
	ticker := time.NewTicker(time.Duration(s.fetchIntervalMinutes) * time.Minute)

	log.Printf("🔄 Scheduler iniciado: Obteniendo datos cada %d minutos para activo_id=%d",
		s.fetchIntervalMinutes, s.defaultActivoID)

	// Ejecutar inmediatamente al inicio
	go s.fetchAndProcessData()

	// Ejecutar periódicamente
	go func() {
		for range ticker.C {
			s.fetchAndProcessData()
		}
	}()
}

// fetchAndProcessData obtiene 1000 datos del ParserService y los procesa
func (s *AnomalyService) fetchAndProcessData() {
	activoID := s.defaultActivoID
	url := fmt.Sprintf("%s/lectura/%d/window?page=1&limit=1000", s.parserServiceURL, activoID)

	log.Printf("📡 Solicitando datos del ParserService: %s", url)

	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("❌ Error al solicitar datos del ParserService: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ ParserService retornó status %d: %s", resp.StatusCode, string(body))
		return
	}

	var parserResp models.ParserServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&parserResp); err != nil {
		log.Printf("❌ Error decodificando respuesta del ParserService: %v", err)
		return
	}

	totalLecturas := 0
	for _, sensor := range parserResp.Sensores {
		totalLecturas += len(sensor.Lecturas)
	}

	log.Printf("✅ Datos recibidos: %d sensores, %d lecturas totales", len(parserResp.Sensores), totalLecturas)

	// Aquí puedes agregar la lógica de detección de anomalías
	// Por ahora, solo registramos los datos recibidos
	s.processDataForAnomalies(parserResp)
}

// processDataForAnomalies procesa los datos para detectar anomalías
// NOTA: Esta es una implementación simple. En producción, aquí iría tu modelo de ML
func (s *AnomalyService) processDataForAnomalies(data models.ParserServiceResponse) {
	// Ejemplo simple: detectar valores extremos como anomalías
	// En producción, aquí usarías tu modelo de ML real

	threshold := 100.0 // Ejemplo de umbral

	for _, sensor := range data.Sensores {
		for _, lectura := range sensor.Lecturas {
			// Simulación simple: valores > threshold son anomalías
			if lectura.Valor > threshold {
				score := (lectura.Valor - threshold) / threshold
				if score > 1.0 {
					score = 1.0
				}

				severidad := s.calculateSeverity(score)

				anomaly := &models.AnomalyInput{
					ActivoID:          data.ActivoID,
					SensorID:          sensor.SensorID,
					Timestamp:         lectura.Timestamp,
					AnomalyScore:      score,
					AnomalyLikelihood: 0.75, // Ejemplo fijo
					Severidad:         severidad,
					Descripcion:       fmt.Sprintf("Valor %.2f excede umbral %.2f", lectura.Valor, threshold),
					Threshold:         threshold,
					IsAnomaly:         true,
				}

				// Guardar anomalía en la base de datos
				if _, err := s.repo.SaveAnomaly(anomaly); err != nil {
					log.Printf("⚠️  Error guardando anomalía: %v", err)
				}
			}
		}
	}
}

// calculateSeverity calcula la severidad basada en el score
func (s *AnomalyService) calculateSeverity(score float64) string {
	if score >= 0.8 {
		return "critica"
	} else if score >= 0.6 {
		return "alta"
	} else if score >= 0.3 {
		return "media"
	}
	return "baja"
}

// SaveAnomaly guarda una anomalía manualmente (para el endpoint POST)
func (s *AnomalyService) SaveAnomaly(anomaly *models.AnomalyInput) (*models.Anomaly, error) {
	return s.repo.SaveAnomaly(anomaly)
}

// GetLatestAnomaly obtiene la anomalía más reciente
func (s *AnomalyService) GetLatestAnomaly() (*models.Anomaly, error) {
	return s.repo.GetLatestAnomaly()
}

// GetAnomaliesBySensor obtiene anomalías por sensor
func (s *AnomalyService) GetAnomaliesBySensor(sensorID string, limit int) ([]models.Anomaly, error) {
	return s.repo.GetAnomaliesBySensor(sensorID, limit)
}
