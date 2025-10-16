package models

import (
	"time"
)

// Tipos de activos permitidos
const (
	TipoCaldera       = "caldera"
	TipoBombaAgua     = "bomba de agua"
	TipoAscensor      = "ascensor"
	TipoTransformador = "transformador"
)

// TiposActivosValidos contiene todos los tipos de activos permitidos
var TiposActivosValidos = []string{
	TipoCaldera,
	TipoBombaAgua,
	TipoAscensor,
	TipoTransformador,
}

// ValidarTipoActivo verifica si un tipo de activo es válido
func ValidarTipoActivo(tipo string) bool {
	for _, tipoValido := range TiposActivosValidos {
		if tipo == tipoValido {
			return true
		}
	}
	return false
}

// Empresa contratista
type Empresa struct {
	ID       int    `json:"id" db:"id"`
	Nombre   string `json:"nombre" db:"nombre"`
	RUT      string `json:"rut" db:"rut"`
	Telefono string `json:"telefono" db:"telefono"`
	Email    string `json:"email" db:"email"`
}

// Técnico especializado para mantenimiento
type Tecnico struct {
	ID           int       `json:"id" db:"id"`
	Nombre       string    `json:"nombre" db:"nombre"`
	Apellido     string    `json:"apellido" db:"apellido"`
	Email        string    `json:"email" db:"email"`
	Telefono     string    `json:"telefono" db:"telefono"`
	Especialidad string    `json:"especialidad" db:"especialidad"`
	Autorizado   bool      `json:"autorizado" db:"autorizado"`
	EmpresaID    int       `json:"empresa_id" db:"empresa_id"`
	Empresa      Empresa   `json:"empresa"`
	CreadoEn     time.Time `json:"creado_en" db:"creado_en"`
}

// Edificio donde se ubican los activos
type Edificio struct {
	ID        int       `json:"id" db:"id"`
	Nombre    string    `json:"nombre" db:"nombre"`
	Direccion string    `json:"direccion" db:"direccion"`
	CreadoEn  time.Time `json:"creado_en" db:"creado_en"`
}

// Usuario del sistema
type Usuario struct {
	ID       int       `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	CreadoEn time.Time `json:"creado_en" db:"creado_en"`
}

// Relación muchos a muchos entre usuarios y edificios
type UsuarioEdificio struct {
	UsuarioID  int       `json:"usuario_id" db:"usuario_id"`
	EdificioID int       `json:"edificio_id" db:"edificio_id"`
	AsignadoEn time.Time `json:"asignado_en" db:"asignado_en"`
}

// Activo gestionado
type Activo struct {
	ID         int       `json:"id" db:"id"`
	Nombre     string    `json:"nombre" db:"nombre"`
	Tipo       string    `json:"tipo" db:"tipo"`
	Ubicacion  string    `json:"ubicacion" db:"ubicacion"`
	EdificioID int       `json:"edificio_id" db:"edificio_id"`
	CreadoEn   time.Time `json:"creado_en" db:"creado_en"`
}

// Acción de mantenimiento colaborativa
type AccionMantenimiento struct {
	ID          int        `json:"id" db:"id"`
	ActivoID    int        `json:"activo_id" db:"activo_id"`
	TecnicoID   int        `json:"tecnico_id" db:"tecnico_id"`
	Tipo        string     `json:"tipo" db:"tipo"` // preventivo, correctivo, emergencia
	Descripcion string     `json:"descripcion" db:"descripcion"`
	Estado      string     `json:"estado" db:"estado"`       // pendiente, en_progreso, completado
	Prioridad   string     `json:"prioridad" db:"prioridad"` // baja, media, alta, critica
	FechaInicio time.Time  `json:"fecha_inicio" db:"fecha_inicio"`
	FechaFin    *time.Time `json:"fecha_fin,omitempty" db:"fecha_fin"`
	CreadoEn    time.Time  `json:"creado_en" db:"creado_en"`
}

// Relación muchos a muchos entre activos y técnicos
type ActivoTecnico struct {
	ActivoID   int       `json:"activo_id" db:"activo_id"`
	TecnicoID  int       `json:"tecnico_id" db:"tecnico_id"`
	AsignadoEn time.Time `json:"asignado_en" db:"asignado_en"`
}

// Reporte automático
type Reporte struct {
	ID          int    `json:"id" db:"id"`
	ActivoID    int    `json:"activo_id" db:"activo_id"`
	TipoReporte string `json:"tipo_reporte" db:"tipo_reporte"` // semanal, mensual, incidente, mantenimiento
	Contenido   string `json:"contenido" db:"contenido"`
	// Nuevos campos para observaciones y estructura
	ObservacionesAnalista string                 `json:"observaciones_analista" db:"observaciones_analista"`
	AutorAnalista         string                 `json:"autor_analista" db:"autor_analista"`
	EstructuraInforme     map[string]interface{} `json:"estructura_informe" db:"estructura_informe"`
	MetadataInforme       map[string]interface{} `json:"metadata_informe" db:"metadata_informe"`
	VersionReporte        int                    `json:"version_reporte" db:"version_reporte"`
	EstadoRevision        string                 `json:"estado_revision" db:"estado_revision"` // pendiente, en_revision, aprobado, rechazado
	FechaRevision         *time.Time             `json:"fecha_revision" db:"fecha_revision"`
	Revisor               *string                `json:"revisor" db:"revisor"`
	//
	GeneradoEn time.Time `json:"generado_en" db:"generado_en"`
	Estado     string    `json:"estado" db:"estado"` // generado, enviado, archivado
}

// DTOs para responses con relaciones
type EmpresaResponse struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	RUT      string `json:"rut"`
	Telefono string `json:"telefono"`
	Email    string `json:"email"`
}

type TecnicoResponse struct {
	ID             int             `json:"id"`
	Nombre         string          `json:"nombre"`
	Apellido       string          `json:"apellido"`
	NombreCompleto string          `json:"nombre_completo"`
	Email          string          `json:"email"`
	Telefono       string          `json:"telefono"`
	Especialidad   string          `json:"especialidad"`
	Autorizado     bool            `json:"autorizado"`
	Empresa        EmpresaResponse `json:"empresa"`
	CreadoEn       time.Time       `json:"creado_en"`
}

type TecnicoConActivos struct {
	Tecnico `json:",inline"`
	Activos []Activo `json:"activos,omitempty"`
}

type ActivoConEdificio struct {
	Activo   `json:",inline"`
	Edificio Edificio `json:"edificio"`
}

type AccionConDetalles struct {
	AccionMantenimiento `json:",inline"`
	Activo              Activo  `json:"activo"`
	Tecnico             Tecnico `json:"tecnico"`
}

type ReporteCompleto struct {
	Reporte           `json:",inline"`
	ActivoConEdificio ActivoConEdificio    `json:"activo"`
	UltimaAccion      *AccionMantenimiento `json:"ultima_accion,omitempty"`
}

// DTOs para requests
type CreateTecnicoRequest struct {
	Nombre       string `json:"nombre" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Telefono     string `json:"telefono" binding:"required"`
	Especialidad string `json:"especialidad" binding:"required"`
}

