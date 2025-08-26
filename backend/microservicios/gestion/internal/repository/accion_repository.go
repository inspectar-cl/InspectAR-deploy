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

// Obtener acciones por técnico
func (r *AccionMantenimientoRepository) GetByTecnico(tecnicoID int) ([]models.AccionConDetalles, error) {
	query := `
		SELECT 
			am.id, am.activo_id, am.tecnico_id, am.tipo, am.descripcion, am.estado, am.prioridad, 
			am.fecha_inicio, am.fecha_fin, am.creado_en,
			a.id, a.activo_id, a.nombre, a.tipo, a.estado, a.ubicacion, a.edificio_id, a.creado_en,
			t.id, t.nombre, t.email, t.telefono, t.especialidad, t.autorizado, t.creado_en
		FROM acciones_mantenimiento am
		JOIN activos a ON am.activo_id = a.id
		JOIN tecnicos t ON am.tecnico_id = t.id
		WHERE am.tecnico_id = $1 
		ORDER BY am.creado_en DESC`

	rows, err := r.db.Query(query, tecnicoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acciones []models.AccionConDetalles
	for rows.Next() {
		var accion models.AccionConDetalles
		err := rows.Scan(
			&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
			&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
			&accion.FechaFin, &accion.CreadoEn,
			&accion.Activo.ID, &accion.Activo.ActivoID, &accion.Activo.Nombre, &accion.Activo.Tipo,
			&accion.Activo.Estado, &accion.Activo.Ubicacion, &accion.Activo.EdificioID, &accion.Activo.CreadoEn,
			&accion.Tecnico.ID, &accion.Tecnico.Nombre, &accion.Tecnico.Email, &accion.Tecnico.Telefono,
			&accion.Tecnico.Especialidad, &accion.Tecnico.Autorizado, &accion.Tecnico.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		acciones = append(acciones, accion)
	}
	return acciones, nil
}

// Obtener acciones por activo con detalles
func (r *AccionMantenimientoRepository) GetByActivo(activoID int) ([]models.AccionConDetalles, error) {
	query := `
		SELECT 
			am.id, am.activo_id, am.tecnico_id, am.tipo, am.descripcion, am.estado, am.prioridad, 
			am.fecha_inicio, am.fecha_fin, am.creado_en,
			a.id, a.activo_id, a.nombre, a.tipo, a.estado, a.ubicacion, a.edificio_id, a.creado_en,
			t.id, t.nombre, t.email, t.telefono, t.especialidad, t.autorizado, t.creado_en
		FROM acciones_mantenimiento am
		JOIN activos a ON am.activo_id = a.id
		JOIN tecnicos t ON am.tecnico_id = t.id
		WHERE am.activo_id = $1 
		ORDER BY am.creado_en DESC`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acciones []models.AccionConDetalles
	for rows.Next() {
		var accion models.AccionConDetalles
		err := rows.Scan(
			&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
			&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
			&accion.FechaFin, &accion.CreadoEn,
			&accion.Activo.ID, &accion.Activo.ActivoID, &accion.Activo.Nombre, &accion.Activo.Tipo,
			&accion.Activo.Estado, &accion.Activo.Ubicacion, &accion.Activo.EdificioID, &accion.Activo.CreadoEn,
			&accion.Tecnico.ID, &accion.Tecnico.Nombre, &accion.Tecnico.Email, &accion.Tecnico.Telefono,
			&accion.Tecnico.Especialidad, &accion.Tecnico.Autorizado, &accion.Tecnico.CreadoEn,
		)
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

// Obtener acciones pendientes con prioridad y detalles
func (r *AccionMantenimientoRepository) GetPendientesConPrioridad() ([]models.AccionConDetalles, error) {
	query := `
		SELECT 
			am.id, am.activo_id, am.tecnico_id, am.tipo, am.descripcion, am.estado, am.prioridad, 
			am.fecha_inicio, am.fecha_fin, am.creado_en,
			a.id, a.activo_id, a.nombre, a.tipo, a.estado, a.ubicacion, a.edificio_id, a.creado_en,
			t.id, t.nombre, t.email, t.telefono, t.especialidad, t.autorizado, t.creado_en
		FROM acciones_mantenimiento am
		JOIN activos a ON am.activo_id = a.id
		JOIN tecnicos t ON am.tecnico_id = t.id
		WHERE am.estado IN ('pendiente', 'en_progreso')
		ORDER BY 
			CASE am.prioridad 
				WHEN 'critica' THEN 1 
				WHEN 'alta' THEN 2 
				WHEN 'media' THEN 3 
				WHEN 'baja' THEN 4 
			END, 
			am.creado_en ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acciones []models.AccionConDetalles
	for rows.Next() {
		var accion models.AccionConDetalles
		err := rows.Scan(
			&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
			&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
			&accion.FechaFin, &accion.CreadoEn,
			&accion.Activo.ID, &accion.Activo.ActivoID, &accion.Activo.Nombre, &accion.Activo.Tipo,
			&accion.Activo.Estado, &accion.Activo.Ubicacion, &accion.Activo.EdificioID, &accion.Activo.CreadoEn,
			&accion.Tecnico.ID, &accion.Tecnico.Nombre, &accion.Tecnico.Email, &accion.Tecnico.Telefono,
			&accion.Tecnico.Especialidad, &accion.Tecnico.Autorizado, &accion.Tecnico.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		acciones = append(acciones, accion)
	}
	return acciones, nil
}

// Obtener última acción de mantenimiento por activo
func (r *AccionMantenimientoRepository) GetUltimaAccionPorActivo(activoID int) (*models.AccionMantenimiento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, tipo, descripcion, estado, prioridad, 
		       fecha_inicio, fecha_fin, creado_en 
		FROM acciones_mantenimiento 
		WHERE activo_id = $1 
		ORDER BY creado_en DESC
		LIMIT 1`

	var accion models.AccionMantenimiento
	err := r.db.QueryRow(query, activoID).Scan(
		&accion.ID, &accion.ActivoID, &accion.TecnicoID, &accion.Tipo,
		&accion.Descripcion, &accion.Estado, &accion.Prioridad, &accion.FechaInicio,
		&accion.FechaFin, &accion.CreadoEn,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No hay acciones para este activo
	}

	if err != nil {
		return nil, err
	}

	return &accion, nil
}
