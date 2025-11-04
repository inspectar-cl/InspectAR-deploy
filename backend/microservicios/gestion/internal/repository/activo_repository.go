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
		SELECT id, nombre, tipo, descripcion, ubicacion, edificio_id, 
		       codigo_activo, codigo_qr, url_qr, qr_generado_en, secuencial, creado_en 
		FROM activos 
		WHERE id = $1
	`

	var activo models.Activo
	err := r.db.QueryRow(query, id).Scan(
		&activo.ID,
		&activo.Nombre,
		&activo.Tipo,
		&activo.Descripcion,
		&activo.Ubicacion,
		&activo.EdificioID,
		&activo.CodigoActivo,
		&activo.CodigoQR,
		&activo.URLQR,
		&activo.QRGeneradoEn,
		&activo.Secuencial,
		&activo.CreadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &activo, nil
}

func (r *ActivoRepository) GetAll() ([]models.Activo, error) {
	query := `
		SELECT id, nombre, tipo, descripcion, ubicacion, edificio_id, 
		       codigo_activo, codigo_qr, url_qr, qr_generado_en, secuencial, creado_en 
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
			&activo.Nombre,
			&activo.Tipo,
			&activo.Descripcion,
			&activo.Ubicacion,
			&activo.EdificioID,
			&activo.CodigoActivo,
			&activo.CodigoQR,
			&activo.URLQR,
			&activo.QRGeneradoEn,
			&activo.Secuencial,
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
		SELECT id, nombre, tipo, descripcion, ubicacion, edificio_id, 
		       codigo_activo, codigo_qr, url_qr, qr_generado_en, secuencial, creado_en 
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
			&activo.Nombre,
			&activo.Tipo,
			&activo.Descripcion,
			&activo.Ubicacion,
			&activo.EdificioID,
			&activo.CodigoActivo,
			&activo.CodigoQR,
			&activo.URLQR,
			&activo.QRGeneradoEn,
			&activo.Secuencial,
			&activo.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		activos = append(activos, activo)
	}

	return activos, nil
}

func (r *ActivoRepository) GetByTipo(tipo string) ([]models.Activo, error) {
	query := `
		SELECT id, nombre, tipo, descripcion, ubicacion, edificio_id, 
		       codigo_activo, codigo_qr, url_qr, qr_generado_en, secuencial, creado_en 
		FROM activos 
		WHERE tipo = $1
		ORDER BY nombre
	`

	rows, err := r.db.Query(query, tipo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activos []models.Activo
	for rows.Next() {
		var activo models.Activo
		err := rows.Scan(
			&activo.ID,
			&activo.Nombre,
			&activo.Tipo,
			&activo.Descripcion,
			&activo.Ubicacion,
			&activo.EdificioID,
			&activo.CodigoActivo,
			&activo.CodigoQR,
			&activo.URLQR,
			&activo.QRGeneradoEn,
			&activo.Secuencial,
			&activo.CreadoEn,
		)
		if err != nil {
			return nil, err
		}
		activos = append(activos, activo)
	}

	return activos, nil
}

// Crear un activo (admin) - El código se genera automáticamente por trigger
func (r *ActivoRepository) Crear(activo models.Activo) (*models.Activo, error) {
	query := `
		INSERT INTO activos (nombre, tipo, descripcion, ubicacion, edificio_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, codigo_activo, secuencial, creado_en
	`

	var codigoActivo string
	var secuencial int
	err := r.db.QueryRow(
		query,
		activo.Nombre,
		activo.Tipo,
		activo.Descripcion,
		activo.Ubicacion,
		activo.EdificioID,
	).Scan(&activo.ID, &codigoActivo, &secuencial, &activo.CreadoEn)

	if err != nil {
		return nil, err
	}

	// Asignar los valores generados
	activo.CodigoActivo = &codigoActivo
	activo.Secuencial = &secuencial

	return &activo, nil
}

// Actualizar un activo (admin)
func (r *ActivoRepository) Actualizar(id int, activo models.Activo) error {
	query := `
		UPDATE activos 
		SET nombre = $1, tipo = $2, descripcion = $3, ubicacion = $4, edificio_id = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		activo.Nombre,
		activo.Tipo,
		activo.Descripcion,
		activo.Ubicacion,
		activo.EdificioID,
		id,
	)

	return err
}

// Eliminar un activo (admin)
func (r *ActivoRepository) Eliminar(id int) error {
	query := `DELETE FROM activos WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ActualizarQR actualiza el código QR y URL del activo
func (r *ActivoRepository) ActualizarQR(activoID int, codigoQR, urlQR string) error {
	query := `
		UPDATE activos 
		SET codigo_qr = $1, 
		    url_qr = $2, 
		    qr_generado_en = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, codigoQR, urlQR, activoID)
	return err
}

// GetByCodigo obtiene un activo por su código único
func (r *ActivoRepository) GetByCodigo(codigo string) (*models.Activo, error) {
	query := `
		SELECT id, nombre, tipo, descripcion, ubicacion, edificio_id, 
		       codigo_activo, codigo_qr, url_qr, qr_generado_en, secuencial, creado_en
		FROM activos 
		WHERE codigo_activo = $1
	`

	var activo models.Activo
	err := r.db.QueryRow(query, codigo).Scan(
		&activo.ID,
		&activo.Nombre,
		&activo.Tipo,
		&activo.Descripcion,
		&activo.Ubicacion,
		&activo.EdificioID,
		&activo.CodigoActivo,
		&activo.CodigoQR,
		&activo.URLQR,
		&activo.QRGeneradoEn,
		&activo.Secuencial,
		&activo.CreadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &activo, nil
}
