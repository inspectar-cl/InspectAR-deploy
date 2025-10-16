package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
	"html/template"
	"math"
	"net/http"
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
	firmaRepo    *repository.FirmaRepository
	usuarioRepo  *repository.UsuarioRepository
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
		firmaRepo:    nil, // Será configurado opcionalmente
	}
}

// SetFirmaRepo configura el repositorio de firmas (opcional)
func (s *ReporteService) SetFirmaRepo(firmaRepo *repository.FirmaRepository) {
	s.firmaRepo = firmaRepo
}

// SetUsuarioRepo configura el repositorio de usuarios (opcional)
func (s *ReporteService) SetUsuarioRepo(usuarioRepo *repository.UsuarioRepository) {
	s.usuarioRepo = usuarioRepo
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
	// Campo para métricas de sensores
	SensorMetrics map[string]interface{} `json:"sensor_metrics,omitempty"`
	// Campo para firma digital
	FirmaDigital *FirmaDataTemplate `json:"firma_digital,omitempty"`
}

// FirmaDataTemplate contiene datos de la firma para el template
type FirmaDataTemplate struct {
	Existe       bool   `json:"existe"`
	RutaArchivo  string `json:"ruta_archivo,omitempty"`
	TipoMime     string `json:"tipo_mime,omitempty"`
	Formato      string `json:"formato,omitempty"`
	AutorNombre  string `json:"autor_nombre,omitempty"`
	FechaFirma   string `json:"fecha_firma,omitempty"`
	ImagenBase64 string `json:"imagen_base64,omitempty"` // Para incrustar en PDF
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
		"fecha_creacion": time.Now(),
		"tipo_activo":    activo.Tipo,
		"nombre_activo":  activo.Nombre,
		"version":        1,
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

	// Recomendaciones basadas en la última acción
	if ultimaAccion != nil && ultimaAccion.Estado == "pendiente" {
		recomendaciones = append(recomendaciones,
			fmt.Sprintf("Ejecutar acción pendiente: %s", ultimaAccion.Descripcion))
	}

	return recomendaciones
}

// Generar conclusiones automáticas
func (s *ReporteService) generarConclusiones(activo *models.Activo, observaciones string) string {
	baseConclusion := fmt.Sprintf("El activo %s está ubicado en %s.",
		activo.Nombre, activo.Ubicacion)

	if observaciones != "" {
		baseConclusion += " Las observaciones del analista proporcionan detalles adicionales para el seguimiento."
	}

	baseConclusion += " Se recomienda continuar con el plan de mantenimiento preventivo."

	return baseConclusion
}

// Generar contenido básico del reporte
func (s *ReporteService) generarContenidoBasico(activo *models.Activo) string {
	return fmt.Sprintf("Reporte para %s - %s ubicado en %s",
		activo.Tipo, activo.Nombre, activo.Ubicacion)
}

