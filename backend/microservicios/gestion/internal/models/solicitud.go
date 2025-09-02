package models

import (
	"time"
)

// TipoSolicitud define los tipos de solicitudes que se pueden enviar
type TipoSolicitud string

const (
	SolicitudMantenimiento TipoSolicitud = "mantenimiento"
	SolicitudReparacion    TipoSolicitud = "reparacion"
	SolicitudInspeccion    TipoSolicitud = "inspeccion"
	SolicitudEmergencia    TipoSolicitud = "emergencia"
	SolicitudConsulta      TipoSolicitud = "consulta"
)

// EstadoSolicitud define los estados posibles de una solicitud
type EstadoSolicitud string

const (
	EstadoPendiente  EstadoSolicitud = "pendiente"
	EstadoEnviada    EstadoSolicitud = "enviada"
	EstadoRecibida   EstadoSolicitud = "recibida"
	EstadoEnProceso  EstadoSolicitud = "en_proceso"
	EstadoCompletada EstadoSolicitud = "completada"
	EstadoCancelada  EstadoSolicitud = "cancelada"
	EstadoRechazada  EstadoSolicitud = "rechazada"
)

// PrioridadSolicitud define la prioridad de la solicitud
type PrioridadSolicitud string

const (
	PrioridadBaja       PrioridadSolicitud = "baja"
	PrioridadMedia      PrioridadSolicitud = "media"
	PrioridadAlta       PrioridadSolicitud = "alta"
	PrioridadCritica    PrioridadSolicitud = "critica"
	PrioridadEmergencia PrioridadSolicitud = "emergencia"
)

// SolicitudTecnico representa una solicitud enviada a un técnico
type SolicitudTecnico struct {
	ID          int `json:"id" gorm:"primaryKey;autoIncrement"`
	TecnicoID   int `json:"tecnico_id" gorm:"not null" validate:"required"`
	ResidenteID int `json:"residente_id" gorm:"not null" validate:"required"`
	ActivoID    int `json:"activo_id" gorm:"not null" validate:"required"`
	EdificioID  int `json:"edificio_id" gorm:"not null" validate:"required"`

	// Contenido de la solicitud
	Tipo        TipoSolicitud      `json:"tipo" gorm:"not null" validate:"required"`
	Asunto      string             `json:"asunto" gorm:"not null;size:200" validate:"required,min=5,max=200"`
	Descripcion string             `json:"descripcion" gorm:"type:text" validate:"required,min=10"`
	Prioridad   PrioridadSolicitud `json:"prioridad" gorm:"not null;default:media"`

	// Estado y seguimiento
	Estado             EstadoSolicitud `json:"estado" gorm:"not null;default:pendiente"`
	FechaCreacion      time.Time       `json:"fecha_creacion" gorm:"autoCreateTime"`
	FechaEnvio         *time.Time      `json:"fecha_envio,omitempty"`
	FechaRecepcion     *time.Time      `json:"fecha_recepcion,omitempty"`
	FechaCompletado    *time.Time      `json:"fecha_completado,omitempty"`
	FechaActualizacion time.Time       `json:"fecha_actualizacion" gorm:"autoUpdateTime"`

	// Información de contacto y respuesta
	MedioContacto    string `json:"medio_contacto" gorm:"size:50"` // email, telefono, sms
	TelefonoContacto string `json:"telefono_contacto" gorm:"size:20"`
	EmailContacto    string `json:"email_contacto" gorm:"size:150"`

	// Respuesta del técnico
	RespuestaTecnico string `json:"respuesta_tecnico,omitempty" gorm:"type:text"`
	NotasInternas    string `json:"notas_internas,omitempty" gorm:"type:text"`

	// Relaciones
	Tecnico Tecnico `json:"tecnico" gorm:"foreignKey:TecnicoID"`

	// Archivos adjuntos (opcional)
	ArchivosAdjuntos []ArchivoSolicitud `json:"archivos_adjuntos,omitempty" gorm:"foreignKey:SolicitudID"`
}

// ArchivoSolicitud representa archivos adjuntos a una solicitud
type ArchivoSolicitud struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SolicitudID   int       `json:"solicitud_id" gorm:"not null"`
	NombreArchivo string    `json:"nombre_archivo" gorm:"not null;size:255"`
	RutaArchivo   string    `json:"ruta_archivo" gorm:"not null;size:500"`
	TipoArchivo   string    `json:"tipo_archivo" gorm:"size:50"`
	TamanoBytes   int64     `json:"tamano_bytes"`
	FechaSubida   time.Time `json:"fecha_subida" gorm:"autoCreateTime"`
}

