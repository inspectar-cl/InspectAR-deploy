package repository

import (
	"database/sql"
	"gestion/internal/models"
	"time"
)

type TecnicoRepository struct {
	db *sql.DB
}

func NewTecnicoRepository(db *sql.DB) *TecnicoRepository {
	return &TecnicoRepository{db: db}
}

func (r *TecnicoRepository) Create(tecnico *models.CreateTecnicoRequest) (*models.Tecnico, error) {
	query := `
		INSERT INTO tecnicos (nombre, email, telefono, especialidad)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nombre, email, telefono, especialidad, disponible, creado_en`

	var result models.Tecnico
	err := r.db.QueryRow(query, tecnico.Nombre, tecnico.Email, tecnico.Telefono, tecnico.Especialidad).
		Scan(&result.ID, &result.Nombre, &result.Email, &result.Telefono, &result.Especialidad, &result.Disponible, &result.CreadoEn)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *TecnicoRepository) GetAll() ([]models.Tecnico, error) {
	query := `SELECT id, nombre, email, telefono, especialidad, disponible, creado_en FROM tecnicos ORDER BY nombre`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tecnicos []models.Tecnico
	for rows.Next() {
		var tecnico models.Tecnico
		err := rows.Scan(&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
			&tecnico.Especialidad, &tecnico.Disponible, &tecnico.CreadoEn)
		if err != nil {
			return nil, err
		}
		tecnicos = append(tecnicos, tecnico)
	}
	return tecnicos, nil
}

func (r *TecnicoRepository) GetByID(id int) (*models.Tecnico, error) {
	query := `SELECT id, nombre, email, telefono, especialidad, disponible, creado_en FROM tecnicos WHERE id = $1`
	
	var tecnico models.Tecnico
	err := r.db.QueryRow(query, id).Scan(&tecnico.ID, &tecnico.Nombre, &tecnico.Email,
		&tecnico.Telefono, &tecnico.Especialidad, &tecnico.Disponible, &tecnico.CreadoEn)
	
	if err != nil {
		return nil, err
	}
	return &tecnico, nil
}

func (r *TecnicoRepository) UpdateDisponibilidad(id int, disponible bool) error {
	query := `UPDATE tecnicos SET disponible = $1 WHERE id = $2`
	_, err := r.db.Exec(query, disponible, id)
	return err
}