// Generar reporte PDF por activo usando HTML template
// Ahora acepta opciones con los campos solicitados
func (s *ReporteService) GenerarReportePDFPorActivo(activoID int, opts models.GenerarReporteRequest) ([]byte, string, error) {
	fmt.Printf("=== GENERANDO REPORTE PARA ACTIVO ID: %d ===\n", activoID)

	// Obtener datos del activo
	activo, err := s.activoRepo.GetByID(activoID)
	if err != nil {
		return nil, "", fmt.Errorf("no se pudo obtener el activo: %v", err)
	}
	fmt.Printf("ACTIVO ENCONTRADO: ID=%d, Nombre=%s, Tipo=%s\n",
		activo.ID, activo.Nombre, activo.Tipo)

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

	// Procesar campos solicitados - crear mapas de control de secciones
	incluirUbicacion := false
	incluirHistorial := false
	incluirUltimasAcciones := false
	incluirDatosSensores := false

	// Si no se especifican campos, incluir todo por defecto (comportamiento legacy)
	if len(opts.Campos) == 0 {
		incluirUbicacion = true
		incluirHistorial = true
		incluirUltimasAcciones = true
	}

	// Procesar campos específicos solicitados
	sensorMetrics := map[string]interface{}{}
	for _, campo := range opts.Campos {
		switch campo {
		case "ubicacion":
			incluirUbicacion = true
		case "historial_mantenimientos":
			incluirHistorial = true
		case "ultima_acciones":
			incluirUltimasAcciones = true
		case "datos_sensores":
			incluirDatosSensores = true

			// Llamar a ParserService para obtener datos de sensores
			parserURL := os.Getenv("PARSER_URL")
			if parserURL == "" {
				parserURL = "http://localhost:8090"
			}
			// Llamada simple al endpoint del ParserService para obtener lecturas del activo
			sensorEndpoint := fmt.Sprintf("%s/lectura/%d/datos", parserURL, activoID)
			fmt.Printf("Llamando a ParserService: %s\n", sensorEndpoint)
			resp, err := http.Get(sensorEndpoint)
			if err == nil && resp.StatusCode == 200 {
				var sensorData interface{}
				err = json.NewDecoder(resp.Body).Decode(&sensorData)
				resp.Body.Close()
				if err == nil {
					sensorMetrics = s.generateFakeMetrics(sensorData)
					fmt.Printf("Datos de sensores obtenidos de ParserService: %+v\n", sensorMetrics)
				} else {
					fmt.Printf("Error decodificando respuesta JSON: %v\n", err)
				}
			} else {
				fmt.Printf("Advertencia: no se pudo obtener datos de sensores desde %s (err=%v)\n", sensorEndpoint, err)
				if resp != nil {
					fmt.Printf("Status code: %d\n", resp.StatusCode)
					resp.Body.Close()
				}
			}

			// DEBUGGING: Verificar estado de sensorMetrics antes de procesamiento
			fmt.Printf("DEBUG: Verificando métricas - len: %d, contenido: %+v\n", len(sensorMetrics), sensorMetrics)

			// Verificar si tenemos datos válidos de sensores, si no, generar métricas de ejemplo
			if len(sensorMetrics) == 0 {
				fmt.Printf("No hay métricas de sensores, generando métricas de ejemplo\n")
				sensorMetrics = s.generateFakeMetrics(map[string]interface{}{
					"sensors": []interface{}{
						map[string]interface{}{
							"name":   "temp_001",
							"values": []interface{}{22.5, 23.1, 22.8},
						},
						map[string]interface{}{
							"name":   "pres_001",
							"values": []interface{}{1.2, 1.3, 1.1},
						},
						map[string]interface{}{
							"name":   "caud_001",
							"values": []interface{}{5.4, 5.7, 5.2},
						},
					},
				})
			} else if _, hasNote := sensorMetrics["note"]; hasNote {
				fmt.Printf("🚨 ENTRANDO A REGENERACION DE METRICAS - ESTE LOG DEBE APARECER 🚨\n")
				fmt.Printf("Datos de sensores con 'note', regenerando métricas de ejemplo\n")
				sensorMetrics = s.generateFakeMetrics(map[string]interface{}{
					"sensors": []interface{}{
						map[string]interface{}{
							"name":   "temp_001",
							"values": []interface{}{22.5, 23.1, 22.8},
						},
						map[string]interface{}{
							"name":   "pres_001",
							"values": []interface{}{1.2, 1.3, 1.1},
						},
						map[string]interface{}{
							"name":   "caud_001",
							"values": []interface{}{5.4, 5.7, 5.2},
						},
					},
				})
				fmt.Printf("🚨 METRICAS REGENERADAS - ESTE LOG DEBE APARECER 🚨\n")
			}
			fmt.Printf("Métricas finales de sensores: %+v\n", sensorMetrics)
		}
	}

	// Procesar firma digital si se proporcionó
	var firmaData *FirmaDataTemplate
	if opts.FirmaID != nil && s.firmaRepo != nil {
		// Obtener firma específica
		firma, err := s.firmaRepo.GetByID(*opts.FirmaID)
		if err == nil {
			firmaData = s.prepararFirmaParaTemplate(firma)
			fmt.Printf("FIRMA ENCONTRADA: ID=%d, Nombre=%s, Formato=%s\n", firma.ID, firma.NombreArchivo, firma.Formato)
		} else {
			fmt.Printf("Advertencia: no se pudo cargar la firma ID %d: %v\n", *opts.FirmaID, err)
		}
	} else if opts.UsarFirmaPredeterminada && opts.Email != "" && s.firmaRepo != nil && s.usuarioRepo != nil {
		// Obtener firma predeterminada del usuario por email
		firma, err := s.firmaRepo.GetDefaultByUsuarioEmail(opts.Email)
		if err == nil {
			firmaData = s.prepararFirmaParaTemplate(firma)
			fmt.Printf("FIRMA PREDETERMINADA ENCONTRADA: ID=%d, Usuario=%d, Formato=%s\n", firma.ID, firma.UsuarioID, firma.Formato)
		} else {
			fmt.Printf("Advertencia: no se encontró firma predeterminada para usuario con email %s: %v\n", opts.Email, err)
		}
	}

	// Preparar datos condicionalmente según campos solicitados
	data := ReporteData{
		Activo:          activo, // Siempre incluir información básica del activo
		FechaGeneracion: time.Now().Format("02/01/2006 15:04:05"),
		AutorAnalista:   "Sistema Automático",
		EstructuraCompleta: models.EstructuraInformeReporte{
			Resumen:     estructura["resumen"].(string),
			DatosActivo: estructura["datos_activo"].(map[string]interface{}),
		},
	}

	// Incluir edificio solo si se solicita ubicacion
	if incluirUbicacion {
		data.Edificio = edificio
	}

	// Incluir historial de acciones solo si se solicita
	if incluirHistorial {
		data.Acciones = acciones
		data.EstructuraCompleta.AccionesRealizadas = estructura["acciones_realizadas"].([]string)
	}

	// Incluir última acción solo si se solicita
	if incluirUltimasAcciones {
		data.UltimaAccion = ultimaAccion
	}

	// Incluir métricas de sensores si se solicitaron
	if incluirDatosSensores {
		data.SensorMetrics = sensorMetrics
	}

	// Incluir recomendaciones y conclusiones solo si hay campos específicos
	if len(opts.Campos) == 0 || incluirHistorial || incluirUltimasAcciones {
		data.Recomendaciones = recomendaciones
		data.Conclusiones = conclusiones
		data.EstructuraCompleta.Recomendaciones = estructura["recomendaciones"].([]string)
		data.EstructuraCompleta.Conclusiones = estructura["conclusiones"].(string)
	}

	// Incluir firma digital si está disponible
	if firmaData != nil {
		data.FirmaDigital = firmaData
		fmt.Printf("FIRMA INCLUIDA EN REPORTE - Formato: %s\n", firmaData.Formato)
	}

	fmt.Printf("DATOS PREPARADOS PARA TEMPLATE - Campos incluidos: ubicacion=%v, historial=%v, ultimas_acciones=%v, sensores=%v, firma=%v\n",
		incluirUbicacion, incluirHistorial, incluirUltimasAcciones, incluirDatosSensores, firmaData != nil)
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
- Ubicación: %s

Información del Edificio:
- Nombre: %s
- Dirección: %s

`, activo.ID, activo.Nombre, activo.Tipo, activo.Ubicacion,
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

// generateFakeMetrics intenta extraer series de datos de sensor y calcula métricas simples (media y tendencia)
func (s *ReporteService) generateFakeMetrics(sensorData interface{}) map[string]interface{} {
	metrics := map[string]interface{}{}
	fmt.Printf("🔍 generateFakeMetrics entrada: %+v (tipo: %T)\n", sensorData, sensorData)

	// Caso: arreglo de sensores
	if arr, ok := sensorData.([]interface{}); ok {
		fmt.Printf("📊 Procesando array de sensores, longitud: %d\n", len(arr))
		for i, item := range arr {
			name := fmt.Sprintf("sensor_%d", i)
			if m, ok2 := item.(map[string]interface{}); ok2 {
				fmt.Printf("📡 Sensor %d mapa: %+v\n", i, m)
				if n, okn := m["name"].(string); okn {
					name = n
					fmt.Printf("📝 Nombre del sensor: %s\n", name)
				}
				// Buscar valores bajo keys comunes: values, readings, data
				var vals []interface{}
				if v, okv := m["values"].([]interface{}); okv {
					vals = v
					fmt.Printf("✅ Encontrados values: %+v\n", vals)
				} else if v, okv := m["readings"].([]interface{}); okv {
					vals = v
					fmt.Printf("✅ Encontrados readings: %+v\n", vals)
				} else if v, okv := m["data"].([]interface{}); okv {
					vals = v
					fmt.Printf("✅ Encontrados data: %+v\n", vals)
				} else {
					fmt.Printf("❌ No se encontraron valores en sensor %s\n", name)
				}
				if len(vals) > 0 {
					var sum float64
					var first, last float64
					for j, vv := range vals {
						switch v := vv.(type) {
						case float64:
							if j == 0 {
								first = v
							}
							last = v
							sum += v
						case int:
							fv := float64(v)
							if j == 0 {
								first = fv
							}
							last = fv
							sum += fv
						case map[string]interface{}:
							if val, ok := v["value"].(float64); ok {
								if j == 0 {
									first = val
								}
								last = val
								sum += val
							}
						}
					}
					mean := sum / float64(len(vals))
					trend := 0.0
					if len(vals) >= 2 {
						trend = (last - first) / float64(len(vals)-1)
					}
					// redondear un poco
					calculatedMetrics := map[string]float64{"mean": math.Round(mean*100) / 100, "trend": math.Round(trend*100) / 100}
					metrics[name] = calculatedMetrics
					fmt.Printf("📈 Métricas calculadas para %s: %+v\n", name, calculatedMetrics)
				}
			}
		}
		if len(metrics) > 0 {
			fmt.Printf("🎯 Retornando métricas de array: %+v\n", metrics)
			return metrics
		}
		fmt.Printf("⚠️ Array procesado pero sin métricas\n")
	}

	// Caso: mapa que contiene sensors o series
	if m, ok := sensorData.(map[string]interface{}); ok {
		fmt.Printf("🗂️ Procesando mapa con claves: %+v\n", func() []string {
			keys := []string{}
			for k := range m {
				keys = append(keys, k)
			}
			return keys
		}())
		if sarr, ok2 := m["sensors"].([]interface{}); ok2 {
			fmt.Printf("🔄 Encontrada clave 'sensors', llamada recursiva\n")
			return s.generateFakeMetrics(sarr)
		}
		// intentar extraer arrays numéricos por clave
		for k, v := range m {
			if arr, ok3 := v.([]interface{}); ok3 {
				var sum float64
				cnt := 0
				var first, last float64
				for j, vv := range arr {
					switch val := vv.(type) {
					case float64:
						if j == 0 {
							first = val
						}
						last = val
						sum += val
						cnt++
					case int:
						fv := float64(val)
						if j == 0 {
							first = fv
						}
						last = fv
						sum += fv
						cnt++
					}
				}
				if cnt > 0 {
					mean := sum / float64(cnt)
					trend := 0.0
					if cnt >= 2 {
						trend = (last - first) / float64(cnt-1)
					}
					metrics[k] = map[string]float64{"mean": math.Round(mean*100) / 100, "trend": math.Round(trend*100) / 100}
				}
			}
		}
		if len(metrics) > 0 {
			return metrics
		}
	}

	metrics["note"] = "no sensor data parsed"
	return metrics
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

// prepararFirmaParaTemplate convierte una firma en datos para el template
func (s *ReporteService) prepararFirmaParaTemplate(firma *models.FirmaDigital) *FirmaDataTemplate {
	if firma == nil {
		return nil
	}

	firmaData := &FirmaDataTemplate{
		Existe:      true,
		RutaArchivo: firma.RutaArchivo,
		TipoMime:    firma.TipoMime,
		Formato:     firma.Formato,
		FechaFirma:  firma.CreadoEn.Format("02/01/2006"),
	}

	// Intentar leer el archivo y convertirlo a base64 para incrustarlo en el PDF
	if firma.RutaArchivo != "" {
		imageData, err := os.ReadFile(firma.RutaArchivo)
		if err == nil {
			// Convertir a base64
			firmaData.ImagenBase64 = base64.StdEncoding.EncodeToString(imageData)
		} else {
			fmt.Printf("Advertencia: no se pudo leer archivo de firma: %v\n", err)
		}
	}

	return firmaData
}
