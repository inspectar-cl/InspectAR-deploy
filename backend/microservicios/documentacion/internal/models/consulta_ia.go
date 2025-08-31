package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ConsultaIA representa una consulta interactiva realizada sobre un documento
type ConsultaIA struct {
	ID                int             `json:"id" db:"id"`
	DocumentoID       int             `json:"documento_id" db:"documento_id"`
	Pregunta          string          `json:"pregunta" db:"pregunta" binding:"required"`
	Respuesta         string          `json:"respuesta" db:"respuesta"`
	Confianza         *float64        `json:"confianza" db:"confianza"` // 0.0 - 1.0
	Fuentes           json.RawMessage `json:"fuentes" db:"fuentes"`
	TiempoRespuestaMs *int            `json:"tiempo_respuesta_ms" db:"tiempo_respuesta_ms"`
	CreadoEn          time.Time       `json:"creado_en" db:"creado_en"`
}

// RequestConsulta representa el payload para hacer una consulta
type RequestConsulta struct {
	Pregunta string `json:"pregunta" binding:"required,min=5,max=500"`
}

// ResponseConsulta representa la respuesta a una consulta
type ResponseConsulta struct {
	Pregunta          string    `json:"pregunta"`
	Respuesta         string    `json:"respuesta"`
	DocumentoID       int       `json:"documento_id"`
	Confianza         *float64  `json:"confianza,omitempty"`
	Fuentes           []string  `json:"fuentes,omitempty"`
	TiempoRespuestaMs *int      `json:"tiempo_respuesta_ms,omitempty"`
	ProcesadoEn       time.Time `json:"procesado_en"`
}

// ConsultaSimilar representa una consulta similar encontrada
type ConsultaSimilar struct {
	ID          int     `json:"id"`
	Pregunta    string  `json:"pregunta"`
	Respuesta   string  `json:"respuesta"`
	Confianza   float64 `json:"confianza"`
	Similaridad float64 `json:"similaridad"`
}

// HistorialConsultas representa el histórico de consultas de un documento
type HistorialConsultas struct {
	DocumentoID int          `json:"documento_id"`
	Total       int          `json:"total"`
	Consultas   []ConsultaIA `json:"consultas"`
}

// EstadisticasConsultas representa estadísticas de uso de consultas
type EstadisticasConsultas struct {
	TotalConsultas           int                   `json:"total_consultas"`
	ConfianzaPromedio        float64               `json:"confianza_promedio"`
	TiempoPromedioMs         int                   `json:"tiempo_promedio_ms"`
	ConsultasHoy             int                   `json:"consultas_hoy"`
	DocumentosMasConsultados []DocumentoConsultado `json:"documentos_mas_consultados"`
}

// DocumentoConsultado representa un documento y su frecuencia de consultas
type DocumentoConsultado struct {
	DocumentoID    int    `json:"documento_id"`
	NombreDoc      string `json:"nombre_documento"`
	TotalConsultas int    `json:"total_consultas"`
}

// Métodos helper para trabajar con JSON

// GetFuentesAsSlice convierte las fuentes JSON a un slice de strings
func (c *ConsultaIA) GetFuentesAsSlice() ([]string, error) {
	if c.Fuentes == nil {
		return []string{}, nil
	}

	var fuentes []string
	err := json.Unmarshal(c.Fuentes, &fuentes)
	return fuentes, err
}

// SetFuentesFromSlice convierte un slice de strings a JSON para almacenar
func (c *ConsultaIA) SetFuentesFromSlice(fuentes []string) error {
	if fuentes == nil {
		c.Fuentes = nil
		return nil
	}

	jsonData, err := json.Marshal(fuentes)
	if err != nil {
		return err
	}
	c.Fuentes = jsonData
	return nil
}

// ToResponse convierte ConsultaIA a ResponseConsulta
func (c *ConsultaIA) ToResponse() ResponseConsulta {
	fuentes, _ := c.GetFuentesAsSlice()

	return ResponseConsulta{
		Pregunta:          c.Pregunta,
		Respuesta:         c.Respuesta,
		DocumentoID:       c.DocumentoID,
		Confianza:         c.Confianza,
		Fuentes:           fuentes,
		TiempoRespuestaMs: c.TiempoRespuestaMs,
		ProcesadoEn:       c.CreadoEn,
	}
}

// Validaciones

// ValidarPregunta valida que la pregunta sea apropiada
func (r *RequestConsulta) ValidarPregunta() error {
	if len(r.Pregunta) < 5 {
		return errors.New("la pregunta debe tener al menos 5 caracteres")
	}
	if len(r.Pregunta) > 500 {
		return errors.New("la pregunta no puede exceder 500 caracteres")
	}
	return nil
}

// EsPreguntaTecnica verifica si la pregunta es de naturaleza técnica
func (r *RequestConsulta) EsPreguntaTecnica() bool {
	palabrasTecnicas := []string{
		"presión", "temperatura", "capacidad", "eficiencia", "potencia",
		"mantenimiento", "mantención", "calibración", "especificación",
		"parámetro", "medición", "valor", "rango", "límite", "tolerancia",
		"cuando", "cuanto", "cuál", "cómo", "dónde", "por qué",
	}

	preguntaLower := strings.ToLower(r.Pregunta)
	for _, palabra := range palabrasTecnicas {
		if strings.Contains(preguntaLower, palabra) {
			return true
		}
	}
	return false
}
