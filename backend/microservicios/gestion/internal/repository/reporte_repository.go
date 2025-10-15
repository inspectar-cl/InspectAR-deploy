package repository

import (
	"database/sql"
	"encoding/json"
	"gestion/internal/models"
	"time"
)

type ReporteRepository struct {
	db *sql.DB
}

func NewReporteRepository(db *sql.DB) *ReporteRepository {
	return &ReporteRepository{db: db}
}

func (r *ReporteRepository) Create(reporte *models.Reporte) (*models.Reporte, error) {
	// Serializar JSON fields
	estructuraJSON, _ := json.Marshal(reporte.EstructuraInforme)
	metadataJSON, _ := json.Marshal(reporte.MetadataInforme)

	query := `
		INSERT INTO reportes (
			activo_id, tipo_reporte, contenido, 
			observaciones_analista, autor_analista, 
			estructura_informe, metadata_informe, version_reporte,
			estado_revision, generado_en, estado
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, generado_en
	`

	err := r.db.QueryRow(
		query,
		reporte.ActivoID,
		reporte.TipoReporte,
		reporte.Contenido,
		reporte.ObservacionesAnalista,
		reporte.AutorAnalista,
		estructuraJSON,
		metadataJSON,
		reporte.VersionReporte,
		reporte.EstadoRevision,
		reporte.GeneradoEn,
		reporte.Estado,
	).Scan(&reporte.ID, &reporte.GeneradoEn)

	if err != nil {
		return nil, err
	}

	return reporte, nil
}

func (r *ReporteRepository) GetByID(id int) (*models.Reporte, error) {
	query := `
		SELECT 
			id, activo_id, tipo_reporte, contenido,
			observaciones_analista, autor_analista,
			estructura_informe, metadata_informe, version_reporte,
			estado_revision, fecha_revision, revisor,
			generado_en, estado
		FROM reportes
		WHERE id = $1
	`

	var reporte models.Reporte
	var estructuraJSON, metadataJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&reporte.ID,
		&reporte.ActivoID,
		&reporte.TipoReporte,
		&reporte.Contenido,
		&reporte.ObservacionesAnalista,
		&reporte.AutorAnalista,
		&estructuraJSON,
		&metadataJSON,
		&reporte.VersionReporte,
		&reporte.EstadoRevision,
		&reporte.FechaRevision,
		&reporte.Revisor,
		&reporte.GeneradoEn,
		&reporte.Estado,
	)

	if err != nil {
		return nil, err
	}

	// Deserializar JSON fields
	json.Unmarshal(estructuraJSON, &reporte.EstructuraInforme)
	json.Unmarshal(metadataJSON, &reporte.MetadataInforme)

	return &reporte, nil
}

func (r *ReporteRepository) UpdateObservaciones(id int, observaciones, autor string) error {
	query := `
		UPDATE reportes 
		SET observaciones_analista = $1, autor_analista = $2, version_reporte = version_reporte + 1
		WHERE id = $3
	`

	_, err := r.db.Exec(query, observaciones, autor, id)
	return err
}

func (r *ReporteRepository) UpdateEstadoRevision(id int, estado, revisor, observaciones string) error {
	query := `
		UPDATE reportes 
		SET estado_revision = $1, revisor = $2, fecha_revision = $3,
		    observaciones_analista = CASE 
		        WHEN $4 != '' THEN $4 
		        ELSE observaciones_analista 
		    END
		WHERE id = $5
	`

	_, err := r.db.Exec(query, estado, revisor, time.Now(), observaciones, id)
	return err
}

