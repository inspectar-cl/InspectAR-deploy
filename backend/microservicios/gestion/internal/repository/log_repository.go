package repository

import (
	"database/sql"
	"encoding/json"
	"gestion/internal/models"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{db: db}
}

// CrearLog registra una acción en la tabla de auditoría
func (r *LogRepository) CrearLog(log models.LogAuditoria) error {
	query := `
		INSERT INTO logs_auditoria (
			usuario_id, usuario_email, accion, entidad, entidad_id, 
			datos_anteriores, datos_nuevos, descripcion,
			ip_origen, user_agent, fecha_accion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
	`

	// Convertir maps a JSONB (o NULL si están vacíos)
	var datosAnterioresJSON interface{}
	var datosNuevosJSON interface{}

	if log.DatosAnteriores != nil && len(log.DatosAnteriores) > 0 {
		datosAnterioresJSON, _ = json.Marshal(log.DatosAnteriores)
	} else {
		datosAnterioresJSON = nil
	}

	if log.DatosNuevos != nil && len(log.DatosNuevos) > 0 {
		datosNuevosJSON, _ = json.Marshal(log.DatosNuevos)
	} else {
		datosNuevosJSON = nil
	}

	_, err := r.db.Exec(
		query,
		log.UsuarioID,
		log.UsuarioEmail,
		log.Accion,
		log.Entidad,
		log.EntidadID,
		datosAnterioresJSON,
		datosNuevosJSON,
		log.Descripcion,
		log.IPOrigen,
		log.UserAgent,
	)

	return err
}

// ObtenerLogsPorUsuario obtiene todos los logs de un usuario
func (r *LogRepository) ObtenerLogsPorUsuario(usuarioID int, limit, offset int) ([]models.LogAuditoria, error) {
	query := `
		SELECT id, usuario_id, usuario_email, accion, entidad, entidad_id,
			   datos_anteriores, datos_nuevos, descripcion,
			   ip_origen, user_agent, fecha_accion
		FROM logs_auditoria
		WHERE usuario_id = $1
		ORDER BY fecha_accion DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, usuarioID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// ObtenerLogsPorEntidad obtiene todos los logs de una entidad específica
func (r *LogRepository) ObtenerLogsPorEntidad(entidad string, entidadID int) ([]models.LogAuditoria, error) {
	query := `
		SELECT id, usuario_id, usuario_email, accion, entidad, entidad_id,
			   datos_anteriores, datos_nuevos, descripcion,
			   ip_origen, user_agent, fecha_accion
		FROM logs_auditoria
		WHERE entidad = $1 AND entidad_id = $2
		ORDER BY fecha_accion DESC
	`

	rows, err := r.db.Query(query, entidad, entidadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// ObtenerTodosLosLogs obtiene todos los logs con paginación
func (r *LogRepository) ObtenerTodosLosLogs(limit, offset int) ([]models.LogAuditoria, error) {
	query := `
		SELECT id, usuario_id, usuario_email, accion, entidad, entidad_id,
			   datos_anteriores, datos_nuevos, descripcion,
			   ip_origen, user_agent, fecha_accion
		FROM logs_auditoria
		ORDER BY fecha_accion DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// scanLogs helper para escanear filas de logs
func (r *LogRepository) scanLogs(rows *sql.Rows) ([]models.LogAuditoria, error) {
	logs := []models.LogAuditoria{}

	for rows.Next() {
		var log models.LogAuditoria
		var datosAnterioresJSON, datosNuevosJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UsuarioID,
			&log.UsuarioEmail,
			&log.Accion,
			&log.Entidad,
			&log.EntidadID,
			&datosAnterioresJSON,
			&datosNuevosJSON,
			&log.Descripcion,
			&log.IPOrigen,
			&log.UserAgent,
			&log.FechaAccion,
		)
		if err != nil {
			return nil, err
		}

		// Decodificar JSON
		if datosAnterioresJSON != nil {
			json.Unmarshal(datosAnterioresJSON, &log.DatosAnteriores)
		}
		if datosNuevosJSON != nil {
			json.Unmarshal(datosNuevosJSON, &log.DatosNuevos)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
