package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type ReporteRepository struct {
	db *sql.DB
}

func NewReporteRepository(db *sql.DB) *ReporteRepository {
	return &ReporteRepository{db: db}
}

func (r *ReporteRepository) Create(reporte *models.Reporte) (*models.Reporte, error) {
	query := `
		INSERT INTO reportes (activo_id, tipo_reporte, contenido, generado_en, estado)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, creado_en
	`

	err := r.db.QueryRow(
		query,
		reporte.ActivoID,
		reporte.TipoReporte,
		reporte.Contenido,
		reporte.GeneradoEn,
		reporte.Estado,
	).Scan(&reporte.ID, &reporte.GeneradoEn)

	if err != nil {
		return nil, err
	}

	return reporte, nil
}

func (r *ReporteRepository) GetByActivo(activoID int) ([]models.ReporteCompleto, error) {
	query := `
		SELECT 
			r.id, r.activo_id, r.tipo_reporte, r.contenido, r.generado_en, r.estado,
			a.id, a.activo_id, a.nombre, a.tipo, a.estado, a.ubicacion, a.edificio_id, a.creado_en,
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
		err := rows.Scan(
			&reporte.ID,
			&reporte.ActivoID,
			&reporte.TipoReporte,
			&reporte.Contenido,
			&reporte.GeneradoEn,
			&reporte.Estado,
			&reporte.ActivoConEdificio.ID,
			&reporte.ActivoConEdificio.ActivoID,
			&reporte.ActivoConEdificio.Nombre,
			&reporte.ActivoConEdificio.Tipo,
			&reporte.ActivoConEdificio.Estado,
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
		reportes = append(reportes, reporte)
	}

	return reportes, nil
}
