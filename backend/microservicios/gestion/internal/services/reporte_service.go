package services

import (
	"bytes"
	"fmt"
	"time"
	"github.com/jung-kurt/gofpdf/v2"
	"gestion/internal/models"
	"gestion/internal/repository"
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

// Generar PDF con información completa del activo
func (s *ReporteService) generarPDF(activo *models.Activo, edificio *models.Edificio, ultimaAccion *models.AccionMantenimiento, contenido string) ([]byte, error) {
	// Crear nuevo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configurar fuentes
	pdf.SetFont("Arial", "B", 20)
	
	// Header del documento
	pdf.SetTextColor(0, 51, 102) // Azul corporativo
	pdf.CellFormat(0, 15, "REPORTE DE ACTIVO - InspectAR", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Información del activo
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "INFORMACION DEL ACTIVO", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "ID del Activo:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, activo.ActivoID, "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Nombre:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, activo.Nombre, "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Tipo:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, activo.Tipo, "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Estado:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	// Colorear el estado según su valor
	switch activo.Estado {
	case "operativo":
		pdf.SetTextColor(0, 128, 0) // Verde
	case "mantenimiento":
		pdf.SetTextColor(255, 140, 0) // Naranja
	case "fuera_servicio":
		pdf.SetTextColor(255, 0, 0) // Rojo
	default:
		pdf.SetTextColor(0, 0, 0) // Negro
	}
	pdf.CellFormat(0, 8, activo.Estado, "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0) // Resetear color
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Ubicacion:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, activo.Ubicacion, "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Creado en:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, activo.CreadoEn.Format("2006-01-02 15:04:05"), "", 1, "L", false, 0, "")
	
	pdf.Ln(10)

	// Información del edificio
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "INFORMACION DEL EDIFICIO", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Nombre:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, edificio.Nombre, "", 1, "L", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Direccion:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, edificio.Direccion, "", 1, "L", false, 0, "")
	
	pdf.Ln(10)

	// Información de mantenimiento
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "HISTORIAL DE MANTENIMIENTO", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	if ultimaAccion != nil {
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(40, 8, "Ultimo Tipo:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, ultimaAccion.Tipo, "", 1, "L", false, 0, "")
		
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(40, 8, "Descripcion:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 12)
		// Texto multilinea para descripción
		pdf.CellFormat(0, 8, ultimaAccion.Descripcion, "", 1, "L", false, 0, "")
		
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(40, 8, "Estado:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 12)
		// Colorear estado de la acción
		switch ultimaAccion.Estado {
		case "completado":
			pdf.SetTextColor(0, 128, 0) // Verde
		case "en_progreso":
			pdf.SetTextColor(255, 140, 0) // Naranja
		case "pendiente":
			pdf.SetTextColor(255, 0, 0) // Rojo
		default:
			pdf.SetTextColor(0, 0, 0) // Negro
		}
		pdf.CellFormat(0, 8, ultimaAccion.Estado, "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Resetear color
		
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(40, 8, "Prioridad:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 12)
		// Colorear prioridad
		switch ultimaAccion.Prioridad {
		case "critica":
			pdf.SetTextColor(128, 0, 128) // Púrpura
		case "alta":
			pdf.SetTextColor(255, 0, 0) // Rojo
		case "media":
			pdf.SetTextColor(255, 140, 0) // Naranja
		case "baja":
			pdf.SetTextColor(0, 128, 0) // Verde
		default:
			pdf.SetTextColor(0, 0, 0) // Negro
		}
		pdf.CellFormat(0, 8, ultimaAccion.Prioridad, "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Resetear color
		
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(40, 8, "Fecha Inicio:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, ultimaAccion.FechaInicio.Format("2006-01-02 15:04:05"), "", 1, "L", false, 0, "")
	} else {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(128, 128, 128) // Gris
		pdf.CellFormat(0, 8, "No hay acciones de mantenimiento registradas para este activo.", "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Resetear color
	}
	
	pdf.Ln(15)

	// Obtener acciones adicionales
	acciones, err := s.accionRepo.GetByActivo(activo.ID)
	if err == nil && len(acciones) > 1 {
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, "HISTORIAL COMPLETO DE ACCIONES", "", 1, "L", false, 0, "")
		pdf.Ln(3)
		
		// Mostrar hasta 5 acciones más recientes para no sobrecargar el PDF
		maxAcciones := len(acciones)
		if maxAcciones > 5 {
			maxAcciones = 5
		}
		
		for i := 0; i < maxAcciones; i++ {
			accion := acciones[i]
			pdf.SetFont("Arial", "B", 11)
			pdf.CellFormat(0, 6, fmt.Sprintf("%d. %s - %s", i+1, accion.Tipo, accion.Estado), "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.CellFormat(5, 5, "", "", 0, "L", false, 0, "") // Indentación
			pdf.CellFormat(0, 5, fmt.Sprintf("Fecha: %s | Prioridad: %s", 
				accion.FechaInicio.Format("2006-01-02"), accion.Prioridad), "", 1, "L", false, 0, "")
			pdf.Ln(2)
		}
	}

	// Agregar información estadística
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "ESTADISTICAS DEL ACTIVO", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// Obtener estadísticas básicas
	totalAcciones, _ := s.accionRepo.GetByActivo(activo.ID)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(60, 8, "Total de acciones registradas:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("%d", len(totalAcciones)), "", 1, "L", false, 0, "")

	// Footer con información de generación
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, fmt.Sprintf("Reporte generado automaticamente el: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "InspectAR - Sistema de Gestion de Activos Industriales", "", 1, "C", false, 0, "")

	// Convertir PDF a bytes
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %v", err)
	}
	
	return buf.Bytes(), nil
}
