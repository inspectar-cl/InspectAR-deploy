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

// Crear acciones de mantenimiento (preventivo, correctivo, emergencia)
func (s *AccionMantenimientoService) CrearAccion(req *models.CreateAccionMantenimientoRequest) (*models.AccionMantenimiento, error) {
	return s.repo.Create(req)
}

// Obtener acciones donde el técnico puede buscar sus acciones
func (s *AccionMantenimientoService) ObtenerAccionesPorTecnico(tecnicoID int) ([]models.AccionConDetalles, error) {
	return s.repo.GetByTecnico(tecnicoID)
}

// Consultar acciones por activo
func (s *AccionMantenimientoService) ObtenerAccionesPorActivo(activoID int) ([]models.AccionConDetalles, error) {
	return s.repo.GetByActivo(activoID)
}

// Actualizar estado de acciones (pendiente, en_progreso, completado)
func (s *AccionMantenimientoService) ActualizarEstado(id int, estado string) error {
	return s.repo.UpdateEstado(id, estado)
}

// Listar acciones pendientes con prioridad
func (s *AccionMantenimientoService) ObtenerAccionesPendientes() ([]models.AccionConDetalles, error) {
	return s.repo.GetPendientesConPrioridad()
}
