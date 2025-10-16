package repository

import (
	"database/sql"
	"documentacion/internal/models"
)

type AnalisisRepository struct {
	db *sql.DB
}

func NewAnalisisRepository(db *sql.DB) *AnalisisRepository {
	return &AnalisisRepository{db: db}
}

// Create crea un nuevo análisis de IA
func (r *AnalisisRepository) Create(req *models.CreateAnalisisRequest) (*models.AnalisisIA, error) {
	query := `
		INSERT INTO analisis_ia (documento_id, estado, creado_en)
		VALUES ($1, $2, NOW())
		RETURNING id, creado_en`

	var analisis models.AnalisisIA
	err := r.db.QueryRow(query, req.DocumentoID, models.EstadoProcesando).Scan(
		&analisis.ID, &analisis.CreadoEn)

	if err != nil {
		return nil, err
	}

	analisis.DocumentoID = req.DocumentoID
	analisis.Estado = models.EstadoProcesando

	return &analisis, nil
}

// Update actualiza un análisis con los resultados
func (r *AnalisisRepository) Update(id int, resumen, puntosClaves, graficos, estado string) error {
	query := `
		UPDATE analisis_ia 
		SET resumen = $2, puntos_claves = $3, graficos = $4, estado = $5
		WHERE id = $1`

	_, err := r.db.Exec(query, id, resumen, puntosClaves, graficos, estado)
	return err
}

// GetByDocumentoID obtiene análisis por documento
func (r *AnalisisRepository) GetByDocumentoID(documentoID int) (*models.AnalisisIA, error) {
	query := `
		SELECT id, documento_id, resumen, puntos_claves, graficos, estado, creado_en
		FROM analisis_ia 
		WHERE documento_id = $1
		ORDER BY creado_en DESC
		LIMIT 1`

	var analisis models.AnalisisIA
	err := r.db.QueryRow(query, documentoID).Scan(
		&analisis.ID, &analisis.DocumentoID, &analisis.Resumen,
		&analisis.PuntosClaves, &analisis.Graficos, &analisis.Estado,
		&analisis.CreadoEn)

	if err != nil {
		return nil, err
	}

	return &analisis, nil
}

// GetByActivoID obtiene todos los análisis de documentos de un activo
func (r *AnalisisRepository) GetByActivoID(activoID int) ([]models.AnalisisIA, error) {
	query := `
		SELECT a.id, a.documento_id, a.resumen, a.puntos_claves, a.graficos, a.estado, a.creado_en
		FROM analisis_ia a
		INNER JOIN documentos d ON a.documento_id = d.id
		WHERE d.activo_id = $1
		ORDER BY a.creado_en DESC`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analisisList []models.AnalisisIA
	for rows.Next() {
		var analisis models.AnalisisIA
		err := rows.Scan(
			&analisis.ID, &analisis.DocumentoID, &analisis.Resumen,
			&analisis.PuntosClaves, &analisis.Graficos, &analisis.Estado,
			&analisis.CreadoEn)
		if err != nil {
			return nil, err
		}
		analisisList = append(analisisList, analisis)
	}

	return analisisList, nil
}
