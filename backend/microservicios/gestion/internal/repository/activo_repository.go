package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type ActivoRepository struct {
	db *sql.DB
}

func NewActivoRepository(db *sql.DB) *ActivoRepository {
	return &ActivoRepository{db: db}
}

func (r *ActivoRepository) GetByID(id int) (*models.Activo, error) {
	query := `
		SELECT id, activo_id, nombre, tipo, estado, ubicacion, edificio_id, creado_en 
		FROM activos 
		WHERE id = $1
	`

	var activo models.Activo
	err := r.db.QueryRow(query, id).Scan(
		&activo.ID,
		&activo.ActivoID,
		&activo.Nombre,
		&activo.Tipo,
		&activo.Estado,
		&activo.Ubicacion,
		&activo.EdificioID,
		&activo.CreadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &activo, nil
}

func (r *ActivoRepository) GetAll() ([]models.Activo, error) {
	query := `
		SELECT id, activo_id, nombre, tipo, estado, ubicacion, edificio_id, creado_en 
		FROM activos 
		ORDER BY nombre
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activos []models.Activo
	for rows.Next() {
		var activo models.Activo
		err := rows.Scan(
			&activo.ID,
			&activo.ActivoID,
			&activo.Nombre,
			&activo.Tipo,
			&activo.Estado,
			&activo.Ubicacion,
			&activo.EdificioID,
			&activo.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		activos = append(activos, activo)
	}

	return activos, nil
}

func (r *ActivoRepository) GetByEdificio(edificioID int) ([]models.Activo, error) {
	query := `
		SELECT id, activo_id, nombre, tipo, estado, ubicacion, edificio_id, creado_en 
		FROM activos 
		WHERE edificio_id = $1
		ORDER BY nombre
	`

	rows, err := r.db.Query(query, edificioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activos []models.Activo
	for rows.Next() {
		var activo models.Activo
		err := rows.Scan(
			&activo.ID,
			&activo.ActivoID,
			&activo.Nombre,
			&activo.Tipo,
			&activo.Estado,
			&activo.Ubicacion,
			&activo.EdificioID,
			&activo.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		activos = append(activos, activo)
	}

	return activos, nil
}
