package repository

import (
	"database/sql"
	"documentacion/internal/models"
	"fmt"
)

type ConsultaRepository struct {
	db *sql.DB
}

func NewConsultaRepository(db *sql.DB) *ConsultaRepository {
	return &ConsultaRepository{db: db}
}

// Create guarda una nueva consulta en la base de datos
func (r *ConsultaRepository) Create(consulta *models.ConsultaIA) error {
	query := `
		INSERT INTO consultas_ia (documento_id, pregunta, respuesta, confianza, fuentes, tiempo_respuesta_ms, creado_en)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		consulta.DocumentoID,
		consulta.Pregunta,
		consulta.Respuesta,
		consulta.Confianza,
		consulta.Fuentes,
		consulta.TiempoRespuestaMs,
		consulta.CreadoEn,
	).Scan(&consulta.ID)

	return err
}

// BuscarSimilares busca consultas similares usando la función de PostgreSQL
func (r *ConsultaRepository) BuscarSimilares(documentoID int, pregunta string) ([]models.ConsultaSimilar, error) {
	query := `
		SELECT id, pregunta, respuesta, 
		       COALESCE(confianza, 0) as confianza,
		       similarity(pregunta, $2) as similaridad
		FROM consultas_ia 
		WHERE documento_id = $1 
		  AND similarity(pregunta, $2) > 0.3
		ORDER BY similaridad DESC 
		LIMIT 5`

	rows, err := r.db.Query(query, documentoID, pregunta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var similares []models.ConsultaSimilar
	for rows.Next() {
		var similar models.ConsultaSimilar
		err := rows.Scan(
			&similar.ID,
			&similar.Pregunta,
			&similar.Respuesta,
			&similar.Confianza,
			&similar.Similaridad,
		)
		if err != nil {
			return nil, err
		}
		similares = append(similares, similar)
	}

	return similares, rows.Err()
}

// ObtenerPorDocumento obtiene consultas de un documento con paginación
func (r *ConsultaRepository) ObtenerPorDocumento(documentoID int, limit, offset int) ([]models.ConsultaIA, int, error) {
	// Primero obtener el total
	var total int
	countQuery := `SELECT COUNT(*) FROM consultas_ia WHERE documento_id = $1`
	err := r.db.QueryRow(countQuery, documentoID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Luego obtener los registros paginados
	query := `
		SELECT id, documento_id, pregunta, respuesta, confianza, fuentes, 
		       tiempo_respuesta_ms, creado_en
		FROM consultas_ia 
		WHERE documento_id = $1 
		ORDER BY creado_en DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, documentoID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var consultas []models.ConsultaIA
	for rows.Next() {
		var consulta models.ConsultaIA
		err := rows.Scan(
			&consulta.ID,
			&consulta.DocumentoID,
			&consulta.Pregunta,
			&consulta.Respuesta,
			&consulta.Confianza,
			&consulta.Fuentes,
			&consulta.TiempoRespuestaMs,
			&consulta.CreadoEn,
		)
		if err != nil {
			return nil, 0, err
		}
		consultas = append(consultas, consulta)
	}

	return consultas, total, rows.Err()
}

