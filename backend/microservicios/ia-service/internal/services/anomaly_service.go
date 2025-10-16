package services

import (
	"bytes"
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
	mlEngineURL          string
	fetchIntervalMinutes int
	defaultActivoID      int
}

// NewAnomalyService crea una nueva instancia del servicio
func NewAnomalyService(repo *repository.AnomalyRepository, parserServiceURL, mlEngineURL string, fetchIntervalMinutes int, defaultActivoID int) *AnomalyService {
	return &AnomalyService{
		repo:                 repo,
		client:               &http.Client{Timeout: 60 * time.Second},
		parserServiceURL:     parserServiceURL,
		mlEngineURL:          mlEngineURL,
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
		totalLecturas += len(sensor.Datos)
	}

	log.Printf("✅ Datos recibidos: %d sensores, %d lecturas totales", len(parserResp.Sensores), totalLecturas)

	// Transformar datos y enviar al ML Engine
	s.processDataWithMLEngine(parserResp)
}

// transformToMLFormat transforma los datos del ParserService al formato del ML Engine
// ParserService agrupa por sensor, ML Engine espera registros agrupados por timestamp
func (s *AnomalyService) transformToMLFormat(data models.ParserServiceResponse) []models.MLSensorRecord {
	// Mapa para agrupar por timestamp: timestamp -> map[sensor_id]valor
	timestampMap := make(map[string]map[string]float64)

	// Agrupar datos por timestamp
	for _, sensor := range data.Sensores {
		for _, dataPoint := range sensor.Datos {
			timestampStr := dataPoint.Timestamp.Format(time.RFC3339)

			if _, exists := timestampMap[timestampStr]; !exists {
				timestampMap[timestampStr] = make(map[string]float64)
			}

			// Agregar el valor del sensor al timestamp
			timestampMap[timestampStr][sensor.IDSensor] = dataPoint.Valor
		}
	}

	// Convertir el mapa a lista de MLSensorRecord
	var mlRecords []models.MLSensorRecord
	for timestamp, sensorValues := range timestampMap {
		record := models.MLSensorRecord{
			Timestamp:      timestamp,
			AACRMotPV:      sensorValues["A_ACR_Mot.PV"],
			AACRMotSV:      sensorValues["A_ACR_Mot.SV"],
			ADescBomPV:     sensorValues["A_Desc_Bom.PV"],
			ADescBomSV:     sensorValues["A_Desc_Bom.SV"],
			FlujoCaudal:    sensorValues["Flujo_Caudal"],
			NivelPresion:   sensorValues["Nivel_Presion"],
			PotenciaActiva: sensorValues["Potencia_Activa"],
			TempAgua:       sensorValues["Temp_Agua"],
			TempAmbiente:   sensorValues["Temp_Ambiente"],
			Vibracion:      sensorValues["Vibracion"],
		}
		mlRecords = append(mlRecords, record)
	}

	log.Printf("📊 Transformados %d registros agrupados por timestamp", len(mlRecords))
	return mlRecords
} // processDataWithMLEngine procesa los datos enviándolos al ML Engine
func (s *AnomalyService) processDataWithMLEngine(data models.ParserServiceResponse) {
	// Transformar datos al formato del ML Engine
	mlRecords := s.transformToMLFormat(data)

	if len(mlRecords) < 180 {
		log.Printf("⚠️  Insuficientes registros (%d). El ML Engine requiere al menos 180", len(mlRecords))
		return
	}

	// Crear la petición para el ML Engine
	mlRequest := models.MLPredictionRequest{
		Pump:    fmt.Sprintf("activo_%d", data.ActivoID),
		Records: mlRecords,
	}

	// Enviar al ML Engine
	predictions, err := s.sendToMLEngine(mlRequest)
	if err != nil {
		log.Printf("❌ Error al obtener predicciones del ML Engine: %v", err)
		return
	}

	// Guardar anomalías detectadas
	savedCount := 0
	for _, result := range predictions.Results {
		if result.IsAnomaly == 1 {
			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, result.Timestamp)
			if err != nil {
				log.Printf("⚠️  Error parseando timestamp %s: %v", result.Timestamp, err)
				continue
			}

			anomaly := &models.StoreAnomalyRequest{
				ActivoID:          data.ActivoID,
				SensorID:          "combined", // Anomalía basada en múltiples sensores
				Timestamp:         timestamp,
				AnomalyScore:      result.AnomalyScore,
				AnomalyLikelihood: result.AnomalyLikelihood,
				Severidad:         result.Severity,
				Descripcion:       result.Description,
				Threshold:         result.Threshold,
				IsAnomaly:         result.IsAnomaly,
			}

			if _, err := s.SaveAnomaly(anomaly); err != nil {
				log.Printf("⚠️  Error guardando anomalía: %v", err)
			} else {
				savedCount++
			}
		}
	}

	log.Printf("✅ ML Engine procesó %d registros, detectó %d anomalías",
		len(predictions.Results), savedCount)
}

// sendToMLEngine envía los datos al ML Engine y retorna las predicciones
func (s *AnomalyService) sendToMLEngine(request models.MLPredictionRequest) (*models.MLPredictionResponse, error) {
	url := fmt.Sprintf("%s/predict_anomaly", s.mlEngineURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	log.Printf("🤖 Enviando %d registros al ML Engine: %s", len(request.Records), url)

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML Engine returned status %d: %s", resp.StatusCode, string(body))
	}

	var mlResponse models.MLPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &mlResponse, nil
}

// SaveAnomaly guarda una anomalía manualmente (para el endpoint POST)
func (s *AnomalyService) SaveAnomaly(anomaly *models.StoreAnomalyRequest) (*models.Anomaly, error) {
	return s.repo.SaveAnomaly(anomaly)
} // GetLatestAnomaly obtiene la anomalía más reciente
func (s *AnomalyService) GetLatestAnomaly() (*models.Anomaly, error) {
	return s.repo.GetLatestAnomaly()
}

// GetAnomaliesBySensor obtiene anomalías por sensor
func (s *AnomalyService) GetAnomaliesBySensor(sensorID string, limit int) ([]models.Anomaly, error) {
	return s.repo.GetAnomaliesBySensor(sensorID, limit)
}

// GetAnomaliesByActivo obtiene anomalías por activo
func (s *AnomalyService) GetAnomaliesByActivo(activoID, limit, offset int) ([]models.Anomaly, error) {
	return s.repo.GetAnomaliesByActivo(activoID, limit, offset)
}
