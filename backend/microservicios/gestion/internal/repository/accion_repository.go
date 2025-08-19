package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type AccionMantenimientoRepository struct {
	db *sql.DB
}

func NewAccionMantenimientoRepository(db *sql.DB) *AccionMantenimientoRepository {
	return &AccionMantenimientoRepository{db: db}
}

func (r *AccionMantenimientoRepository) Create(accion *models.CreateAccionMantenimientoRequest) (*models.AccionMantenimiento, error) {
	query := `
		INSERT INTO acciones_mantenimiento (activo_id, tecnico_id, tipo, descripcion, prioridad)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, activo_id, tecnico_id, tipo, descripcion, estado, prioridad, fecha_inicio, fecha_fin, creado_en`

	var result models.AccionMantenimiento
	err := r.db.QueryRow(query, accion.ActivoID, accion.TecnicoID, accion.Tipo, accion.Descripcion, accion.Prioridad).
		Scan(&result.ID, &result.ActivoID, &result.TecnicoID, &result.Tipo, &result.Descripcion,
			&result.Estado, &result.Prioridad, &result.FechaInicio, &result.FechaFin, &result.CreadoEn)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AccionMantenimientoRepository) GetByActivo(activoID int) ([]models.AccionMantenimiento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, tipo, descripcion, estado, prioridad, 
		       fecha_inicio, fecha_fin, creado_en 
		FROM acciones_mantenimiento 
		WHERE activo_id = $1 
		ORDER BY creado_en DESC`
	
	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acciones []models.AccionMantenimiento
	for rows.Next() {
		var accion models.AccionMantenimiento
		err := rows.Scan(&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
			&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
			&accion.FechaFin, &accion.CreadoEn)
		if err != nil {
			return nil, err
		}
		acciones = append(acciones, accion)
	}
	return acciones, nil
}

func (r *AccionMantenimientoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE acciones_mantenimiento SET estado = $1 WHERE id = $2`
	
	if estado == "completado" {
		query = `UPDATE acciones_mantenimiento SET estado = $1, fecha_fin = CURRENT_TIMESTAMP WHERE id = $2`
	}
	
	_, err := r.db.Exec(query, estado, id)
	return err
}

func (r *AccionMantenimientoRepository) GetPendientes() ([]models.AccionMantenimiento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, tipo, descripcion, estado, prioridad, 
		       fecha_inicio, fecha_fin, creado_en 
		FROM acciones_mantenimiento 
		WHERE estado IN ('pendiente', 'en_progreso')
		ORDER BY prioridad DESC, creado_en ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acciones []models.AccionMantenimiento
	for rows.Next() {
		var accion models.AccionMantenimiento
		err := rows.Scan(&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
			&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
			&accion.FechaFin, &accion.CreadoEn)
		if err != nil {
			return nil, err
		}
		acciones = append(acciones, accion)
	}
	return acciones, nil
}