type CreateActivoRequest struct {
	Nombre     string `json:"nombre" binding:"required"`
	Tipo       string `json:"tipo" binding:"required"`
	Ubicacion  string `json:"ubicacion" binding:"required"`
	EdificioID int    `json:"edificio_id" binding:"required"`
}

type UpdateAutorizadoTecnicoRequest struct {
	Autorizado *bool `json:"autorizado" binding:"required"`
}

type AsignarTecnicoActivoRequest struct {
	TecnicoID int `json:"tecnico_id" binding:"required"`
}

type CreateAccionMantenimientoRequest struct {
	ActivoID    int    `json:"activo_id" binding:"required"`
	TecnicoID   int    `json:"tecnico_id" binding:"required"`
	Tipo        string `json:"tipo" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	Prioridad   string `json:"prioridad" binding:"required"`
}

// DTOs para reportes con observaciones
type CreateReporteRequest struct {
	ActivoID              int    `json:"activo_id" binding:"required"`
	TipoReporte           string `json:"tipo_reporte" binding:"required"`
	Contenido             string `json:"contenido,omitempty"`
	ObservacionesAnalista string `json:"observaciones_analista,omitempty"`
	AutorAnalista         string `json:"autor_analista,omitempty"`
}

// DTO para generar reporte por activo con campos solicitados
type GenerarReporteRequest struct {
	Campos                  []string `json:"campos" binding:"required"` // e.g. ["ubicacion","historico_mantenimientos","ultima_acciones","datos_sensores"]
	FirmaID                 *int     `json:"firma_id,omitempty"`        // ID de la firma digital a incluir (opcional)
	UsarFirmaPredeterminada bool     `json:"usar_firma_predeterminada"` // Si true, usa la firma predeterminada del usuario
	Email                   string   `json:"email,omitempty"`           // Email del usuario que genera el reporte (para obtener firma predeterminada)
}

type UpdateObservacionesRequest struct {
	ObservacionesAnalista string `json:"observaciones_analista" binding:"required"`
	AutorAnalista         string `json:"autor_analista" binding:"required"`
}

type UpdateEstadoRevisionRequest struct {
	EstadoRevision string `json:"estado_revision" binding:"required"`
	Revisor        string `json:"revisor" binding:"required"`
	Observaciones  string `json:"observaciones,omitempty"`
}

type ReporteConObservaciones struct {
	Reporte
	NombreActivo string `json:"nombre_activo"`
	TipoActivo   string `json:"tipo_activo"`
}

type EstructuraInformeReporte struct {
	Resumen            string                 `json:"resumen"`
	Observaciones      string                 `json:"observaciones"`
	DatosActivo        map[string]interface{} `json:"datos_activo"`
	AccionesRealizadas []string               `json:"acciones_realizadas"`
	Recomendaciones    []string               `json:"recomendaciones"`
	Conclusiones       string                 `json:"conclusiones"`
}

// Modelos para reportes de fallas de usuarios

type TipoFalla struct {
	IDFalla          int       `json:"id_falla" db:"id_falla"`
	Tipo             string    `json:"tipo" db:"tipo"`
	Descripcion      string    `json:"descripcion" db:"descripcion"`
	FechaPublicacion time.Time `json:"fecha_publicacion" db:"fecha_publicacion"`
	IDUsuario        int       `json:"id_usuario" db:"id_usuario"`
	IDEdificio       int       `json:"id_edificio" db:"id_edificio"`
	Estado           string    `json:"estado" db:"estado"`
}

type Comentario struct {
	IDComentario    int       `json:"id_comentario" db:"id_comentario"`
	IDFalla         int       `json:"id_falla" db:"id_falla"`
	IDUsuario       int       `json:"id_usuario" db:"id_usuario"`
	Comentario      string    `json:"comentario" db:"comentario"`
	FechaComentario time.Time `json:"fecha_comentario" db:"fecha_comentario"`
}

// DTOs para las APIs

type CreateTipoFallaRequest struct {
	Email       string `json:"email" binding:"required"`
	Tipo        string `json:"tipo" binding:"required"`
	Descripcion string `json:"descripcion"`
	IDEdificio  int    `json:"id_edificio" binding:"required"`
}

type CreateComentarioRequest struct {
	Email      string `json:"email" binding:"required"`
	IDFalla    int    `json:"id_falla" binding:"required"`
	Comentario string `json:"comentario" binding:"required"`
}

// DTO para respuesta con username en lugar de id_usuario
type TipoFallaResponse struct {
	IDFalla          int                  `json:"id_falla"`
	Tipo             string               `json:"tipo"`
	Descripcion      string               `json:"descripcion"`
	FechaPublicacion time.Time            `json:"fecha_publicacion"`
	Username         string               `json:"username"`
	IDEdificio       int                  `json:"id_edificio"`
	Estado           string               `json:"estado"`
	Comentarios      []ComentarioResponse `json:"comentarios"`
}

type ComentarioResponse struct {
	IDComentario    int       `json:"id_comentario"`
	Username        string    `json:"username"`
	Comentario      string    `json:"comentario"`
	FechaComentario time.Time `json:"fecha_comentario"`
}

type UpdateEstadoAccionRequest struct {
	Estado string `json:"estado" binding:"required"`
}

// Firma Digital para usuarios
type FirmaDigital struct {
	ID               int       `json:"id" db:"id"`
	UsuarioID        int       `json:"usuario_id" db:"usuario_id"`
	NombreArchivo    string    `json:"nombre_archivo" db:"nombre_archivo"`
	RutaArchivo      string    `json:"ruta_archivo" db:"ruta_archivo"`
	TipoMime         string    `json:"tipo_mime" db:"tipo_mime"`               // image/png, image/jpeg, image/svg+xml
	Formato          string    `json:"formato" db:"formato"`                   // png, jpeg, jpg, svg
	DatosFirma       []byte    `json:"datos_firma,omitempty" db:"datos_firma"` // Para SVG o datos binarios
	TamanoBytes      int64     `json:"tamano_bytes" db:"tamano_bytes"`
	EsPredeterminada bool      `json:"es_predeterminada" db:"es_predeterminada"`
	CreadoEn         time.Time `json:"creado_en" db:"creado_en"`
	ActualizadoEn    time.Time `json:"actualizado_en" db:"actualizado_en"`
}

// DTOs para firmas digitales
type CreateFirmaRequest struct {
	Email            string `form:"email" binding:"required,email"`
	NombreArchivo    string `form:"nombre_archivo"`
	EsPredeterminada bool   `form:"es_predeterminada"`
	// El archivo se maneja por separado en el handler mediante c.FormFile()
}

type CreateFirmaSVGRequest struct {
	Email            string `json:"email" binding:"required,email"`
	NombreArchivo    string `json:"nombre_archivo" binding:"required"`
	DatosSVG         string `json:"datos_svg" binding:"required"` // SVG como string
	EsPredeterminada bool   `json:"es_predeterminada"`
}

type UpdateFirmaRequest struct {
	NombreArchivo    string `json:"nombre_archivo"`
	EsPredeterminada bool   `json:"es_predeterminada"`
}

type FirmaResponse struct {
	ID               int       `json:"id"`
	UsuarioID        int       `json:"usuario_id"`
	NombreArchivo    string    `json:"nombre_archivo"`
	RutaArchivo      string    `json:"ruta_archivo"`
	TipoMime         string    `json:"tipo_mime"`
	Formato          string    `json:"formato"`
	TamanoBytes      int64     `json:"tamano_bytes"`
	EsPredeterminada bool      `json:"es_predeterminada"`
	CreadoEn         time.Time `json:"creado_en"`
	ActualizadoEn    time.Time `json:"actualizado_en"`
}
