package services

import (
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
)

type ActivoService struct {
	activoRepo   *repository.ActivoRepository
	edificioRepo *repository.EdificioRepository
}

func NewActivoService(activoRepo *repository.ActivoRepository, edificioRepo *repository.EdificioRepository) *ActivoService {
	return &ActivoService{
		activoRepo:   activoRepo,
		edificioRepo: edificioRepo,
	}
}

// GetActivoByID obtiene un activo por su ID
func (s *ActivoService) GetActivoByID(id int) (*models.Activo, error) {
	return s.activoRepo.GetByID(id)
}

// GetAllActivos obtiene todos los activos
func (s *ActivoService) GetAllActivos() ([]models.Activo, error) {
	return s.activoRepo.GetAll()
}

// GetActivosByEdificio obtiene todos los activos de un edificio específico
func (s *ActivoService) GetActivosByEdificio(edificioID int) ([]models.Activo, error) {
	// Verificar que el edificio existe
	_, err := s.edificioRepo.GetByID(edificioID)
	if err != nil {
		return nil, fmt.Errorf("edificio no encontrado: %v", err)
	}

	return s.activoRepo.GetByEdificio(edificioID)
}

// GetActivosByTipo obtiene todos los activos de un tipo específico
func (s *ActivoService) GetActivosByTipo(tipo string) ([]models.Activo, error) {
	// Validar que el tipo es válido
	if !models.ValidarTipoActivo(tipo) {
		return nil, fmt.Errorf("tipo de activo inválido: %s. Tipos permitidos: %v", tipo, models.TiposActivosValidos)
	}

	return s.activoRepo.GetByTipo(tipo)
}
