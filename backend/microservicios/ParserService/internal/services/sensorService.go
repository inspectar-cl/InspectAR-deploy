package services

import (
	"context"
	"log"
	"time"

	"ParserService/internal/models"
	"ParserService/internal/repository"
)

type SensorService struct {
	influxRepo *repository.InfluxRepository
}

func NewSensorService(influxRepo *repository.InfluxRepository) *SensorService {
	return &SensorService{influxRepo: influxRepo}
}

func (s *SensorService) GetDatosSensor(ctx context.Context, sensorID string, since time.Duration) ([]models.SensorData, error) {
	log.Printf("🔍 GetDatosSensor llamado para sensor: %s", sensorID)
	result, err := s.influxRepo.GetSensorData(ctx, sensorID)
	if err != nil {
		log.Printf("❌ Error en GetDatosSensor para %s: %v", sensorID, err)
	} else {
		log.Printf("✅ GetDatosSensor exitoso para %s: %d registros", sensorID, len(result))
	}
	return result, err
}

func (s *SensorService) InsertarLectura(ctx context.Context, sensorID string, valor float64, timestamp time.Time) error {
	return s.influxRepo.InsertSensorData(ctx, sensorID, valor, timestamp)
}

func (s *SensorService) GetSensorLastData(ctx context.Context, sensorID string) (*models.SensorData, error) {
	return s.influxRepo.GetSensorLastData(ctx, sensorID)
}

func (s *SensorService) GetDatosSensorWindow(ctx context.Context, sensorID string, limit int, offset int) ([]models.SensorData, error) {
	log.Printf("🔍 GetDatosSensorWindow llamado para sensor: %s con límite: %d, offset: %d", sensorID, limit, offset)
	result, err := s.influxRepo.GetSensorDataWindow(ctx, sensorID, limit, offset)
	if err != nil {
		log.Printf("❌ Error en GetDatosSensorWindow para %s: %v", sensorID, err)
	} else {
		log.Printf("✅ GetDatosSensorWindow exitoso para %s: %d registros", sensorID, len(result))
	}
	return result, err
}
