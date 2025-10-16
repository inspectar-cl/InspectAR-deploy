package repository

import (
	"database/sql"
	"fmt"
	"gestion/internal/models"
)

type TipoFallaRepository struct {
	db *sql.DB
}

func NewTipoFallaRepository(db *sql.DB) *TipoFallaRepository {
	return &TipoFallaRepository{db: db}
}

// CrearTipoFalla crea un nuevo reporte de falla por un usuario (usando email)
func (r *TipoFallaRepository) CrearTipoFalla(email, tipo, descripcion string, idEdificio int) (*models.TipoFalla, error) {
	query := `
		INSERT INTO tipos_falla (tipo, descripcion, id_usuario, id_edificio)
		SELECT $1, $2, u.id, $3
		FROM usuarios u
		WHERE u.email = $4
		RETURNING id_falla, tipo, descripcion, fecha_publicacion, id_usuario, id_edificio, estado
	`

	var tipoFalla models.TipoFalla
	err := r.db.QueryRow(query, tipo, descripcion, idEdificio, email).Scan(
		&tipoFalla.IDFalla,
		&tipoFalla.Tipo,
		&tipoFalla.Descripcion,
		&tipoFalla.FechaPublicacion,
		&tipoFalla.IDUsuario,
		&tipoFalla.IDEdificio,
		&tipoFalla.Estado,
	)

	if err != nil {
		return nil, fmt.Errorf("error al crear tipo_falla: %v", err)
	}

	return &tipoFalla, nil
}

// CrearComentario crea un nuevo comentario sobre una falla (usando email)
func (r *TipoFallaRepository) CrearComentario(email string, idFalla int, comentario string) (*models.Comentario, error) {
	query := `
		INSERT INTO comentarios (id_falla, id_usuario, comentario)
		SELECT $1, u.id, $2
		FROM usuarios u
		WHERE u.email = $3
		RETURNING id_comentario, id_falla, id_usuario, comentario, fecha_comentario
	`

	var com models.Comentario
	err := r.db.QueryRow(query, idFalla, comentario, email).Scan(
		&com.IDComentario,
		&com.IDFalla,
		&com.IDUsuario,
		&com.Comentario,
		&com.FechaComentario,
	)

	if err != nil {
		return nil, fmt.Errorf("error al crear comentario: %v", err)
	}

	return &com, nil
}

// ObtenerTiposFallaPorEdificio obtiene los tipos de falla con sus comentarios paginados
// Retorna los 10 últimos tipos_falla del edificio con sus comentarios
func (r *TipoFallaRepository) ObtenerTiposFallaPorEdificio(idEdificio, pagina int) ([]models.TipoFallaResponse, error) {
	offset := (pagina - 1) * 10
	limit := 10

	// Query para obtener los tipos_falla
	queryFallas := `
		SELECT 
			tf.id_falla,
			tf.tipo,
			tf.descripcion,
			tf.fecha_publicacion,
			u.username,
			tf.id_edificio,
			tf.estado
		FROM tipos_falla tf
		INNER JOIN usuarios u ON tf.id_usuario = u.id
		WHERE tf.id_edificio = $1
		ORDER BY tf.fecha_publicacion DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(queryFallas, idEdificio, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al obtener tipos_falla: %v", err)
	}
	defer rows.Close()

	var tiposFalla []models.TipoFallaResponse

	for rows.Next() {
		var tf models.TipoFallaResponse
		err := rows.Scan(
			&tf.IDFalla,
			&tf.Tipo,
			&tf.Descripcion,
			&tf.FechaPublicacion,
			&tf.Username,
			&tf.IDEdificio,
			&tf.Estado,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear tipo_falla: %v", err)
		}

		// Obtener comentarios para este tipo_falla
		comentarios, err := r.obtenerComentariosPorFalla(tf.IDFalla)
		if err != nil {
			return nil, fmt.Errorf("error al obtener comentarios: %v", err)
		}

		tf.Comentarios = comentarios
		tiposFalla = append(tiposFalla, tf)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return tiposFalla, nil
}

// ObtenerTodosTiposFallaPorEdificio obtiene TODOS los tipos de falla con sus comentarios (sin paginación)
func (r *TipoFallaRepository) ObtenerTodosTiposFallaPorEdificio(idEdificio int) ([]models.TipoFallaResponse, error) {
	// Query para obtener todos los tipos_falla sin límite
	queryFallas := `
		SELECT 
			tf.id_falla,
			tf.tipo,
			tf.descripcion,
			tf.fecha_publicacion,
			u.username,
			tf.id_edificio,
			tf.estado
		FROM tipos_falla tf
		INNER JOIN usuarios u ON tf.id_usuario = u.id
		WHERE tf.id_edificio = $1
		ORDER BY tf.fecha_publicacion DESC
	`

	rows, err := r.db.Query(queryFallas, idEdificio)
	if err != nil {
		return nil, fmt.Errorf("error al obtener tipos_falla: %v", err)
	}
	defer rows.Close()

	var tiposFalla []models.TipoFallaResponse

	for rows.Next() {
		var tf models.TipoFallaResponse
		err := rows.Scan(
			&tf.IDFalla,
			&tf.Tipo,
			&tf.Descripcion,
			&tf.FechaPublicacion,
			&tf.Username,
			&tf.IDEdificio,
			&tf.Estado,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear tipo_falla: %v", err)
		}

		// Obtener comentarios para este tipo_falla
		comentarios, err := r.obtenerComentariosPorFalla(tf.IDFalla)
		if err != nil {
			return nil, fmt.Errorf("error al obtener comentarios: %v", err)
		}

		tf.Comentarios = comentarios
		tiposFalla = append(tiposFalla, tf)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return tiposFalla, nil
}

// obtenerComentariosPorFalla obtiene todos los comentarios de una falla específica
func (r *TipoFallaRepository) obtenerComentariosPorFalla(idFalla int) ([]models.ComentarioResponse, error) {
	queryComentarios := `
		SELECT 
			c.id_comentario,
			u.username,
			c.comentario,
			c.fecha_comentario
		FROM comentarios c
		INNER JOIN usuarios u ON c.id_usuario = u.id
		WHERE c.id_falla = $1
		ORDER BY c.fecha_comentario ASC
	`

	rows, err := r.db.Query(queryComentarios, idFalla)
	if err != nil {
		return nil, fmt.Errorf("error al obtener comentarios: %v", err)
	}
	defer rows.Close()

	var comentarios []models.ComentarioResponse

	for rows.Next() {
		var com models.ComentarioResponse
		err := rows.Scan(
			&com.IDComentario,
			&com.Username,
			&com.Comentario,
			&com.FechaComentario,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear comentario: %v", err)
		}

		comentarios = append(comentarios, com)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comentarios rows: %v", err)
	}

	return comentarios, nil
}
