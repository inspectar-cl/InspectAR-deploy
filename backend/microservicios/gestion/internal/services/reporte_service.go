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
	// Nuevos campos para observaciones
	Observaciones      string
	AutorAnalista      string
	Recomendaciones    []string
	Conclusiones       string
	EstructuraCompleta models.EstructuraInformeReporte
}

// Crear reporte con observaciones
func (s *ReporteService) CrearReporte(req models.CreateReporteRequest) (*models.Reporte, error) {
	// Obtener datos del activo para generar estructura
	activo, err := s.activoRepo.GetByID(req.ActivoID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el activo: %v", err)
	}

	// Generar estructura del informe automáticamente
	estructura := s.generarEstructuraInforme(activo, req.ObservacionesAnalista)

	// Crear metadata
	metadata := map[string]interface{}{
		"fecha_creacion":  time.Now(),
		"tipo_activo":     activo.Tipo,
		"nombre_activo":   activo.Nombre,
		"version":         1,
		"estado_original": activo.Estado,
	}

	// Usar contenido proporcionado o generar uno básico
	contenido := req.Contenido
	if contenido == "" {
		contenido = s.generarContenidoBasico(activo)
	}

	reporte := &models.Reporte{
		ActivoID:              req.ActivoID,
		TipoReporte:           req.TipoReporte,
		Contenido:             contenido,
		ObservacionesAnalista: req.ObservacionesAnalista,
		AutorAnalista:         req.AutorAnalista,
		EstructuraInforme:     estructura,
		MetadataInforme:       metadata,
		VersionReporte:        1,
		EstadoRevision:        "pendiente",
		GeneradoEn:            time.Now(),
		Estado:                "generado",
	}

	return s.repo.Create(reporte)
}

// Actualizar observaciones de un reporte
func (s *ReporteService) ActualizarObservaciones(id int, req models.UpdateObservacionesRequest) error {
	return s.repo.UpdateObservaciones(id, req.ObservacionesAnalista, req.AutorAnalista)
}

// Actualizar estado de revisión
func (s *ReporteService) ActualizarEstadoRevision(id int, req models.UpdateEstadoRevisionRequest) error {
	return s.repo.UpdateEstadoRevision(id, req.EstadoRevision, req.Revisor, req.Observaciones)
}

// Obtener reporte por ID con estructura completa
func (s *ReporteService) ObtenerReportePorID(id int) (*models.Reporte, error) {
	return s.repo.GetByID(id)
}

// Generar estructura del informe automáticamente
func (s *ReporteService) generarEstructuraInforme(activo *models.Activo, observaciones string) map[string]interface{} {
	// Obtener última acción de mantenimiento
	ultimaAccion, _ := s.accionRepo.GetUltimaAccionPorActivo(activo.ID)

	// Obtener historial de acciones
	acciones, _ := s.accionRepo.GetByActivo(activo.ID)

	accionesRealizadas := make([]string, 0)
	for _, accion := range acciones {
		if accion.Estado == "completado" {
			accionesRealizadas = append(accionesRealizadas,
				fmt.Sprintf("%s - %s", accion.Tipo, accion.Descripcion))
		}
	}

	recomendaciones := s.generarRecomendaciones(activo, ultimaAccion)

	estructura := map[string]interface{}{
		"resumen": fmt.Sprintf("Reporte %s para %s ubicado en %s",
			activo.Tipo, activo.Nombre, activo.Ubicacion),
		"observaciones": observaciones,
		"datos_activo": map[string]interface{}{
			"id":        activo.ID,
			"nombre":    activo.Nombre,
			"tipo":      activo.Tipo,
			"estado":    activo.Estado,
			"ubicacion": activo.Ubicacion,
		},
		"acciones_realizadas": accionesRealizadas,
		"recomendaciones":     recomendaciones,
		"conclusiones":        s.generarConclusiones(activo, observaciones),
		"fecha_generacion":    time.Now().Format("2006-01-02 15:04:05"),
	}

	return estructura
}

// Generar recomendaciones automáticas
func (s *ReporteService) generarRecomendaciones(activo *models.Activo, ultimaAccion *models.AccionMantenimiento) []string {
	recomendaciones := make([]string, 0)

	// Recomendaciones basadas en el tipo de activo
	switch activo.Tipo {
	case "caldera":
		recomendaciones = append(recomendaciones,
			"Realizar inspección visual mensual de conexiones",
			"Verificar presión de operación semanalmente",
			"Mantener limpieza de quemadores")
	case "bomba de agua", "bomba hidráulica":
		recomendaciones = append(recomendaciones,
			"Verificar niveles de vibración mensualmente",
			"Inspeccionar sellos y empaques",
			"Revisar alineación del motor")
	case "transformador":
		recomendaciones = append(recomendaciones,
			"Verificar niveles de aceite dieléctrico",
			"Realizar termografía semestral",
			"Inspeccionar conexiones eléctricas")
	default:
		recomendaciones = append(recomendaciones,
			"Realizar mantenimiento preventivo según cronograma",
			"Documentar todas las intervenciones")
	}

	// Recomendaciones basadas en el estado
	if activo.Estado == "mantenimiento" {
		recomendaciones = append(recomendaciones,
			"Completar mantenimiento programado lo antes posible",
			"Verificar funcionamiento después de la intervención")
	}

	// Recomendaciones basadas en la última acción
	if ultimaAccion != nil && ultimaAccion.Estado == "pendiente" {
		recomendaciones = append(recomendaciones,
			fmt.Sprintf("Ejecutar acción pendiente: %s", ultimaAccion.Descripcion))
	}

	return recomendaciones
}

