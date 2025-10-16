package repository

import (
	"database/sql"
	"gestion/internal/models"
)

type UsuarioRepository struct {
	db *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

// GetEdificiosByEmail obtiene todos los edificios asociados a un usuario por su email
func (r *UsuarioRepository) GetEdificiosByEmail(email string) ([]models.Edificio, error) {
	query := `
		SELECT DISTINCT e.id, e.nombre, e.direccion, e.creado_en
		FROM edificios e
		INNER JOIN usuarios_edificios ue ON e.id = ue.edificio_id
		INNER JOIN usuarios u ON ue.usuario_id = u.id
		WHERE u.email = $1
		ORDER BY e.nombre
	`

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edificios []models.Edificio
	for rows.Next() {
		var edificio models.Edificio
		err := rows.Scan(&edificio.ID, &edificio.Nombre, &edificio.Direccion, &edificio.CreadoEn)
		if err != nil {
			return nil, err
		}
		edificios = append(edificios, edificio)
	}

	return edificios, nil
}

// GetUsuarioByEmail obtiene un usuario por su email
func (r *UsuarioRepository) GetUsuarioByEmail(email string) (*models.Usuario, error) {
	query := `SELECT id, username, email, creado_en FROM usuarios WHERE email = $1`

	var usuario models.Usuario
	err := r.db.QueryRow(query, email).Scan(&usuario.ID, &usuario.Username, &usuario.Email, &usuario.CreadoEn)
	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

// VerificarAccesoEdificio verifica si un usuario tiene acceso a un edificio específico
func (r *UsuarioRepository) VerificarAccesoEdificio(email string, edificioID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM usuarios u
			INNER JOIN usuarios_edificios ue ON u.id = ue.usuario_id
			WHERE u.email = $1 AND ue.edificio_id = $2
		)
	`

	var tieneAcceso bool
	err := r.db.QueryRow(query, email, edificioID).Scan(&tieneAcceso)
	if err != nil {
		return false, err
	}

	return tieneAcceso, nil
}

// VerificarAccesoActivo verifica si un usuario tiene acceso a un activo específico
// (a través de la relación edificio -> activo)
func (r *UsuarioRepository) VerificarAccesoActivo(email string, activoID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM usuarios u
			INNER JOIN usuarios_edificios ue ON u.id = ue.usuario_id
			INNER JOIN activos a ON a.edificio_id = ue.edificio_id
			WHERE u.email = $1 AND a.id = $2
		)
	`

	var tieneAcceso bool
	err := r.db.QueryRow(query, email, activoID).Scan(&tieneAcceso)
	if err != nil {
		return false, err
	}

	return tieneAcceso, nil
}
