package services

import (
	"documentacion/internal/models"
	"documentacion/internal/repository"
	"documentacion/storage"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ConsultaService struct {
	consultaRepo   *repository.ConsultaRepository
	documentoRepo  *repository.DocumentoRepository
	storageService storage.StorageService
	geminiAPIKey   string
	projectID      string
}

func NewConsultaService(
	consultaRepo *repository.ConsultaRepository,
	documentoRepo *repository.DocumentoRepository,
	storageService storage.StorageService,
	projectID string,
	geminiAPIKey string,
) *ConsultaService {
	return &ConsultaService{
		consultaRepo:   consultaRepo,
		documentoRepo:  documentoRepo,
		storageService: storageService,
		geminiAPIKey:   geminiAPIKey,
		projectID:      projectID,
	}
}

// ConsultarConIA procesa una pregunta específica sobre un documento usando IA simulada
func (s *ConsultaService) ConsultarConIA(documentoID int, pregunta string, documento *models.Documento) (*models.ResponseConsulta, error) {
	// Por ahora simular respuesta inteligente basada en el documento y pregunta
	respuesta := s.generarRespuestaSimulada(pregunta, documento)

	// Simular confianza basada en palabras clave
	confianza := s.calcularConfianzaSimulada(pregunta, documento)

	// Simular fuentes
	fuentes := s.extraerFuentesSimuladas(documento)

	response := &models.ResponseConsulta{
		Pregunta:    pregunta,
		Respuesta:   respuesta,
		DocumentoID: documentoID,
		Confianza:   &confianza,
		Fuentes:     fuentes,
		ProcesadoEn: time.Now(),
	}

	return response, nil
}

// GuardarConsulta guarda una consulta realizada en la base de datos
func (s *ConsultaService) GuardarConsulta(consulta *models.ConsultaIA) error {
	return s.consultaRepo.Create(consulta)
}

// BuscarConsultasSimilares busca consultas similares a una pregunta
func (s *ConsultaService) BuscarConsultasSimilares(documentoID int, pregunta string) ([]models.ConsultaSimilar, error) {
	return s.consultaRepo.BuscarSimilares(documentoID, pregunta)
}

// ObtenerHistorial obtiene el historial de consultas de un documento
func (s *ConsultaService) ObtenerHistorial(documentoID int, limit, offset int) (*models.HistorialConsultas, error) {
	consultas, total, err := s.consultaRepo.ObtenerPorDocumento(documentoID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.HistorialConsultas{
		DocumentoID: documentoID,
		Total:       total,
		Consultas:   consultas,
	}, nil
}

// ObtenerEstadisticas obtiene estadísticas globales de consultas
func (s *ConsultaService) ObtenerEstadisticas() (*models.EstadisticasConsultas, error) {
	return s.consultaRepo.ObtenerEstadisticas()
}

// Métodos de simulación (para desarrollo/testing)

func (s *ConsultaService) generarRespuestaSimulada(pregunta string, documento *models.Documento) string {
	preguntaLower := strings.ToLower(pregunta)

	// Respuestas basadas en palabras clave
	if strings.Contains(preguntaLower, "presión") {
		if documento.Categoria == "ficha_tecnica" {
			return fmt.Sprintf("Según la ficha técnica de %s, la presión máxima de operación es de 150 PSI (10.3 bar). La presión de trabajo recomendada está entre 120-140 PSI para óptimo rendimiento.", documento.Nombre)
		}
		return "La información sobre presión no está claramente especificada en este documento. Recomiendo consultar la ficha técnica oficial del equipo."
	}

	if strings.Contains(preguntaLower, "capacidad") || strings.Contains(preguntaLower, "potencia") {
		if strings.Contains(strings.ToLower(documento.Nombre), "caldera") {
			return fmt.Sprintf("La %s tiene una capacidad térmica de 500kW con una eficiencia energética del 92%%. Opera con gas natural o GLP según las especificaciones del fabricante.", documento.Nombre)
		}
		if strings.Contains(strings.ToLower(documento.Nombre), "bomba") {
			return fmt.Sprintf("La %s tiene una potencia de 50HP y está diseñada para aplicaciones industriales con alta demanda de caudal.", documento.Nombre)
		}
		return "La capacidad específica se detalla en las especificaciones técnicas del documento. Para información precisa, consulte la sección de características principales."
	}

	if strings.Contains(preguntaLower, "mantenimiento") || strings.Contains(preguntaLower, "mantención") {
		if documento.Categoria == "reporte_mantenimiento" {
			return fmt.Sprintf("Según el reporte de mantenimiento, la última intervención fue realizada el 30 de enero de 2024. Se ejecutaron tareas de limpieza, calibración y reemplazo de componentes. El estado general fue evaluado como satisfactorio con próxima inspección programada para abril 2024.")
		}
		return "Este documento contiene información de mantenimiento. Para detalles específicos sobre fechas y procedimientos, consulte la sección de historial de mantenimiento."
	}

	if strings.Contains(preguntaLower, "eficiencia") {
		return fmt.Sprintf("La eficiencia energética del equipo %s es del 92%%, cumpliendo con estándares industriales de alta eficiencia. Está certificado bajo normas ISO 9001 y CE.", documento.Nombre)
	}

	if strings.Contains(preguntaLower, "temperatura") {
		return "La temperatura máxima de operación es de 180°C, con rangos de trabajo recomendados entre 160-175°C para óptimo rendimiento y vida útil del equipo."
	}

	// Respuesta genérica
	return fmt.Sprintf("Basándome en el documento '%s' (categoría: %s), la información solicitada se encuentra en las especificaciones técnicas. %s Para detalles específicos, recomiendo revisar las secciones técnicas del documento.",
		documento.Nombre, documento.Categoria, documento.Descripcion)
}

func (s *ConsultaService) calcularConfianzaSimulada(pregunta string, documento *models.Documento) float64 {
	preguntaLower := strings.ToLower(pregunta)
	confianzaBase := 0.7

	// Incrementar confianza por coincidencias
	if documento.Categoria == "ficha_tecnica" && (strings.Contains(preguntaLower, "presión") || strings.Contains(preguntaLower, "capacidad")) {
		confianzaBase += 0.2
	}

	if documento.Categoria == "reporte_mantenimiento" && strings.Contains(preguntaLower, "mantenimiento") {
		confianzaBase += 0.25
	}

	if strings.Contains(strings.ToLower(documento.PalabrasClave), strings.Split(preguntaLower, " ")[0]) {
		confianzaBase += 0.1
	}

	// Agregar variación aleatoria pequeña
	variacion := (rand.Float64() - 0.5) * 0.1
	confianza := confianzaBase + variacion

	// Mantener entre 0.5 y 1.0
	if confianza > 1.0 {
		confianza = 1.0
	}
	if confianza < 0.5 {
		confianza = 0.5
	}

	return confianza
}

func (s *ConsultaService) extraerFuentesSimuladas(documento *models.Documento) []string {
	fuentes := []string{"Documento principal"}

	switch documento.Categoria {
	case "ficha_tecnica":
		fuentes = append(fuentes, "Tabla de especificaciones técnicas", "Sección de características principales")
	case "reporte_mantenimiento":
		fuentes = append(fuentes, "Reporte de mantenimiento", "Checklist de inspección")
	case "manual_fabricante":
		fuentes = append(fuentes, "Manual del fabricante", "Procedimientos de operación")
	default:
		fuentes = append(fuentes, "Contenido técnico del documento")
	}

	return fuentes
}
