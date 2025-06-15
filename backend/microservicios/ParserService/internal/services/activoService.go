// services/activo_service.go
package services

import (
	"context"
	"forms/internal/models"
	"forms/internal/repository"
)

type ActivoService struct {
	activoRepo *repository.ActivoRepository
}

func NewActivoService(repo *repository.ActivoRepository) *ActivoService {
	return &ActivoService{activoRepo: repo}
}

func (s *ActivoService) CrearActivo(ctx context.Context, activo *models.Activo) (string, error) {
	return s.activoRepo.CreateActivo(ctx, activo)
}

func (s *ActivoService) ObtenerActivo(ctx context.Context, activoID string) (*models.Activo, error) {
	return s.activoRepo.GetActivo(ctx, activoID)
}

func (s *ActivoService) AgregarSensor(ctx context.Context, activoID string, sensor models.Sensor) error {
	return s.activoRepo.AddSensor(ctx, activoID, sensor)
}
