package repository

import (
	"database/sql"
	"fmt"
	"gestion/internal/models"
	"strings"
)

type SolicitudRepository struct {
	db *sql.DB
}

func NewSolicitudRepository(db *sql.DB) *SolicitudRepository {
	return &SolicitudRepository{db: db}
}

// CreateSolicitud crea una nueva solicitud
func (r *SolicitudRepository) CreateSolicitud(req *models.CreateSolicitudRequest) (*models.SolicitudTecnico, error) {
	query := `
		INSERT INTO solicitudes_tecnico (
			tecnico_id, residente_id, activo_id, edificio_id, tipo, asunto, 
			descripcion, prioridad, medio_contacto, telefono_contacto, 
			email_contacto, estado, fecha_creacion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pendiente', NOW())
		RETURNING id, tecnico_id, residente_id, activo_id, edificio_id, tipo, 
				  asunto, descripcion, prioridad, estado, fecha_creacion, 
				  medio_contacto, telefono_contacto, email_contacto`

	var solicitud models.SolicitudTecnico
	err := r.db.QueryRow(
		query,
		req.TecnicoID, req.ResidenteID, req.ActivoID, req.EdificioID,
		req.Tipo, req.Asunto, req.Descripcion, req.Prioridad,
		req.MedioContacto, req.TelefonoContacto, req.EmailContacto,
	).Scan(
		&solicitud.ID, &solicitud.TecnicoID, &solicitud.ResidenteID,
		&solicitud.ActivoID, &solicitud.EdificioID, &solicitud.Tipo,
		&solicitud.Asunto, &solicitud.Descripcion, &solicitud.Prioridad,
		&solicitud.Estado, &solicitud.FechaCreacion, &solicitud.MedioContacto,
		&solicitud.TelefonoContacto, &solicitud.EmailContacto,
	)

	if err != nil {
		return nil, fmt.Errorf("error al crear solicitud: %w", err)
	}

	// Cargar información del técnico
	if err := r.loadTecnicoInfo(&solicitud); err != nil {
		return nil, fmt.Errorf("error al cargar información del técnico: %w", err)
	}

	return &solicitud, nil
}

// GetSolicitudByID obtiene una solicitud por su ID
func (r *SolicitudRepository) GetSolicitudByID(id int) (*models.SolicitudTecnico, error) {
	query := `
		SELECT s.id, s.tecnico_id, s.residente_id, s.activo_id, s.edificio_id,
			   s.tipo, s.asunto, s.descripcion, s.prioridad, s.estado,
			   s.fecha_creacion, s.fecha_envio, s.fecha_recepcion, s.fecha_completado,
			   s.medio_contacto, s.telefono_contacto, s.email_contacto,
			   s.respuesta_tecnico, s.notas_internas,
			   t.id, t.nombre, t.apellido, t.email, t.telefono, t.especialidad, t.empresa_id,
			   e.id, e.nombre, e.rut, e.telefono
		FROM solicitudes_tecnico s
		JOIN tecnicos t ON s.tecnico_id = t.id
		JOIN empresas e ON t.empresa_id = e.id
		WHERE s.id = $1`

	var solicitud models.SolicitudTecnico
	var tecnico models.Tecnico
	var empresa models.Empresa

	err := r.db.QueryRow(query, id).Scan(
		&solicitud.ID, &solicitud.TecnicoID, &solicitud.ResidenteID,
		&solicitud.ActivoID, &solicitud.EdificioID, &solicitud.Tipo,
		&solicitud.Asunto, &solicitud.Descripcion, &solicitud.Prioridad,
		&solicitud.Estado, &solicitud.FechaCreacion, &solicitud.FechaEnvio,
		&solicitud.FechaRecepcion, &solicitud.FechaCompletado,
		&solicitud.MedioContacto, &solicitud.TelefonoContacto, &solicitud.EmailContacto,
		&solicitud.RespuestaTecnico, &solicitud.NotasInternas,
		&tecnico.ID, &tecnico.Nombre, &tecnico.Apellido, &tecnico.Email, &tecnico.Telefono,
		&tecnico.Especialidad, &tecnico.EmpresaID,
		&empresa.ID, &empresa.Nombre, &empresa.RUT, &empresa.Telefono,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("solicitud no encontrada")
		}
		return nil, fmt.Errorf("error al obtener solicitud: %w", err)
	}

	// Asignar datos del técnico y empresa
	tecnico.Empresa = empresa
	solicitud.Tecnico = tecnico

	return &solicitud, nil
}

