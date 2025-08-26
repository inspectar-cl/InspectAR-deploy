package services

import (
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
	"time"
)

type ReporteService struct {
	repo         *repository.ReporteRepository
	activoRepo   *repository.ActivoRepository
	accionRepo   *repository.AccionMantenimientoRepository
	edificioRepo *repository.EdificioRepository
}

func NewReporteService(
	repo *repository.ReporteRepository,
	activoRepo *repository.ActivoRepository,
	accionRepo *repository.AccionMantenimientoRepository,
	edificioRepo *repository.EdificioRepository,
) *ReporteService {
	return &ReporteService{
		repo:         repo,
		activoRepo:   activoRepo,
		accionRepo:   accionRepo,
		edificioRepo: edificioRepo,
	}
}

// Generar reporte PDF por activo
func (s *ReporteService) GenerarReportePDFPorActivo(activoID int) ([]byte, string, error) {
	// Obtener datos del activo
	activo, err := s.activoRepo.GetByID(activoID)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo obtener el activo: %v", err)
	}

	// Obtener datos del edificio
	edificio, err := s.edificioRepo.GetByID(activo.EdificioID)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo obtener el edificio: %v", err)
	}

	// Obtener última acción de mantenimiento
	ultimaAccion, err := s.accionRepo.GetUltimaAccionPorActivo(activoID)
	if err != nil {
		// No es error crítico si no hay acciones
		ultimaAccion = nil
	}

	// Crear el contenido del reporte
	contenido := s.generarContenidoReporte(activo, edificio, ultimaAccion)

	// Crear registro del reporte en la base de datos
	reporte := &models.Reporte{
		ActivoID:    activoID,
		TipoReporte: "mantenimiento",
		Contenido:   contenido,
		GeneradoEn:  time.Now(),
		Estado:      "generado",
	}

	_, err = s.repo.Create(reporte)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo guardar el reporte: %v", err)
	}

	// Generar PDF
	pdfBytes, err := s.generarPDF(activo, edificio, ultimaAccion, contenido)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo generar el PDF: %v", err)
	}

	filename := fmt.Sprintf("reporte_activo_%s_%s.pdf",
		activo.ActivoID,
		time.Now().Format("2006-01-02_15-04-05"))

	return pdfBytes, filename, nil
}

// Obtener reportes de un activo
func (s *ReporteService) ObtenerReportesPorActivo(activoID int) ([]models.ReporteCompleto, error) {
	return s.repo.GetByActivo(activoID)
}

// Generar contenido del reporte
func (s *ReporteService) generarContenidoReporte(activo *models.Activo, edificio *models.Edificio, ultimaAccion *models.AccionMantenimiento) string {
	contenido := fmt.Sprintf(`
REPORTE DE ACTIVO

Información del Activo:
- ID: %s
- Nombre: %s
- Tipo: %s
- Estado: %s
- Ubicación: %s

Información del Edificio:
- Nombre: %s
- Dirección: %s

`, activo.ActivoID, activo.Nombre, activo.Tipo, activo.Estado, activo.Ubicacion,
		edificio.Nombre, edificio.Direccion)

	if ultimaAccion != nil {
		contenido += fmt.Sprintf(`
Última Acción de Mantenimiento:
- Tipo: %s
- Descripción: %s
- Estado: %s
- Prioridad: %s
- Fecha: %s

`, ultimaAccion.Tipo, ultimaAccion.Descripcion, ultimaAccion.Estado,
			ultimaAccion.Prioridad, ultimaAccion.FechaInicio.Format("2006-01-02 15:04:05"))
	} else {
		contenido += "\nNo hay acciones de mantenimiento registradas para este activo.\n"
	}

	contenido += fmt.Sprintf("Reporte generado el: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return contenido
}

// Generar PDF (placeholder - implementar con librería PDF como gofpdf)
func (s *ReporteService) generarPDF(activo *models.Activo, edificio *models.Edificio, ultimaAccion *models.AccionMantenimiento, contenido string) ([]byte, error) {
	// TODO: Implementar generación real de PDF con una librería como gofpdf
	// Por ahora retornamos el contenido como texto plano
	return []byte(contenido), nil
}
