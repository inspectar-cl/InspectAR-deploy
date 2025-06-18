// services/activo_service.go
package services

import (
	"ParserService/internal/models"
	"ParserService/internal/repository"
	"context"
	"errors"
)

type ActivoService struct {
	activoRepo *repository.ActivoRepository
}

func NewActivoService(repo *repository.ActivoRepository) *ActivoService {
	return &ActivoService{activoRepo: repo}
}

func (s *ActivoService) CrearActivo(ctx context.Context, activo *models.Activo) (string, error) {
	existe, err := s.activoRepo.ExistActivo(ctx, activo.ActivoID)
	if err != nil {
		return "", err
	}
	if existe {
		return "", errors.New("El activo ya existe")
	}
	return s.activoRepo.CreateActivo(ctx, activo)
}

func (s *ActivoService) ObtenerActivo(ctx context.Context, activoID string) (*models.Activo, error) {
	return s.activoRepo.GetActivo(ctx, activoID)
}

func (s *ActivoService) AgregarSensor(ctx context.Context, activoID string, sensor models.Sensor) error {
	return s.activoRepo.AddSensor(ctx, activoID, sensor)
}

func (s *ActivoService) GetAllActivos(ctx context.Context) ([]models.Activo, error) {
	return s.activoRepo.GetAllActivos(ctx)
}

func (s *ActivoService) ActualizarEstado(ctx context.Context, activoID string, nuevoEstado string) error {
	return s.activoRepo.ActualizarEstado(ctx, activoID, nuevoEstado)
}