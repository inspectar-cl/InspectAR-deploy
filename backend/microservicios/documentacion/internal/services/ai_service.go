package services

import (
	"context"
	"documentacion/config"
	"documentacion/internal/models"
	"documentacion/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	client         *genai.Client
	model          *genai.GenerativeModel
	analisisRepo   *repository.AnalisisRepository
	documentoRepo  *repository.DocumentoRepository
	storageService interface {
		GetFile(rutaArchivo string) (io.ReadCloser, error)
	}
}

func NewAIService(cfg *config.GeminiConfig, analisisRepo *repository.AnalisisRepository, documentoRepo *repository.DocumentoRepository, storageService interface {
	GetFile(rutaArchivo string) (io.ReadCloser, error)
}) (*AIService, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY no configurado")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("error inicializando cliente Gemini: %v", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.2)

	return &AIService{
		client:         client,
		model:          model,
		analisisRepo:   analisisRepo,
		documentoRepo:  documentoRepo,
		storageService: storageService,
	}, nil
}

// AnalizarDocumento realiza análisis de IA de un documento (HdU19)
func (s *AIService) AnalizarDocumento(documentoID int) (*models.AnalisisIA, error) {
	// Crear registro de análisis en estado "procesando"
	req := &models.CreateAnalisisRequest{DocumentoID: documentoID}
	analisis, err := s.analisisRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("error creando análisis: %v", err)
	}

	// Procesar en background
	go s.procesarAnalisisAsync(analisis.ID, documentoID)

	return analisis, nil
}

// procesarAnalisisAsync procesa el análisis de manera asíncrona
func (s *AIService) procesarAnalisisAsync(analisisID, documentoID int) {
	// Obtener información del documento
	documento, err := s.documentoRepo.GetByID(documentoID)
	if err != nil {
		log.Printf("Error obteniendo documento %d: %v", documentoID, err)
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
		return
	}

	// Solo procesar PDFs por ahora
	if documento.TipoArchivo != "pdf" {
		log.Printf("Tipo de archivo no soportado para análisis: %s", documento.TipoArchivo)
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
		return
	}

	// Obtener archivo del storage
	file, err := s.storageService.GetFile(documento.RutaArchivo)
	if err != nil {
		log.Printf("Error obteniendo archivo %s: %v", documento.RutaArchivo, err)
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
		return
	}
	defer file.Close()

	// Leer contenido del archivo
	contenido, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error leyendo archivo: %v", err)
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
		return
	}

	// Crear prompt para Gemini
	prompt := s.crearPromptAnalisis(documento)

	// Analizar con Gemini
	ctx := context.Background()

	// Para archivos PDF, usamos el contenido como texto
	// En una implementación completa, necesitarías un extractor de texto de PDF
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt+"\n\nContenido del documento (extracto): "+string(contenido[:min(2000, len(contenido))])))
	if err != nil {
		log.Printf("Error en análisis Gemini: %v", err)
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
		return
	}

	// Procesar respuesta
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		resultado := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

		// Extraer información estructurada
		resumen, puntosClaves, graficos := s.procesarResultadoIA(resultado)

		// Actualizar análisis con resultados
		err = s.analisisRepo.Update(analisisID, resumen, puntosClaves, graficos, models.EstadoCompletado)
		if err != nil {
			log.Printf("Error actualizando análisis: %v", err)
		} else {
			log.Printf("Análisis completado para documento %d", documentoID)
		}
	} else {
		log.Printf("No se recibió respuesta válida de Gemini")
		s.analisisRepo.Update(analisisID, "", "", "", models.EstadoError)
	}
}

