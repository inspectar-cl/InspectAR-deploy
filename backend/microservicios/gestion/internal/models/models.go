package models

import (
	"time"
)

// Técnico especializado para mantenimiento
type Tecnico struct {
	ID           int       `json:"id" db:"id"`
	Nombre       string    `json:"nombre" db:"nombre"`
	Email        string    `json:"email" db:"email"`
	Telefono     string    `json:"telefono" db:"telefono"`
	Especialidad string    `json:"especialidad" db:"especialidad"`
	Disponible   bool      `json:"disponible" db:"disponible"`
	CreadoEn     time.Time `json:"creado_en" db:"creado_en"`
}

// Activo gestionado
type Activo struct {
	ID          int       `json:"id" db:"id"`
	ActivoID    string    `json:"activo_id" db:"activo_id"`
	Nombre      string    `json:"nombre" db:"nombre"`
	Tipo        string    `json:"tipo" db:"tipo"`
	Estado      string    `json:"estado" db:"estado"`
	Ubicacion   string    `json:"ubicacion" db:"ubicacion"`
	EdificioID  string    `json:"edificio_id" db:"edificio_id"`
	CreadoEn    time.Time `json:"creado_en" db:"creado_en"`
}

// Acción de mantenimiento colaborativa
type AccionMantenimiento struct {
	ID          int       `json:"id" db:"id"`
	ActivoID    int       `json:"activo_id" db:"activo_id"`
	TecnicoID   int       `json:"tecnico_id" db:"tecnico_id"`
	Tipo        string    `json:"tipo" db:"tipo"` // preventivo, correctivo, emergencia
	Descripcion string    `json:"descripcion" db:"descripcion"`
	Estado      string    `json:"estado" db:"estado"` // pendiente, en_progreso, completado
	Prioridad   string    `json:"prioridad" db:"prioridad"` // baja, media, alta, critica
	FechaInicio time.Time `json:"fecha_inicio" db:"fecha_inicio"`
	FechaFin    *time.Time `json:"fecha_fin,omitempty" db:"fecha_fin"`
	CreadoEn    time.Time `json:"creado_en" db:"creado_en"`
}

// Reporte automático
type Reporte struct {
	ID           int       `json:"id" db:"id"`
	ActivoID     int       `json:"activo_id" db:"activo_id"`
	TipoReporte  string    `json:"tipo_reporte" db:"tipo_reporte"` // semanal, mensual, incidente
	Contenido    string    `json:"contenido" db:"contenido"`
	GeneradoEn   time.Time `json:"generado_en" db:"generado_en"`
	Estado       string    `json:"estado" db:"estado"` // generado, enviado, archivado
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
	EdificioID string `json:"edificio_id" binding:"required"`
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
