package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type EdificioRepository struct {
	db *sql.DB
}

func NewEdificioRepository(db *sql.DB) *EdificioRepository {
	return &EdificioRepository{db: db}
}

func (r *EdificioRepository) GetByID(id int) (*models.Edificio, error) {
	query := `
		SELECT id, nombre, direccion, creado_en 
		FROM edificios 
		WHERE id = $1
	`

	var edificio models.Edificio
	err := r.db.QueryRow(query, id).Scan(
		&edificio.ID,
		&edificio.Nombre,
		&edificio.Direccion,
		&edificio.CreadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &edificio, nil
}

func (r *EdificioRepository) GetAll() ([]models.Edificio, error) {
	query := `
		SELECT id, nombre, direccion, creado_en 
		FROM edificios 
		ORDER BY nombre
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edificios []models.Edificio
	for rows.Next() {
		var edificio models.Edificio
		err := rows.Scan(
			&edificio.ID,
			&edificio.Nombre,
			&edificio.Direccion,
			&edificio.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		edificios = append(edificios, edificio)
	}

	return edificios, nil
}

// Crear un edificio (admin)
func (r *EdificioRepository) Crear(edificio models.Edificio) (*models.Edificio, error) {
	query := `
		INSERT INTO edificios (nombre, direccion, latitud, longitud)
		VALUES ($1, $2, $3, $4)
		RETURNING id, creado_en
	`

	err := r.db.QueryRow(
		query,
		edificio.Nombre,
		edificio.Direccion,
		edificio.Latitud,
		edificio.Longitud,
	).Scan(&edificio.ID, &edificio.CreadoEn)

	if err != nil {
		return nil, err
	}

	return &edificio, nil
}

// Actualizar un edificio (admin)
func (r *EdificioRepository) Actualizar(id int, edificio models.Edificio) error {
	query := `
		UPDATE edificios 
		SET nombre = $1, direccion = $2, latitud = $3, longitud = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		query,
		edificio.Nombre,
		edificio.Direccion,
		edificio.Latitud,
		edificio.Longitud,
		id,
	)

	return err
}

// Eliminar un edificio (admin)
func (r *EdificioRepository) Eliminar(id int) error {
	query := `DELETE FROM edificios WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