func (r *ReporteRepository) UpdateEstructuraInforme(id int, estructura map[string]interface{}) error {
	estructuraJSON, _ := json.Marshal(estructura)

	query := `
		UPDATE reportes 
		SET estructura_informe = $1, version_reporte = version_reporte + 1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, estructuraJSON, id)
	return err
}

func (r *ReporteRepository) GetByActivo(activoID int) ([]models.ReporteCompleto, error) {
	query := `
		SELECT 
			r.id, r.activo_id, r.tipo_reporte, r.contenido,
			r.observaciones_analista, r.autor_analista,
			r.estructura_informe, r.metadata_informe, r.version_reporte,
			r.estado_revision, r.fecha_revision, r.revisor,
			r.generado_en, r.estado,
			a.id, a.nombre, a.tipo, a.ubicacion, a.edificio_id, a.creado_en,
			e.id, e.nombre, e.direccion, e.creado_en
		FROM reportes r
		JOIN activos a ON r.activo_id = a.id
		JOIN edificios e ON a.edificio_id = e.id
		WHERE r.activo_id = $1
		ORDER BY r.generado_en DESC
	`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []models.ReporteCompleto
	for rows.Next() {
		var reporte models.ReporteCompleto
		var estructuraJSON, metadataJSON []byte

		err := rows.Scan(
			&reporte.ID,
			&reporte.ActivoID,
			&reporte.TipoReporte,
			&reporte.Contenido,
			&reporte.ObservacionesAnalista,
			&reporte.AutorAnalista,
			&estructuraJSON,
			&metadataJSON,
			&reporte.VersionReporte,
			&reporte.EstadoRevision,
			&reporte.FechaRevision,
			&reporte.Revisor,
			&reporte.GeneradoEn,
			&reporte.Estado,
			&reporte.ActivoConEdificio.ID,
			&reporte.ActivoConEdificio.Nombre,
			&reporte.ActivoConEdificio.Tipo,
			&reporte.ActivoConEdificio.Ubicacion,
			&reporte.ActivoConEdificio.EdificioID,
			&reporte.ActivoConEdificio.CreadoEn,
			&reporte.ActivoConEdificio.Edificio.ID,
			&reporte.ActivoConEdificio.Edificio.Nombre,
			&reporte.ActivoConEdificio.Edificio.Direccion,
			&reporte.ActivoConEdificio.Edificio.CreadoEn,
		)
		if err != nil {
			return nil, err
		}

		// Deserializar JSON fields
		json.Unmarshal(estructuraJSON, &reporte.EstructuraInforme)
		json.Unmarshal(metadataJSON, &reporte.MetadataInforme)

		reportes = append(reportes, reporte)
	}

	return reportes, nil
}

func (r *ReporteRepository) GetConObservaciones(activoID int) ([]models.ReporteConObservaciones, error) {
	query := `
		SELECT 
			r.id, r.activo_id, r.tipo_reporte, r.contenido,
			r.observaciones_analista, r.autor_analista,
			r.estructura_informe, r.metadata_informe, r.version_reporte,
			r.estado_revision, r.fecha_revision, r.revisor,
			r.generado_en, r.estado,
			a.nombre, a.tipo
		FROM reportes r
		JOIN activos a ON r.activo_id = a.id
		WHERE r.activo_id = $1
		ORDER BY r.generado_en DESC
	`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []models.ReporteConObservaciones
	for rows.Next() {
		var reporte models.ReporteConObservaciones
		var estructuraJSON, metadataJSON []byte

		err := rows.Scan(
			&reporte.ID,
			&reporte.ActivoID,
			&reporte.TipoReporte,
			&reporte.Contenido,
			&reporte.ObservacionesAnalista,
			&reporte.AutorAnalista,
			&estructuraJSON,
			&metadataJSON,
			&reporte.VersionReporte,
			&reporte.EstadoRevision,
			&reporte.FechaRevision,
			&reporte.Revisor,
			&reporte.GeneradoEn,
			&reporte.Estado,
			&reporte.NombreActivo,
			&reporte.TipoActivo,
		)
		if err != nil {
			return nil, err
		}

		// Deserializar JSON fields
		json.Unmarshal(estructuraJSON, &reporte.EstructuraInforme)
		json.Unmarshal(metadataJSON, &reporte.MetadataInforme)

		reportes = append(reportes, reporte)
	}

	return reportes, nil
}

func (r *ReporteRepository) GetAll() ([]models.ReporteConObservaciones, error) {
	query := `
		SELECT 
			r.id, r.activo_id, r.tipo_reporte, r.contenido,
			r.observaciones_analista, r.autor_analista,
			r.estructura_informe, r.metadata_informe, r.version_reporte,
			r.estado_revision, r.fecha_revision, r.revisor,
			r.generado_en, r.estado,
			a.nombre, a.tipo
		FROM reportes r
		JOIN activos a ON r.activo_id = a.id
		ORDER BY r.generado_en DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportes []models.ReporteConObservaciones
	for rows.Next() {
		var reporte models.ReporteConObservaciones
		var estructuraJSON, metadataJSON []byte

		err := rows.Scan(
			&reporte.ID,
			&reporte.ActivoID,
			&reporte.TipoReporte,
			&reporte.Contenido,
			&reporte.ObservacionesAnalista,
			&reporte.AutorAnalista,
			&estructuraJSON,
			&metadataJSON,
			&reporte.VersionReporte,
			&reporte.EstadoRevision,
			&reporte.FechaRevision,
			&reporte.Revisor,
			&reporte.GeneradoEn,
			&reporte.Estado,
			&reporte.NombreActivo,
			&reporte.TipoActivo,
		)
		if err != nil {
			return nil, err
		}

		// Deserializar JSON fields
		json.Unmarshal(estructuraJSON, &reporte.EstructuraInforme)
		json.Unmarshal(metadataJSON, &reporte.MetadataInforme)

		reportes = append(reportes, reporte)
	}

	return reportes, nil
}
