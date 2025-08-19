package services

import (
	"gestion/internal/models"
	"gestion/internal/repository"
)

type AccionMantenimientoService struct {
	repo *repository.AccionMantenimientoRepository
}

func NewAccionMantenimientoService(repo *repository.AccionMantenimientoRepository) *AccionMantenimientoService {
	return &AccionMantenimientoService{repo: repo}
}

func (s *AccionMantenimientoService) CrearAccion(req *models.CreateAccionMantenimientoRequest) (*models.AccionMantenimiento, error) {
	return s.repo.Create(req)
}

func (s *AccionMantenimientoService) ObtenerAccionesPorActivo(activoID int) ([]models.AccionMantenimiento, error) {
	return s.repo.GetByActivo(activoID)
}

func (s *AccionMantenimientoService) ActualizarEstado(id int, estado string) error {
	return s.repo.UpdateEstado(id, estado)
}

func (s *AccionMantenimientoService) ObtenerAccionesPendientes() ([]models.AccionMantenimiento, error) {
	return s.repo.GetPendientes()
}
