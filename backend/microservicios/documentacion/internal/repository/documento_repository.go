package repository

import (
	"database/sql"
	"documentacion/internal/models"
	"fmt"
	"strings"
	"time"
)

type DocumentoRepository struct {
	db *sql.DB
}

func NewDocumentoRepository(db *sql.DB) *DocumentoRepository {
	return &DocumentoRepository{db: db}
}

// Create crea un nuevo documento en la base de datos
func (r *DocumentoRepository) Create(req *models.CreateDocumentoRequest, rutaArchivo string, tamanoBytes int64, tipoArchivo string) (*models.Documento, error) {
	query := `
		INSERT INTO documentos (
			activo_id, tecnico_id, nombre, descripcion, categoria, 
			tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision, 
			subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING id, creado_en, actualizado_en`

	var doc models.Documento
	err := r.db.QueryRow(
		query,
		req.ActivoID, req.TecnicoID, req.Nombre, req.Descripcion, req.Categoria,
		tipoArchivo, rutaArchivo, tamanoBytes, req.FechaEmision,
		req.SubidoPor, req.PalabrasClave, req.EsFichaTecnica,
	).Scan(&doc.ID, &doc.CreadoEn, &doc.ActualizadoEn)

	if err != nil {
		return nil, err
	}

	// Llenar el resto de los campos
	doc.ActivoID = req.ActivoID
	doc.TecnicoID = req.TecnicoID
	doc.Nombre = req.Nombre
	doc.Descripcion = req.Descripcion
	doc.Categoria = req.Categoria
	doc.TipoArchivo = tipoArchivo
	doc.RutaArchivo = rutaArchivo
	doc.TamanoBytes = tamanoBytes
	doc.FechaEmision = req.FechaEmision
	doc.SubidoPor = req.SubidoPor
	doc.PalabrasClave = req.PalabrasClave
	doc.EsFichaTecnica = req.EsFichaTecnica

	return &doc, nil
}

// GetByID obtiene un documento por su ID
func (r *DocumentoRepository) GetByID(id int) (*models.Documento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, nombre, descripcion, categoria,
			   tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
			   subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		FROM documentos WHERE id = $1`

	var doc models.Documento
	err := r.db.QueryRow(query, id).Scan(
		&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
		&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
		&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
		&doc.CreadoEn, &doc.ActualizadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// GetAll obtiene todos los documentos con filtros opcionales
func (r *DocumentoRepository) GetAll(activoID *int, soloFichasTecnicas bool) ([]models.Documento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, nombre, descripcion, categoria,
			   tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
			   subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		FROM documentos WHERE 1=1`

	var args []interface{}
	argCount := 0

	// Filtro por activo_id si se proporciona
	if activoID != nil {
		argCount++
		query += fmt.Sprintf(" AND activo_id = $%d", argCount)
		args = append(args, *activoID)
	}

	// Filtro por fichas técnicas si se solicita
	if soloFichasTecnicas {
		query += " AND es_ficha_tecnica = true"
	}

	query += " ORDER BY creado_en DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var doc models.Documento
		err := rows.Scan(
			&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
			&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
			&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
			&doc.CreadoEn, &doc.ActualizadoEn,
		)
		if err != nil {
			return nil, err
		}
		documentos = append(documentos, doc)
	}

	return documentos, nil
}