// CreateSolicitudRequest representa la petición para crear una solicitud
type CreateSolicitudRequest struct {
	TecnicoID        int                `json:"tecnico_id" validate:"required"`
	ResidenteID      int                `json:"residente_id" validate:"required"`
	ActivoID         int                `json:"activo_id" validate:"required"`
	EdificioID       int                `json:"edificio_id" validate:"required"`
	Tipo             TipoSolicitud      `json:"tipo" validate:"required"`
	Asunto           string             `json:"asunto" validate:"required,min=5,max=200"`
	Descripcion      string             `json:"descripcion" validate:"required,min=10"`
	Prioridad        PrioridadSolicitud `json:"prioridad"`
	MedioContacto    string             `json:"medio_contacto" validate:"required"`
	TelefonoContacto string             `json:"telefono_contacto,omitempty"`
	EmailContacto    string             `json:"email_contacto,omitempty" validate:"omitempty,email"`
}

// UpdateSolicitudRequest representa la petición para actualizar una solicitud
type UpdateSolicitudRequest struct {
	Estado           *EstadoSolicitud `json:"estado,omitempty"`
	RespuestaTecnico *string          `json:"respuesta_tecnico,omitempty"`
	NotasInternas    *string          `json:"notas_internas,omitempty"`
	FechaRecepcion   *time.Time       `json:"fecha_recepcion,omitempty"`
	FechaCompletado  *time.Time       `json:"fecha_completado,omitempty"`
}

// SolicitudResponse representa la respuesta simplificada para solicitudes
type SolicitudResponse struct {
	ID               int                `json:"id"`
	TecnicoID        int                `json:"tecnico_id"`
	Tecnico          TecnicoResponse    `json:"tecnico"`
	ActivoID         int                `json:"activo_id"`
	EdificioID       int                `json:"edificio_id"`
	Tipo             TipoSolicitud      `json:"tipo"`
	Asunto           string             `json:"asunto"`
	Descripcion      string             `json:"descripcion"`
	Prioridad        PrioridadSolicitud `json:"prioridad"`
	Estado           EstadoSolicitud    `json:"estado"`
	FechaCreacion    time.Time          `json:"fecha_creacion"`
	FechaEnvio       *time.Time         `json:"fecha_envio,omitempty"`
	FechaRecepcion   *time.Time         `json:"fecha_recepcion,omitempty"`
	FechaCompletado  *time.Time         `json:"fecha_completado,omitempty"`
	MedioContacto    string             `json:"medio_contacto"`
	TelefonoContacto string             `json:"telefono_contacto,omitempty"`
	EmailContacto    string             `json:"email_contacto,omitempty"`
	RespuestaTecnico string             `json:"respuesta_tecnico,omitempty"`
}

// SolicitudFilter representa los filtros para buscar solicitudes
type SolicitudFilter struct {
	TecnicoID   *int                `json:"tecnico_id" form:"tecnico_id"`
	ResidenteID *int                `json:"residente_id" form:"residente_id"`
	ActivoID    *int                `json:"activo_id" form:"activo_id"`
	EdificioID  *int                `json:"edificio_id" form:"edificio_id"`
	Tipo        *TipoSolicitud      `json:"tipo" form:"tipo"`
	Estado      *EstadoSolicitud    `json:"estado" form:"estado"`
	Prioridad   *PrioridadSolicitud `json:"prioridad" form:"prioridad"`
	FechaDesde  *time.Time          `json:"fecha_desde" form:"fecha_desde"`
	FechaHasta  *time.Time          `json:"fecha_hasta" form:"fecha_hasta"`
	Busqueda    *string             `json:"busqueda" form:"busqueda"`
}

// EstadisticasSolicitudes representa estadísticas de solicitudes
type EstadisticasSolicitudes struct {
	TotalSolicitudes        int     `json:"total_solicitudes"`
	SolicitudesPendientes   int     `json:"solicitudes_pendientes"`
	SolicitudesEnProceso    int     `json:"solicitudes_en_proceso"`
	SolicitudesCompletadas  int     `json:"solicitudes_completadas"`
	SolicitudesCanceladas   int     `json:"solicitudes_canceladas"`
	TiempoPromedioRespuesta float64 `json:"tiempo_promedio_respuesta_horas"`
}

// SolicitudListResponse representa una lista paginada de solicitudes
type SolicitudListResponse struct {
	Solicitudes []SolicitudResponse `json:"solicitudes"`
	Total       int                 `json:"total"`
	Page        int                 `json:"page"`
	Limit       int                 `json:"limit"`
	TotalPages  int                 `json:"total_pages"`
}
