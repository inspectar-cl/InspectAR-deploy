package repository

import (
	"database/sql"
	"fmt"
	"gestion/internal/models"
)

type FirmaRepository struct {
	db *sql.DB
}

func NewFirmaRepository(db *sql.DB) *FirmaRepository {
	return &FirmaRepository{db: db}
}

// Create crea una nueva firma digital
func (r *FirmaRepository) Create(firma *models.FirmaDigital) error {
	query := `
		INSERT INTO firmas_digitales 
		(usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, datos_firma, tamano_bytes, es_predeterminada)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, creado_en, actualizado_en
	`

	err := r.db.QueryRow(
		query,
		firma.UsuarioID,
		firma.NombreArchivo,
		firma.RutaArchivo,
		firma.TipoMime,
		firma.Formato,
		firma.DatosFirma,
		firma.TamanoBytes,
		firma.EsPredeterminada,
	).Scan(&firma.ID, &firma.CreadoEn, &firma.ActualizadoEn)

	if err != nil {
		return fmt.Errorf("error al crear firma digital: %v", err)
	}

	// Si es predeterminada, desmarcar otras firmas del mismo usuario
	if firma.EsPredeterminada {
		return r.UnsetOtherDefaults(firma.UsuarioID, firma.ID)
	}

	return nil
}