// GetByActivoID obtiene todos los documentos de un activo
func (r *DocumentoRepository) GetByActivoID(activoID int) ([]models.Documento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, nombre, descripcion, categoria,
			   tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
			   subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		FROM documentos 
		WHERE activo_id = $1
		ORDER BY creado_en DESC`

	rows, err := r.db.Query(query, activoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var doc models.Documento
		err := rows.Scan(
			&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
			&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
			&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
			&doc.CreadoEn, &doc.ActualizadoEn,
		)
		if err != nil {
			return nil, err
		}
		documentos = append(documentos, doc)
	}

	return documentos, nil
}

// GetFichaTecnicaPorActivo obtiene la ficha técnica de un activo (HdU23)
func (r *DocumentoRepository) GetFichaTecnicaPorActivo(activoID int) (*models.Documento, error) {
	query := `
		SELECT id, activo_id, tecnico_id, nombre, descripcion, categoria,
			   tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
			   subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		FROM documentos 
		WHERE activo_id = $1 AND es_ficha_tecnica = true
		ORDER BY creado_en DESC
		LIMIT 1`

	var doc models.Documento
	err := r.db.QueryRow(query, activoID).Scan(
		&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
		&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
		&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
		&doc.CreadoEn, &doc.ActualizadoEn,
	)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// BuscarConFiltros busca documentos con filtros múltiples (HdU05)
func (r *DocumentoRepository) BuscarConFiltros(filtros models.DocumentoFiltros) ([]models.Documento, error) {
	baseQuery := `
		SELECT id, activo_id, tecnico_id, nombre, descripcion, categoria,
			   tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
			   subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
		FROM documentos WHERE 1=1`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Filtro por activo
	if filtros.ActivoID != nil {
		conditions = append(conditions, fmt.Sprintf("activo_id = $%d", argIndex))
		args = append(args, *filtros.ActivoID)
		argIndex++
	}

	// Filtro por categoría
	if filtros.Categoria != "" {
		conditions = append(conditions, fmt.Sprintf("categoria = $%d", argIndex))
		args = append(args, filtros.Categoria)
		argIndex++
	}

	// Filtro por fecha desde
	if filtros.FechaDesde != nil {
		conditions = append(conditions, fmt.Sprintf("fecha_emision >= $%d", argIndex))
		args = append(args, *filtros.FechaDesde)
		argIndex++
	}

	// Filtro por fecha hasta
	if filtros.FechaHasta != nil {
		conditions = append(conditions, fmt.Sprintf("fecha_emision <= $%d", argIndex))
		args = append(args, *filtros.FechaHasta)
		argIndex++
	}

	// Filtro por palabra clave
	if filtros.PalabraClave != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(nombre ILIKE $%d OR descripcion ILIKE $%d OR palabras_clave ILIKE $%d)",
			argIndex, argIndex, argIndex))
		args = append(args, "%"+filtros.PalabraClave+"%")
		argIndex++
	}

	// Filtro por ficha técnica
	if filtros.EsFichaTecnica != nil {
		conditions = append(conditions, fmt.Sprintf("es_ficha_tecnica = $%d", argIndex))
		args = append(args, *filtros.EsFichaTecnica)
		argIndex++
	}

	// Construir query final
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY creado_en DESC"

	// Paginación
	if filtros.Limite > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filtros.Limite)
		argIndex++

		if filtros.Offset > 0 {
			baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filtros.Offset)
		}
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documentos []models.Documento
	for rows.Next() {
		var doc models.Documento
		err := rows.Scan(
			&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
			&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
			&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
			&doc.CreadoEn, &doc.ActualizadoEn,
		)
		if err != nil {
			return nil, err
		}
		documentos = append(documentos, doc)
	}

	return documentos, nil
}

// GetHistorialMantenimiento construye el historial de mantenimiento de un activo
func (r *DocumentoRepository) GetHistorialMantenimiento(activoID int) (*models.HistorialMantenimiento, error) {
	// Obtener documentos del activo
	documentos, err := r.GetByActivoID(activoID)
	if err != nil {
		return nil, err
	}

	// Buscar la fecha del último mantenimiento
	var ultimoMantenimiento *time.Time
	for _, doc := range documentos {
		if doc.Categoria == models.CategoriaInformeMantenimiento {
			if ultimoMantenimiento == nil || doc.FechaEmision.After(*ultimoMantenimiento) {
				ultimoMantenimiento = &doc.FechaEmision
			}
		}
	}

	historial := &models.HistorialMantenimiento{
		ActivoID:            activoID,
		DocumentosTotal:     len(documentos),
		UltimoMantenimiento: ultimoMantenimiento,
		Documentos:          documentos,
	}

	return historial, nil
}

// Update actualiza un documento existente
func (r *DocumentoRepository) Update(id int, req *models.UpdateDocumentoRequest) (*models.Documento, error) {
	// Construir query dinámicamente basado en los campos presentes
	setParts := []string{}
	args := []interface{}{id} // ID siempre es el primer argumento
	argCount := 2

	if req.Nombre != nil {
		setParts = append(setParts, fmt.Sprintf("nombre = $%d", argCount))
		args = append(args, *req.Nombre)
		argCount++
	}

	if req.Descripcion != nil {
		setParts = append(setParts, fmt.Sprintf("descripcion = $%d", argCount))
		args = append(args, *req.Descripcion)
		argCount++
	}

	if req.Categoria != nil {
		setParts = append(setParts, fmt.Sprintf("categoria = $%d", argCount))
		args = append(args, *req.Categoria)
		argCount++
	}

	if req.PalabrasClave != nil {
		setParts = append(setParts, fmt.Sprintf("palabras_clave = $%d", argCount))
		args = append(args, *req.PalabrasClave)
		argCount++
	}

	if req.EsFichaTecnica != nil {
		setParts = append(setParts, fmt.Sprintf("es_ficha_tecnica = $%d", argCount))
		args = append(args, *req.EsFichaTecnica)
		argCount++
	}

	// Siempre actualizar timestamp
	setParts = append(setParts, "actualizado_en = NOW()")

	if len(setParts) == 1 { // Solo timestamp, no hay campos para actualizar
		return nil, fmt.Errorf("no hay campos para actualizar")
	}

	query := fmt.Sprintf(`
		UPDATE documentos 
		SET %s
		WHERE id = $1
		RETURNING id, activo_id, tecnico_id, nombre, descripcion, categoria,
				  tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision,
				  subido_por, palabras_clave, es_ficha_tecnica, creado_en, actualizado_en
	`, strings.Join(setParts, ", "))

	var doc models.Documento
	err := r.db.QueryRow(query, args...).Scan(
		&doc.ID, &doc.ActivoID, &doc.TecnicoID, &doc.Nombre, &doc.Descripcion,
		&doc.Categoria, &doc.TipoArchivo, &doc.RutaArchivo, &doc.TamanoBytes,
		&doc.FechaEmision, &doc.SubidoPor, &doc.PalabrasClave, &doc.EsFichaTecnica,
		&doc.CreadoEn, &doc.ActualizadoEn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("documento no encontrado")
		}
		return nil, err
	}

	return &doc, nil
}
