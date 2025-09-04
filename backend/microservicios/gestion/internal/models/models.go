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

// Activo gestionado
type Activo struct {
	ID         int       `json:"id" db:"id"`
	ActivoID   string    `json:"activo_id" db:"activo_id"`
	Nombre     string    `json:"nombre" db:"nombre"`
	Tipo       string    `json:"tipo" db:"tipo"`
	Estado     string    `json:"estado" db:"estado"`
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
	ID          int       `json:"id" db:"id"`
	ActivoID    int       `json:"activo_id" db:"activo_id"`
	TipoReporte string    `json:"tipo_reporte" db:"tipo_reporte"` // semanal, mensual, incidente, mantenimiento
	Contenido   string    `json:"contenido" db:"contenido"`
	GeneradoEn  time.Time `json:"generado_en" db:"generado_en"`
	Estado      string    `json:"estado" db:"estado"` // generado, enviado, archivado
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
	ActivoID   string `json:"activo_id" binding:"required"`
	Nombre     string `json:"nombre" binding:"required"`
	Tipo       string `json:"tipo" binding:"required"`
	Ubicacion  string `json:"ubicacion" binding:"required"`
	EdificioID int    `json:"edificio_id" binding:"required"`
}

type UpdateAutorizadoTecnicoRequest struct {
	Autorizado bool `json:"autorizado" binding:"required"`
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

type UpdateEstadoAccionRequest struct {
	Estado string `json:"estado" binding:"required"`
}
