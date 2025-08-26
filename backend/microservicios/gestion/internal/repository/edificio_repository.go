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
