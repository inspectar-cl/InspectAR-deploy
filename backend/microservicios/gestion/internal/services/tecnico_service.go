package services

import (
	"gestion/internal/models"
	"gestion/internal/repository"
)

type TecnicoService struct {
	repo *repository.TecnicoRepository
}

func NewTecnicoService(repo *repository.TecnicoRepository) *TecnicoService {
	return &TecnicoService{repo: repo}
}

func (s *TecnicoService) CrearTecnico(req *models.CreateTecnicoRequest) (*models.Tecnico, error) {
	return s.repo.Create(req)
}

func (s *TecnicoService) ListarTecnicos() ([]models.Tecnico, error) {
	return s.repo.GetAll()
}

func (s *TecnicoService) ObtenerTecnico(id int) (*models.Tecnico, error) {
	return s.repo.GetByID(id)
}

func (s *TecnicoService) ActualizarDisponibilidad(id int, disponible bool) error {
	return s.repo.UpdateDisponibilidad(id, disponible)
}
