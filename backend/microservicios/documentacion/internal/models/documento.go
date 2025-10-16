package models

import "time"

// Documento representa un documento técnico subido al sistema
type Documento struct {
	ID             int       `json:"id" db:"id"`
	ActivoID       int       `json:"activo_id" db:"activo_id"`
	TecnicoID      *int      `json:"tecnico_id,omitempty" db:"tecnico_id"` // Opcional, para HdU19
	Nombre         string    `json:"nombre" db:"nombre"`
	Descripcion    string    `json:"descripcion" db:"descripcion"`
	Categoria      string    `json:"categoria" db:"categoria"`       // ficha_tecnica, informe_mantenimiento, diagnostico, manual_fabricante
	TipoArchivo    string    `json:"tipo_archivo" db:"tipo_archivo"` // pdf, docx
	RutaArchivo    string    `json:"ruta_archivo" db:"ruta_archivo"`
	TamanoBytes    int64     `json:"tamano_bytes" db:"tamano_bytes"`
	FechaEmision   time.Time `json:"fecha_emision" db:"fecha_emision"`
	SubidoPor      string    `json:"subido_por" db:"subido_por"`             // usuario que subió el documento
	PalabrasClave  string    `json:"palabras_clave" db:"palabras_clave"`     // separadas por comas
	EsFichaTecnica bool      `json:"es_ficha_tecnica" db:"es_ficha_tecnica"` // Para HdU23
	CreadoEn       time.Time `json:"creado_en" db:"creado_en"`
	ActualizadoEn  time.Time `json:"actualizado_en" db:"actualizado_en"`
}

// CreateDocumentoRequest para subir un nuevo documento
type CreateDocumentoRequest struct {
	ActivoID       int       `json:"activo_id" binding:"required"`
	TecnicoID      *int      `json:"tecnico_id,omitempty"`
	Nombre         string    `json:"nombre" binding:"required"`
	Descripcion    string    `json:"descripcion"`
	Categoria      string    `json:"categoria" binding:"required,oneof=ficha_tecnica informe_mantenimiento diagnostico manual_fabricante certificacion"`
	FechaEmision   time.Time `json:"fecha_emision" binding:"required"`
	SubidoPor      string    `json:"subido_por" binding:"required"`
	PalabrasClave  string    `json:"palabras_clave"`
	EsFichaTecnica bool      `json:"es_ficha_tecnica"`
}

// UpdateDocumentoRequest para actualizar un documento existente
type UpdateDocumentoRequest struct {
	Nombre         *string `json:"nombre,omitempty"`
	Descripcion    *string `json:"descripcion,omitempty"`
	Categoria      *string `json:"categoria,omitempty" binding:"omitempty,oneof=ficha_tecnica informe_mantenimiento diagnostico manual_fabricante certificacion"`
	PalabrasClave  *string `json:"palabras_clave,omitempty"`
	EsFichaTecnica *bool   `json:"es_ficha_tecnica,omitempty"`
}

// DocumentoFiltros para búsqueda y filtrado
type DocumentoFiltros struct {
	ActivoID       *int       `json:"activo_id,omitempty"`
	Categoria      string     `json:"categoria,omitempty"`
	FechaDesde     *time.Time `json:"fecha_desde,omitempty"`
	FechaHasta     *time.Time `json:"fecha_hasta,omitempty"`
	PalabraClave   string     `json:"palabra_clave,omitempty"`
	EsFichaTecnica *bool      `json:"es_ficha_tecnica,omitempty"`
	Limite         int        `json:"limite,omitempty"`
	Offset         int        `json:"offset,omitempty"`
}

// AnalisisIA representa el análisis de IA de un documento (HdU19)
type AnalisisIA struct {
	ID           int       `json:"id" db:"id"`
	DocumentoID  int       `json:"documento_id" db:"documento_id"`
	Resumen      string    `json:"resumen" db:"resumen"`
	PuntosClaves string    `json:"puntos_claves" db:"puntos_claves"` // JSON array de puntos importantes
	Graficos     string    `json:"graficos" db:"graficos"`           // JSON array de descripciones de gráficos
	Estado       string    `json:"estado" db:"estado"`               // procesando, completado, error
	CreadoEn     time.Time `json:"creado_en" db:"creado_en"`
}

// CreateAnalisisRequest para solicitar análisis de IA
type CreateAnalisisRequest struct {
	DocumentoID int `json:"documento_id" binding:"required"`
}

// HistorialMantenimiento combina documentos con análisis para construir historial
type HistorialMantenimiento struct {
	ActivoID            int          `json:"activo_id"`
	DocumentosTotal     int          `json:"documentos_total"`
	UltimoMantenimiento *time.Time   `json:"ultimo_mantenimiento,omitempty"`
	Documentos          []Documento  `json:"documentos"`
	AnalisisIA          []AnalisisIA `json:"analisis_ia,omitempty"`
}

// Categorías válidas para documentos
const (
	CategoriaFichaTecnica         = "ficha_tecnica"
	CategoriaInformeMantenimiento = "informe_mantenimiento"
	CategoriaDiagnostico          = "diagnostico"
	CategoriaManualFabricante     = "manual_fabricante"
	CategoriaCertificacion        = "certificacion"
)

// Estados de análisis de IA
const (
	EstadoProcesando = "procesando"
	EstadoCompletado = "completado"
	EstadoError      = "error"
)
