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

// Crear técnico con especialidades
func (s *TecnicoService) CrearTecnico(req *models.CreateTecnicoRequest) (*models.Tecnico, error) {
	return s.repo.Create(req)
}

// Listar todos los técnicos
func (s *TecnicoService) ListarTodosLosTecnicos() ([]models.Tecnico, error) {
	return s.repo.GetAll()
}

// Listar técnicos relacionados con un activo específico
func (s *TecnicoService) ListarTecnicosPorActivo(activoID int, soloAutorizados bool) ([]models.Tecnico, error) {
	return s.repo.GetByActivo(activoID, soloAutorizados)
}

// Listar técnicos relacionados con un edificio (indirectamente a través de activos)
func (s *TecnicoService) ListarTecnicosPorEdificio(edificioID int, soloAutorizados bool) ([]models.Tecnico, error) {
	return s.repo.GetByEdificio(edificioID, soloAutorizados)
}

// Consultar información de un técnico específico
func (s *TecnicoService) ObtenerTecnico(id int) (*models.Tecnico, error) {
	return s.repo.GetByID(id)
}

// Actualizar estado autorizado de técnicos
func (s *TecnicoService) ActualizarAutorizado(id int, autorizado bool) error {
	return s.repo.UpdateAutorizado(id, autorizado)
}

// Asignar técnico a activo (relación muchos a muchos)
func (s *TecnicoService) AsignarTecnicoAActivo(activoID, tecnicoID int) error {
	return s.repo.AsignarTecnicoAActivo(activoID, tecnicoID)
}

// Obtener activos asociados a un técnico (RUTA PRINCIPAL)
func (s *TecnicoService) ObtenerActivosPorTecnico(tecnicoID int) ([]models.Activo, error) {
	return s.repo.ObtenerActivosPorTecnico(tecnicoID)
}