// GetSolicitudesByFilter obtiene solicitudes aplicando filtros
func (r *SolicitudRepository) GetSolicitudesByFilter(filter models.SolicitudFilter, limit, offset int) ([]models.SolicitudTecnico, int64, error) {
	// Construir query base
	baseQuery := `
		FROM solicitudes_tecnico s
		JOIN tecnicos t ON s.tecnico_id = t.id
		JOIN empresas e ON t.empresa_id = e.id`

	// Construir condiciones WHERE
	whereConditions, args := r.buildWhereConditions(filter)
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Contar total de registros
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("error al contar solicitudes: %w", err)
	}

	// Obtener registros con paginación
	selectQuery := `
		SELECT s.id, s.tecnico_id, s.residente_id, s.activo_id, s.edificio_id,
			   s.tipo, s.asunto, s.descripcion, s.prioridad, s.estado,
			   s.fecha_creacion, s.fecha_envio, s.fecha_recepcion, s.fecha_completado,
			   s.medio_contacto, s.telefono_contacto, s.email_contacto,
			   s.respuesta_tecnico,
			   t.nombre, t.apellido, t.email, t.telefono, t.especialidad,
			   e.nombre as empresa_nombre, e.rut ` +
		baseQuery + whereClause +
		` ORDER BY s.fecha_creacion DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) +
		` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	args = append(args, limit, offset)
	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener solicitudes: %w", err)
	}
	defer rows.Close()

	var solicitudes []models.SolicitudTecnico
	for rows.Next() {
		var solicitud models.SolicitudTecnico
		var tecnico models.Tecnico
		var empresa models.Empresa

		err := rows.Scan(
			&solicitud.ID, &solicitud.TecnicoID, &solicitud.ResidenteID,
			&solicitud.ActivoID, &solicitud.EdificioID, &solicitud.Tipo,
			&solicitud.Asunto, &solicitud.Descripcion, &solicitud.Prioridad,
			&solicitud.Estado, &solicitud.FechaCreacion, &solicitud.FechaEnvio,
			&solicitud.FechaRecepcion, &solicitud.FechaCompletado,
			&solicitud.MedioContacto, &solicitud.TelefonoContacto, &solicitud.EmailContacto,
			&solicitud.RespuestaTecnico, &tecnico.Nombre, &tecnico.Apellido,
			&tecnico.Email, &tecnico.Telefono, &tecnico.Especialidad,
			&empresa.Nombre, &empresa.RUT,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al escanear solicitud: %w", err)
		}

		tecnico.Empresa = empresa
		solicitud.Tecnico = tecnico
		solicitudes = append(solicitudes, solicitud)
	}

	return solicitudes, total, nil
}

// MarkSolicitudAsEnviada marca una solicitud como enviada
func (r *SolicitudRepository) MarkSolicitudAsEnviada(id int) error {
	query := `
		UPDATE solicitudes_tecnico 
		SET estado = 'enviada', fecha_envio = NOW() 
		WHERE id = $1 AND estado = 'pendiente'`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al marcar solicitud como enviada: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("solicitud no encontrada o ya enviada")
	}

	return nil
}

// UpdateSolicitudEstado actualiza el estado de una solicitud
func (r *SolicitudRepository) UpdateSolicitudEstado(id int, estado models.EstadoSolicitud, respuesta string) error {
	var query string
	var args []interface{}

	switch estado {
	case models.EstadoRecibida:
		query = `UPDATE solicitudes_tecnico SET estado = $1, fecha_recepcion = NOW() WHERE id = $2`
		args = []interface{}{estado, id}
	case models.EstadoCompletada:
		query = `UPDATE solicitudes_tecnico SET estado = $1, fecha_completado = NOW(), respuesta_tecnico = $2 WHERE id = $3`
		args = []interface{}{estado, respuesta, id}
	default:
		query = `UPDATE solicitudes_tecnico SET estado = $1 WHERE id = $2`
		args = []interface{}{estado, id}
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de solicitud: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("solicitud no encontrada")
	}

	return nil
}

