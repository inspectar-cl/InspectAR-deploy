// services/sensor_service.go
package services

import (
	"context"
	"time"

	"forms/internal/models"
	"forms/internal/repository"
)

type SensorService struct {
	influxRepo *repository.InfluxRepository
}

func NewSensorService(influxRepo *repository.InfluxRepository) *SensorService {
	return &SensorService{influxRepo: influxRepo}
}

func (s *SensorService) GetDatosSensor(ctx context.Context, sensorID string, since time.Duration) ([]models.SensorData, error) {
	return s.influxRepo.GetSensorData(ctx, sensorID, since)
}