// crearPromptAnalisis crea el prompt específico según el tipo de documento
func (s *AIService) crearPromptAnalisis(documento *models.Documento) string {
	basePrompt := fmt.Sprintf(`
Analiza este documento técnico de mantenimiento industrial. El documento es de categoría "%s" y se relaciona con el activo ID %d.

Por favor proporciona:
1. RESUMEN: Un resumen ejecutivo del documento (máximo 200 palabras)
2. PUNTOS_CLAVES: Lista de 5-10 puntos importantes en formato JSON array
3. GRAFICOS: Descripción de gráficos, tablas o diagramas encontrados en formato JSON array

Formato de respuesta esperado:
RESUMEN: [tu resumen aquí]
PUNTOS_CLAVES: ["punto1", "punto2", "punto3", ...]
GRAFICOS: ["descripción gráfico 1", "descripción gráfico 2", ...]

Enfócate especialmente en:
- Información técnica relevante para mantenimiento
- Especificaciones críticas
- Procedimientos de seguridad
- Recomendaciones del fabricante
- Datos de rendimiento o diagnóstico
`, documento.Categoria, documento.ActivoID)

	// Personalizar según categoría
	switch documento.Categoria {
	case models.CategoriaFichaTecnica:
		basePrompt += "\nEste es una ficha técnica. Enfócate en especificaciones, capacidades y características técnicas."
	case models.CategoriaInformeMantenimiento:
		basePrompt += "\nEste es un informe de mantenimiento. Enfócate en acciones realizadas, hallazgos y recomendaciones."
	case models.CategoriaDiagnostico:
		basePrompt += "\nEste es un diagnóstico. Enfócate en problemas identificados, causas y soluciones propuestas."
	case models.CategoriaManualFabricante:
		basePrompt += "\nEste es un manual del fabricante. Enfócate en procedimientos, especificaciones y guías de operación."
	}

	return basePrompt
}

// procesarResultadoIA extrae la información estructurada de la respuesta de Gemini
func (s *AIService) procesarResultadoIA(resultado string) (resumen, puntosClaves, graficos string) {
	// Parser simple para extraer las secciones
	// En una implementación más robusta, usarías regex o parsing más sofisticado

	lines := strings.Split(resultado, "\n")
	var currentSection string
	var resumenLines, puntosLines, graficosLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RESUMEN:") {
			currentSection = "resumen"
			resumenLines = append(resumenLines, strings.TrimPrefix(line, "RESUMEN:"))
		} else if strings.HasPrefix(line, "PUNTOS_CLAVES:") {
			currentSection = "puntos"
			puntosLines = append(puntosLines, strings.TrimPrefix(line, "PUNTOS_CLAVES:"))
		} else if strings.HasPrefix(line, "GRAFICOS:") {
			currentSection = "graficos"
			graficosLines = append(graficosLines, strings.TrimPrefix(line, "GRAFICOS:"))
		} else if line != "" {
			switch currentSection {
			case "resumen":
				resumenLines = append(resumenLines, line)
			case "puntos":
				puntosLines = append(puntosLines, line)
			case "graficos":
				graficosLines = append(graficosLines, line)
			}
		}
	}

	resumen = strings.TrimSpace(strings.Join(resumenLines, " "))

	// Intentar parsear JSON para puntos claves y gráficos
	puntosStr := strings.TrimSpace(strings.Join(puntosLines, " "))
	if strings.HasPrefix(puntosStr, "[") {
		puntosClaves = puntosStr
	} else {
		// Fallback: crear array JSON simple
		defaultPuntos := []string{"Análisis pendiente de procesamiento manual"}
		puntosJSON, _ := json.Marshal(defaultPuntos)
		puntosClaves = string(puntosJSON)
	}

	graficosStr := strings.TrimSpace(strings.Join(graficosLines, " "))
	if strings.HasPrefix(graficosStr, "[") {
		graficos = graficosStr
	} else {
		// Fallback: crear array JSON simple
		defaultGraficos := []string{"No se detectaron gráficos específicos"}
		graficosJSON, _ := json.Marshal(defaultGraficos)
		graficos = string(graficosJSON)
	}

	return resumen, puntosClaves, graficos
}

// ObtenerAnalisis obtiene el análisis de un documento
func (s *AIService) ObtenerAnalisis(documentoID int) (*models.AnalisisIA, error) {
	return s.analisisRepo.GetByDocumentoID(documentoID)
}

// ObtenerAnalisisPorActivo obtiene todos los análisis de documentos de un activo
func (s *AIService) ObtenerAnalisisPorActivo(activoID int) ([]models.AnalisisIA, error) {
	return s.analisisRepo.GetByActivoID(activoID)
}

// Close cierra el cliente de Gemini
func (s *AIService) Close() error {
	return s.client.Close()
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