// GetEstadisticasSolicitudes obtiene estadísticas de solicitudes
func (r *SolicitudRepository) GetEstadisticasSolicitudes(tecnicoID *int, edificioID *int) (*models.EstadisticasSolicitudes, error) {
	baseQuery := "FROM solicitudes_tecnico s"
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if tecnicoID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("s.tecnico_id = $%d", argIndex))
		args = append(args, *tecnicoID)
		argIndex++
	}

	if edificioID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("s.edificio_id = $%d", argIndex))
		args = append(args, *edificioID)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	var stats models.EstadisticasSolicitudes

	// Total de solicitudes
	totalQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	if err := r.db.QueryRow(totalQuery, args...).Scan(&stats.TotalSolicitudes); err != nil {
		return nil, fmt.Errorf("error al contar total de solicitudes: %w", err)
	}

	// Solicitudes por estado
	estadoQuery := `
		SELECT estado, COUNT(*) 
		` + baseQuery + whereClause + `
		GROUP BY estado`

	rows, err := r.db.Query(estadoQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener conteos por estado: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var estado models.EstadoSolicitud
		var count int
		if err := rows.Scan(&estado, &count); err != nil {
			return nil, fmt.Errorf("error al escanear estado: %w", err)
		}

		switch estado {
		case models.EstadoPendiente:
			stats.SolicitudesPendientes = count
		case models.EstadoEnProceso:
			stats.SolicitudesEnProceso = count
		case models.EstadoCompletada:
			stats.SolicitudesCompletadas = count
		case models.EstadoCancelada:
			stats.SolicitudesCanceladas = count
		}
	}

	return &stats, nil
}

// Helper functions

func (r *SolicitudRepository) loadTecnicoInfo(solicitud *models.SolicitudTecnico) error {
	query := `
		SELECT t.nombre, t.apellido, t.email, t.telefono, t.especialidad,
			   e.id, e.nombre, e.rut, e.telefono
		FROM tecnicos t
		JOIN empresas e ON t.empresa_id = e.id
		WHERE t.id = $1`

	var tecnico models.Tecnico
	var empresa models.Empresa

	err := r.db.QueryRow(query, solicitud.TecnicoID).Scan(
		&tecnico.Nombre, &tecnico.Apellido, &tecnico.Email, &tecnico.Telefono,
		&tecnico.Especialidad, &empresa.ID, &empresa.Nombre, &empresa.RUT, &empresa.Telefono,
	)

	if err != nil {
		return err
	}

	tecnico.ID = solicitud.TecnicoID
	tecnico.Empresa = empresa
	solicitud.Tecnico = tecnico
	return nil
}

func (r *SolicitudRepository) buildWhereConditions(filter models.SolicitudFilter) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.TecnicoID != nil {
		conditions = append(conditions, fmt.Sprintf("s.tecnico_id = $%d", argIndex))
		args = append(args, *filter.TecnicoID)
		argIndex++
	}

	if filter.ResidenteID != nil {
		conditions = append(conditions, fmt.Sprintf("s.residente_id = $%d", argIndex))
		args = append(args, *filter.ResidenteID)
		argIndex++
	}

	if filter.ActivoID != nil {
		conditions = append(conditions, fmt.Sprintf("s.activo_id = $%d", argIndex))
		args = append(args, *filter.ActivoID)
		argIndex++
	}

	if filter.EdificioID != nil {
		conditions = append(conditions, fmt.Sprintf("s.edificio_id = $%d", argIndex))
		args = append(args, *filter.EdificioID)
		argIndex++
	}

	if filter.Tipo != nil {
		conditions = append(conditions, fmt.Sprintf("s.tipo = $%d", argIndex))
		args = append(args, *filter.Tipo)
		argIndex++
	}

	if filter.Estado != nil {
		conditions = append(conditions, fmt.Sprintf("s.estado = $%d", argIndex))
		args = append(args, *filter.Estado)
		argIndex++
	}

	if filter.Prioridad != nil {
		conditions = append(conditions, fmt.Sprintf("s.prioridad = $%d", argIndex))
		args = append(args, *filter.Prioridad)
		argIndex++
	}

	if filter.FechaDesde != nil {
		conditions = append(conditions, fmt.Sprintf("s.fecha_creacion >= $%d", argIndex))
		args = append(args, *filter.FechaDesde)
		argIndex++
	}

	if filter.FechaHasta != nil {
		conditions = append(conditions, fmt.Sprintf("s.fecha_creacion <= $%d", argIndex))
		args = append(args, *filter.FechaHasta)
		argIndex++
	}

	if filter.Busqueda != nil && *filter.Busqueda != "" {
		searchTerm := "%" + *filter.Busqueda + "%"
		conditions = append(conditions, fmt.Sprintf("(s.asunto ILIKE $%d OR s.descripcion ILIKE $%d)", argIndex, argIndex))
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}
