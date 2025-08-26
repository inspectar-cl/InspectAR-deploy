package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type TecnicoRepository struct {
	db *sql.DB
}

func NewTecnicoRepository(db *sql.DB) *TecnicoRepository {
	return &TecnicoRepository{db: db}
}

// Create - Crear técnico con especialidades
func (r *TecnicoRepository) Create(req *models.CreateTecnicoRequest) (*models.Tecnico, error) {
	query := `
		INSERT INTO tecnicos (nombre, email, telefono, especialidad, autorizado) 
		VALUES ($1, $2, $3, $4, true) 
		RETURNING id, nombre, email, telefono, especialidad, autorizado, creado_en`

	var tecnico models.Tecnico
	err := r.db.QueryRow(query, req.Nombre, req.Email, req.Telefono, req.Especialidad).Scan(
		&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
		&tecnico.Especialidad, &tecnico.Autorizado, &tecnico.CreadoEn,
	)

	return &tecnico, err
}

// GetAll - Listar todos los técnicos
func (r *TecnicoRepository) GetAll() ([]models.Tecnico, error) {
	query := `SELECT id, nombre, email, telefono, especialidad, autorizado, creado_en FROM tecnicos ORDER BY nombre`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tecnicos []models.Tecnico
	for rows.Next() {
		var tecnico models.Tecnico
		err := rows.Scan(&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
			&tecnico.Especialidad, &tecnico.Autorizado, &tecnico.CreadoEn)
		if err != nil {
			return nil, err
		}
		tecnicos = append(tecnicos, tecnico)
	}

	return tecnicos, nil
}

// GetByActivo - Listar técnicos relacionados con un activo específico
func (r *TecnicoRepository) GetByActivo(activoID int, soloAutorizados bool) ([]models.Tecnico, error) {
	query := `
		SELECT DISTINCT t.id, t.nombre, t.email, t.telefono, t.especialidad, t.autorizado, t.creado_en 
		FROM tecnicos t
		JOIN activos_tecnicos at ON t.id = at.tecnico_id
		WHERE at.activo_id = $1`

	if soloAutorizados {
		query += ` AND t.autorizado = true`
	}
	query += ` ORDER BY t.nombre`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tecnicos []models.Tecnico
	for rows.Next() {
		var tecnico models.Tecnico
		err := rows.Scan(&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
			&tecnico.Especialidad, &tecnico.Autorizado, &tecnico.CreadoEn)
		if err != nil {
			return nil, err
		}
		tecnicos = append(tecnicos, tecnico)
	}

	return tecnicos, nil
}

// GetByEdificio - Listar técnicos relacionados con un edificio (indirectamente a través de activos)
func (r *TecnicoRepository) GetByEdificio(edificioID int, soloAutorizados bool) ([]models.Tecnico, error) {
	query := `
		SELECT DISTINCT t.id, t.nombre, t.email, t.telefono, t.especialidad, t.autorizado, t.creado_en 
		FROM tecnicos t
		JOIN activos_tecnicos at ON t.id = at.tecnico_id
		JOIN activos a ON at.activo_id = a.id
		WHERE a.edificio_id = $1`

	if soloAutorizados {
		query += ` AND t.autorizado = true`
	}
	query += ` ORDER BY t.nombre`

	rows, err := r.db.Query(query, edificioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tecnicos []models.Tecnico
	for rows.Next() {
		var tecnico models.Tecnico
		err := rows.Scan(&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
			&tecnico.Especialidad, &tecnico.Autorizado, &tecnico.CreadoEn)
		if err != nil {
			return nil, err
		}
		tecnicos = append(tecnicos, tecnico)
	}

	return tecnicos, nil
}

// GetByID - Consultar información de un técnico específico
func (r *TecnicoRepository) GetByID(id int) (*models.Tecnico, error) {
	query := `SELECT id, nombre, email, telefono, especialidad, autorizado, creado_en FROM tecnicos WHERE id = $1`

	var tecnico models.Tecnico
	err := r.db.QueryRow(query, id).Scan(
		&tecnico.ID, &tecnico.Nombre, &tecnico.Email, &tecnico.Telefono,
		&tecnico.Especialidad, &tecnico.Autorizado, &tecnico.CreadoEn,
	)

	return &tecnico, err
}

// UpdateAutorizado - Actualizar estado autorizado de técnicos
func (r *TecnicoRepository) UpdateAutorizado(id int, autorizado bool) error {
	query := `UPDATE tecnicos SET autorizado = $1 WHERE id = $2`
	_, err := r.db.Exec(query, autorizado, id)
	return err
}

// AsignarTecnicoAActivo - Asignar técnico a activo (relación muchos a muchos)
func (r *TecnicoRepository) AsignarTecnicoAActivo(activoID, tecnicoID int) error {
	query := `INSERT INTO activos_tecnicos (activo_id, tecnico_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, activoID, tecnicoID)
	return err
}

// ObtenerActivosPorTecnico - Obtener activos asociados a un técnico (RUTA PRINCIPAL)
func (r *TecnicoRepository) ObtenerActivosPorTecnico(tecnicoID int) ([]models.Activo, error) {
	query := `
		SELECT DISTINCT a.id, a.activo_id, a.nombre, a.tipo, a.estado, a.ubicacion, a.edificio_id, a.creado_en
		FROM activos a
		JOIN activos_tecnicos at ON a.id = at.activo_id
		WHERE at.tecnico_id = $1
		ORDER BY a.nombre`

	rows, err := r.db.Query(query, tecnicoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activos []models.Activo
	for rows.Next() {
		var activo models.Activo

		err := rows.Scan(
			&activo.ID, &activo.ActivoID, &activo.Nombre, &activo.Tipo,
			&activo.Estado, &activo.Ubicacion, &activo.EdificioID, &activo.CreadoEn,
		)
		if err != nil {
			return nil, err
		}

		activos = append(activos, activo)
	}

	return activos, nil
}
