package services

import (
	"bytes"
	"context"
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

// Estructura para los datos del template
type ReporteData struct {
	Activo          *models.Activo
	Edificio        *models.Edificio
	UltimaAccion    *models.AccionMantenimiento
	Acciones        []models.AccionMantenimiento
	FechaGeneracion string
}

// Generar reporte PDF por activo usando HTML template
func (s *ReporteService) GenerarReportePDFPorActivo(activoID int) ([]byte, string, error) {
	fmt.Printf("=== GENERANDO REPORTE PARA ACTIVO ID: %d ===\n", activoID)

	// Obtener datos del activo
	activo, err := s.activoRepo.GetByID(activoID)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo obtener el activo: %v", err)
	}
	fmt.Printf("ACTIVO ENCONTRADO: ID=%d, Nombre=%s, Tipo=%s, Estado=%s\n",
		activo.ID, activo.Nombre, activo.Tipo, activo.Estado)

	// Obtener datos del edificio si existe
	var edificio *models.Edificio
	if activo.EdificioID != 0 {
		edificio, err = s.edificioRepo.GetByID(activo.EdificioID)
		if err != nil {
			// No es error crítico si no se encuentra el edificio
			fmt.Printf("Advertencia: no se pudo obtener el edificio: %v\n", err)
		} else {
			fmt.Printf("EDIFICIO ENCONTRADO: ID=%d, Nombre=%s, Direccion=%s\n",
				edificio.ID, edificio.Nombre, edificio.Direccion)
		}
	} else {
		fmt.Printf("ACTIVO SIN EDIFICIO ASIGNADO\n")
	}

	// Obtener última acción de mantenimiento
	ultimaAccion, err := s.accionRepo.GetUltimaAccionPorActivo(activoID)
	if err != nil {
		// No es error crítico si no hay acciones
		fmt.Printf("Advertencia: no se encontraron acciones de mantenimiento: %v\n", err)
	} else {
		fmt.Printf("ULTIMA ACCION ENCONTRADA: ID=%d, Tipo=%s, Estado=%s, Descripcion=%s\n",
			ultimaAccion.ID, ultimaAccion.Tipo, ultimaAccion.Estado, ultimaAccion.Descripcion)
	}

	// Obtener historial de acciones
	accionesCompletas, err := s.accionRepo.GetByActivo(activoID)
	if err != nil {
		// No es error crítico si no hay acciones
		fmt.Printf("Advertencia: no se pudo obtener historial de acciones: %v\n", err)
		accionesCompletas = []models.AccionConDetalles{}
	}

	// Convertir AccionConDetalles a AccionMantenimiento para el template
	acciones := make([]models.AccionMantenimiento, len(accionesCompletas))
	for i, accion := range accionesCompletas {
		acciones[i] = accion.AccionMantenimiento
	}
	fmt.Printf("ACCIONES ENCONTRADAS: %d acciones\n", len(acciones))

	// Preparar datos para el template
	data := ReporteData{
		Activo:          activo,
		Edificio:        edificio,
		UltimaAccion:    ultimaAccion,
		Acciones:        acciones,
		FechaGeneracion: time.Now().Format("02/01/2006 15:04:05"),
	}
	fmt.Printf("DATOS PREPARADOS PARA TEMPLATE\n")

	// Renderizar HTML template
	htmlContent, err := s.renderHTMLTemplate(data)
	if err != nil {
		return nil, "", fmt.Errorf("error renderizando template HTML: %v", err)
	}

	fmt.Printf("HTML generado: %d caracteres\n", len(htmlContent))

	// Guardar HTML para debug - GUARDAR HTML COMPLETO EN ARCHIVO TEMPORAL
	if err := os.WriteFile("/tmp/debug_report.html", []byte(htmlContent), 0644); err == nil {
		fmt.Printf("HTML guardado en /tmp/debug_report.html para inspección\n")
	}

	// Guardar HTML para debug
	if len(htmlContent) > 100 {
		fmt.Printf("PRIMEROS 200 CARACTERES DEL HTML: %s\n", htmlContent[:200])
	}

	// Generar PDF desde HTML
	pdfBytes, err := s.generatePDFFromHTML(htmlContent)
	if err != nil {
		return nil, "", fmt.Errorf("error generando PDF: %v", err)
	}

	// Crear registro del reporte
	reporte := &models.Reporte{
		ActivoID:    activoID,
		TipoReporte: "PDF_ACTIVO",
		Contenido:   "Reporte PDF generado desde template HTML",
		GeneradoEn:  time.Now(),
		Estado:      "generado",
	}

	_, err = s.repo.Create(reporte)
	if err != nil {
		// Log del error pero no falla la generación del PDF
		fmt.Printf("Advertencia: no se pudo crear el registro del reporte: %v\n", err)
	}

	filename := fmt.Sprintf("reporte_activo_%d_%s.pdf", activo.ID, time.Now().Format("20060102_150405"))
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
- ID: %d
- Nombre: %s
- Tipo: %s
- Estado: %s
- Ubicación: %s

Información del Edificio:
- Nombre: %s
- Dirección: %s

`, activo.ID, activo.Nombre, activo.Tipo, activo.Estado, activo.Ubicacion,
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

// Renderizar template HTML con los datos
func (s *ReporteService) renderHTMLTemplate(data ReporteData) (string, error) {
	templatePath := filepath.Join("templates", "reporte_activo.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error cargando template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("error ejecutando template: %v", err)
	}

	return buf.String(), nil
}

// Generar PDF desde HTML usando chromedp
func (s *ReporteService) generatePDFFromHTML(htmlContent string) ([]byte, error) {
	// Crear contexto con opciones para headless browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Guardar HTML en archivo temporal para evitar problemas con data URLs largas
	tmpFile := "/tmp/temp_report.html"
	if err := os.WriteFile(tmpFile, []byte(htmlContent), 0644); err != nil {
		return nil, fmt.Errorf("error guardando HTML temporal: %v", err)
	}
	defer os.Remove(tmpFile)

	// Variable para almacenar el PDF
	var pdfBytes []byte

	// Navegar al archivo HTML y generar PDF
	fileURL := "file://" + tmpFile
	fmt.Printf("Navegando a: %s\n", fileURL)

	err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.Sleep(5*time.Second), // Aumentar tiempo de espera considerablemente
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				WithScale(0.9). // Aumentar escala
				WithDisplayHeaderFooter(false).
				WithPreferCSSPageSize(false).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("error generando PDF con chromedp: %v", err)
	}

	// Validar que se generó correctamente (reducir umbral)
	fmt.Printf("PDF generado con %d bytes\n", len(pdfBytes))
	if len(pdfBytes) < 500 {
		return nil, fmt.Errorf("PDF generado es muy pequeño: %d bytes", len(pdfBytes))
	}

	return pdfBytes, nil
}