// Generar conclusiones automáticas
func (s *ReporteService) generarConclusiones(activo *models.Activo, observaciones string) string {
	baseConclusion := fmt.Sprintf("El activo %s se encuentra en estado %s.",
		activo.Nombre, activo.Estado)

	if observaciones != "" {
		baseConclusion += " Las observaciones del analista proporcionan detalles adicionales para el seguimiento."
	}

	if activo.Estado == "operativo" {
		baseConclusion += " Se recomienda continuar con el plan de mantenimiento preventivo."
	} else if activo.Estado == "mantenimiento" {
		baseConclusion += " Se requiere completar las tareas de mantenimiento programadas."
	}

	return baseConclusion
}

// Generar contenido básico del reporte
func (s *ReporteService) generarContenidoBasico(activo *models.Activo) string {
	return fmt.Sprintf("Reporte para %s - %s ubicado en %s. Estado actual: %s",
		activo.Tipo, activo.Nombre, activo.Ubicacion, activo.Estado)
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

	// Generar estructura del reporte para el PDF
	estructura := s.generarEstructuraInforme(activo, "")
	recomendaciones := s.generarRecomendaciones(activo, ultimaAccion)
	conclusiones := s.generarConclusiones(activo, "")

	// Preparar datos para el template con estructura completa
	data := ReporteData{
		Activo:          activo,
		Edificio:        edificio,
		UltimaAccion:    ultimaAccion,
		Acciones:        acciones,
		FechaGeneracion: time.Now().Format("02/01/2006 15:04:05"),
		AutorAnalista:   "Sistema Automático",
		Recomendaciones: recomendaciones,
		Conclusiones:    conclusiones,
		EstructuraCompleta: models.EstructuraInformeReporte{
			Resumen:            estructura["resumen"].(string),
			Observaciones:      estructura["observaciones"].(string),
			DatosActivo:        estructura["datos_activo"].(map[string]interface{}),
			AccionesRealizadas: estructura["acciones_realizadas"].([]string),
			Recomendaciones:    estructura["recomendaciones"].([]string),
			Conclusiones:       estructura["conclusiones"].(string),
		},
	}
	fmt.Printf("DATOS PREPARADOS PARA TEMPLATE CON ESTRUCTURA COMPLETA\n")

	// Renderizar HTML template
	htmlContent, err := s.renderHTMLTemplate(data)
	if err != nil {
		return nil, "", fmt.Errorf("error renderizando template HTML: %v", err)
	}

	fmt.Printf("HTML generado: %d caracteres\n", len(htmlContent))

	// Guardar HTML para debug
	if err := os.WriteFile("/tmp/debug_report.html", []byte(htmlContent), 0644); err == nil {
		fmt.Printf("HTML guardado en /tmp/debug_report.html para inspección\n")
	}

	// Generar PDF desde HTML
	pdfBytes, err := s.generatePDFFromHTML(htmlContent)
	if err != nil {
		return nil, "", fmt.Errorf("error generando PDF: %v", err)
	}

	// Crear registro del reporte con estructura completa
	reporte := &models.Reporte{
		ActivoID:              activoID,
		TipoReporte:           "PDF_ACTIVO",
		Contenido:             "Reporte PDF generado desde template HTML con estructura completa",
		ObservacionesAnalista: "",
		AutorAnalista:         "Sistema Automático",
		EstructuraInforme:     estructura,
		MetadataInforme: map[string]interface{}{
			"fecha_generacion": time.Now(),
			"tipo_generacion":  "automatica",
			"version_template": "2.0",
		},
		VersionReporte: 1,
		EstadoRevision: "pendiente",
		GeneradoEn:     time.Now(),
		Estado:         "generado",
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

// Obtener reportes con observaciones
func (s *ReporteService) ObtenerReportesConObservaciones(activoID int) ([]models.ReporteConObservaciones, error) {
	return s.repo.GetConObservaciones(activoID)
}

// Obtener todos los reportes con observaciones
func (s *ReporteService) ObtenerTodosLosReportes() ([]models.ReporteConObservaciones, error) {
	return s.repo.GetAll()
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
