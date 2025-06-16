// services/sensor_service.go
package services

import (
	"context"
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
	return s.influxRepo.GetSensorData(ctx, sensorID)
}

func (s *SensorService) InsertarLectura(ctx context.Context, sensorID string, valor float64, timestamp time.Time) error {
	return s.influxRepo.InsertSensorData(ctx, sensorID, valor, timestamp)
}