// ObtenerEstadisticas calcula estadísticas globales de consultas
func (r *ConsultaRepository) ObtenerEstadisticas() (*models.EstadisticasConsultas, error) {
	var stats models.EstadisticasConsultas

	// Consulta principal para estadísticas básicas
	query := `
		SELECT 
			COUNT(*) as total_consultas,
			COALESCE(AVG(confianza), 0) as confianza_promedio,
			COALESCE(AVG(tiempo_respuesta_ms), 0) as tiempo_promedio_ms,
			COUNT(CASE WHEN DATE(creado_en) = CURRENT_DATE THEN 1 END) as consultas_hoy
		FROM consultas_ia 
		WHERE confianza IS NOT NULL`

	err := r.db.QueryRow(query).Scan(
		&stats.TotalConsultas,
		&stats.ConfianzaPromedio,
		&stats.TiempoPromedioMs,
		&stats.ConsultasHoy,
	)
	if err != nil {
		return nil, err
	}

	// Obtener documentos más consultados
	docsQuery := `
		SELECT 
			c.documento_id,
			d.nombre as nombre_documento,
			COUNT(*) as total_consultas
		FROM consultas_ia c
		JOIN documentos d ON c.documento_id = d.id
		GROUP BY c.documento_id, d.nombre
		ORDER BY total_consultas DESC
		LIMIT 5`

	rows, err := r.db.Query(docsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docsConsultados []models.DocumentoConsultado
	for rows.Next() {
		var doc models.DocumentoConsultado
		err := rows.Scan(
			&doc.DocumentoID,
			&doc.NombreDoc,
			&doc.TotalConsultas,
		)
		if err != nil {
			return nil, err
		}
		docsConsultados = append(docsConsultados, doc)
	}

	stats.DocumentosMasConsultados = docsConsultados
	return &stats, rows.Err()
}

// ObtenerPorID obtiene una consulta específica por ID
func (r *ConsultaRepository) ObtenerPorID(id int) (*models.ConsultaIA, error) {
	query := `
		SELECT id, documento_id, pregunta, respuesta, confianza, fuentes, 
		       tiempo_respuesta_ms, creado_en
		FROM consultas_ia 
		WHERE id = $1`

	var consulta models.ConsultaIA
	err := r.db.QueryRow(query, id).Scan(
		&consulta.ID,
		&consulta.DocumentoID,
		&consulta.Pregunta,
		&consulta.Respuesta,
		&consulta.Confianza,
		&consulta.Fuentes,
		&consulta.TiempoRespuestaMs,
		&consulta.CreadoEn,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("consulta no encontrada")
	}
	if err != nil {
		return nil, err
	}

	return &consulta, nil
}

// EliminarPorDocumento elimina todas las consultas de un documento
func (r *ConsultaRepository) EliminarPorDocumento(documentoID int) error {
	query := `DELETE FROM consultas_ia WHERE documento_id = $1`
	_, err := r.db.Exec(query, documentoID)
	return err
}

// ObtenerConsultasRecientes obtiene las consultas más recientes del sistema
func (r *ConsultaRepository) ObtenerConsultasRecientes(limit int) ([]models.ConsultaIA, error) {
	query := `
		SELECT c.id, c.documento_id, c.pregunta, c.respuesta, c.confianza, 
		       c.fuentes, c.tiempo_respuesta_ms, c.creado_en
		FROM consultas_ia c
		ORDER BY c.creado_en DESC 
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultas []models.ConsultaIA
	for rows.Next() {
		var consulta models.ConsultaIA
		err := rows.Scan(
			&consulta.ID,
			&consulta.DocumentoID,
			&consulta.Pregunta,
			&consulta.Respuesta,
			&consulta.Confianza,
			&consulta.Fuentes,
			&consulta.TiempoRespuestaMs,
			&consulta.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		consultas = append(consultas, consulta)
	}

	return consultas, rows.Err()
}

// BuscarPorTexto busca consultas que contengan cierto texto
func (r *ConsultaRepository) BuscarPorTexto(texto string, limit int) ([]models.ConsultaIA, error) {
	query := `
		SELECT id, documento_id, pregunta, respuesta, confianza, fuentes, 
		       tiempo_respuesta_ms, creado_en
		FROM consultas_ia 
		WHERE pregunta ILIKE $1 OR respuesta ILIKE $1
		ORDER BY creado_en DESC 
		LIMIT $2`

	searchTerm := "%" + texto + "%"
	rows, err := r.db.Query(query, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultas []models.ConsultaIA
	for rows.Next() {
		var consulta models.ConsultaIA
		err := rows.Scan(
			&consulta.ID,
			&consulta.DocumentoID,
			&consulta.Pregunta,
			&consulta.Respuesta,
			&consulta.Confianza,
			&consulta.Fuentes,
			&consulta.TiempoRespuestaMs,
			&consulta.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		consultas = append(consultas, consulta)
	}

	return consultas, rows.Err()
}

// GetEstadisticas obtiene estadísticas generales de consultas
func (r *ConsultaRepository) GetEstadisticas() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_consultas,
			COUNT(DISTINCT documento_id) as documentos_consultados,
			AVG(confianza) as confianza_promedio,
			AVG(tiempo_respuesta_ms) as tiempo_promedio_ms,
			COUNT(CASE WHEN confianza >= 0.8 THEN 1 END) as consultas_alta_confianza,
			DATE_TRUNC('day', MAX(creado_en)) as ultima_consulta
		FROM consultas_ia`

	var totalConsultas, documentosConsultados, consultasAltaConfianza int
	var confianzaPromedio, tiempoPromedioMs float64
	var ultimaConsulta sql.NullTime

	err := r.db.QueryRow(query).Scan(
		&totalConsultas,
		&documentosConsultados,
		&confianzaPromedio,
		&tiempoPromedioMs,
		&consultasAltaConfianza,
		&ultimaConsulta,
	)

	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %v", err)
	}

	estadisticas := map[string]interface{}{
		"total_consultas":          totalConsultas,
		"documentos_consultados":   documentosConsultados,
		"confianza_promedio":       fmt.Sprintf("%.2f", confianzaPromedio),
		"tiempo_promedio_ms":       fmt.Sprintf("%.2f", tiempoPromedioMs),
		"consultas_alta_confianza": consultasAltaConfianza,
	}

	if ultimaConsulta.Valid {
		estadisticas["ultima_consulta"] = ultimaConsulta.Time
	} else {
		estadisticas["ultima_consulta"] = nil
	}

	return estadisticas, nil
}