// GetByID obtiene una firma por ID
func (r *FirmaRepository) GetByID(id int) (*models.FirmaDigital, error) {
	query := `
		SELECT id, usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, 
		       datos_firma, tamano_bytes, es_predeterminada, creado_en, actualizado_en
		FROM firmas_digitales
		WHERE id = $1
	`

	var firma models.FirmaDigital
	err := r.db.QueryRow(query, id).Scan(
		&firma.ID,
		&firma.UsuarioID,
		&firma.NombreArchivo,
		&firma.RutaArchivo,
		&firma.TipoMime,
		&firma.Formato,
		&firma.DatosFirma,
		&firma.TamanoBytes,
		&firma.EsPredeterminada,
		&firma.CreadoEn,
		&firma.ActualizadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &firma, nil
}

// GetByUsuarioID obtiene todas las firmas de un usuario
func (r *FirmaRepository) GetByUsuarioID(usuarioID int) ([]models.FirmaDigital, error) {
	query := `
		SELECT id, usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, 
		       datos_firma, tamano_bytes, es_predeterminada, creado_en, actualizado_en
		FROM firmas_digitales
		WHERE usuario_id = $1
		ORDER BY es_predeterminada DESC, creado_en DESC
	`

	rows, err := r.db.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var firmas []models.FirmaDigital
	for rows.Next() {
		var firma models.FirmaDigital
		err := rows.Scan(
			&firma.ID,
			&firma.UsuarioID,
			&firma.NombreArchivo,
			&firma.RutaArchivo,
			&firma.TipoMime,
			&firma.Formato,
			&firma.DatosFirma,
			&firma.TamanoBytes,
			&firma.EsPredeterminada,
			&firma.CreadoEn,
			&firma.ActualizadoEn,
		)
		if err != nil {
			return nil, err
		}
		firmas = append(firmas, firma)
	}

	return firmas, nil
}

// GetByUsuarioEmail obtiene todas las firmas de un usuario por email
func (r *FirmaRepository) GetByUsuarioEmail(email string) ([]models.FirmaDigital, error) {
	query := `
		SELECT f.id, f.usuario_id, f.nombre_archivo, f.ruta_archivo, f.tipo_mime, f.formato, 
		       f.datos_firma, f.tamano_bytes, f.es_predeterminada, f.creado_en, f.actualizado_en
		FROM firmas_digitales f
		INNER JOIN usuarios u ON f.usuario_id = u.id
		WHERE u.email = $1
		ORDER BY f.es_predeterminada DESC, f.creado_en DESC
	`

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var firmas []models.FirmaDigital
	for rows.Next() {
		var firma models.FirmaDigital
		err := rows.Scan(
			&firma.ID,
			&firma.UsuarioID,
			&firma.NombreArchivo,
			&firma.RutaArchivo,
			&firma.TipoMime,
			&firma.Formato,
			&firma.DatosFirma,
			&firma.TamanoBytes,
			&firma.EsPredeterminada,
			&firma.CreadoEn,
			&firma.ActualizadoEn,
		)
		if err != nil {
			return nil, err
		}
		firmas = append(firmas, firma)
	}

	return firmas, nil
}

// GetDefaultByUsuario obtiene la firma predeterminada de un usuario
func (r *FirmaRepository) GetDefaultByUsuario(usuarioID int) (*models.FirmaDigital, error) {
	query := `
		SELECT id, usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, 
		       datos_firma, tamano_bytes, es_predeterminada, creado_en, actualizado_en
		FROM firmas_digitales
		WHERE usuario_id = $1 AND es_predeterminada = true
		LIMIT 1
	`

	var firma models.FirmaDigital
	err := r.db.QueryRow(query, usuarioID).Scan(
		&firma.ID,
		&firma.UsuarioID,
		&firma.NombreArchivo,
		&firma.RutaArchivo,
		&firma.TipoMime,
		&firma.Formato,
		&firma.DatosFirma,
		&firma.TamanoBytes,
		&firma.EsPredeterminada,
		&firma.CreadoEn,
		&firma.ActualizadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &firma, nil
}

// GetDefaultByUsuarioEmail obtiene la firma predeterminada de un usuario por email
func (r *FirmaRepository) GetDefaultByUsuarioEmail(email string) (*models.FirmaDigital, error) {
	query := `
		SELECT f.id, f.usuario_id, f.nombre_archivo, f.ruta_archivo, f.tipo_mime, f.formato, 
		       f.datos_firma, f.tamano_bytes, f.es_predeterminada, f.creado_en, f.actualizado_en
		FROM firmas_digitales f
		INNER JOIN usuarios u ON f.usuario_id = u.id
		WHERE u.email = $1 AND f.es_predeterminada = true
		LIMIT 1
	`

	var firma models.FirmaDigital
	err := r.db.QueryRow(query, email).Scan(
		&firma.ID,
		&firma.UsuarioID,
		&firma.NombreArchivo,
		&firma.RutaArchivo,
		&firma.TipoMime,
		&firma.Formato,
		&firma.DatosFirma,
		&firma.TamanoBytes,
		&firma.EsPredeterminada,
		&firma.CreadoEn,
		&firma.ActualizadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &firma, nil
}

// Update actualiza una firma
func (r *FirmaRepository) Update(firma *models.FirmaDigital) error {
	query := `
		UPDATE firmas_digitales
		SET nombre_archivo = $1, es_predeterminada = $2, actualizado_en = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING actualizado_en
	`

	err := r.db.QueryRow(query, firma.NombreArchivo, firma.EsPredeterminada, firma.ID).Scan(&firma.ActualizadoEn)
	if err != nil {
		return fmt.Errorf("error al actualizar firma: %v", err)
	}

	// Si se marca como predeterminada, desmarcar otras
	if firma.EsPredeterminada {
		return r.UnsetOtherDefaults(firma.UsuarioID, firma.ID)
	}

	return nil
}

// Delete elimina una firma
func (r *FirmaRepository) Delete(id int) error {
	query := `DELETE FROM firmas_digitales WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar firma: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("firma no encontrada")
	}

	return nil
}

// SetAsDefault marca una firma como predeterminada
func (r *FirmaRepository) SetAsDefault(id int, usuarioID int) error {
	// Primero desmarcar todas las firmas del usuario
	err := r.UnsetOtherDefaults(usuarioID, 0)
	if err != nil {
		return err
	}

	// Luego marcar la firma especificada como predeterminada
	query := `
		UPDATE firmas_digitales
		SET es_predeterminada = true, actualizado_en = CURRENT_TIMESTAMP
		WHERE id = $1 AND usuario_id = $2
	`

	result, err := r.db.Exec(query, id, usuarioID)
	if err != nil {
		return fmt.Errorf("error al establecer firma predeterminada: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("firma no encontrada o no pertenece al usuario")
	}

	return nil
}

// UnsetOtherDefaults desmarca todas las firmas predeterminadas excepto la especificada
func (r *FirmaRepository) UnsetOtherDefaults(usuarioID int, exceptID int) error {
	query := `
		UPDATE firmas_digitales
		SET es_predeterminada = false, actualizado_en = CURRENT_TIMESTAMP
		WHERE usuario_id = $1 AND id != $2 AND es_predeterminada = true
	`

	_, err := r.db.Exec(query, usuarioID, exceptID)
	return err
}
