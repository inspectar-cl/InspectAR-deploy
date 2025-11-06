package models

import (
	"time"
)

// Edificio representa un edificio en el sistema
type Edificio struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Direccion     string    `gorm:"size:255;not null" json:"direccion"`
	NumeroActivos int       `gorm:"default:0" json:"numero_activos"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Usuarios      []Usuario `gorm:"foreignKey:EdificioID" json:"usuarios,omitempty"`
	Activos       []Activo  `gorm:"foreignKey:EdificioID" json:"activos,omitempty"`
}

// Usuario representa un usuario en el sistema
type Usuario struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Scope      string    `gorm:"size:100;not null" json:"scope"`
	Usuario    string    `gorm:"size:100;uniqueIndex;not null" json:"usuario"`
	Correo     string    `gorm:"size:255;uniqueIndex;not null" json:"correo"`
	Numero     string    `gorm:"size:20" json:"numero"`
	EdificioID *uint     `gorm:"index" json:"edificio_id"`
	Edificio   *Edificio `gorm:"foreignKey:EdificioID" json:"edificio,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Activo representa un activo en el sistema
type Activo struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nombre     string    `gorm:"size:255;not null" json:"nombre"`
	EdificioID uint      `gorm:"not null;index" json:"edificio_id"`
	Edificio   Edificio  `gorm:"foreignKey:EdificioID" json:"edificio,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TipoNotificacion representa los tipos de notificación disponibles
type TipoNotificacion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Nombre      string    `gorm:"size:50;uniqueIndex;not null" json:"nombre"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Prioridad   int       `gorm:"default:1" json:"prioridad"` // 1=Low, 2=Medium, 3=High, 4=Critical
	Activo      bool      `gorm:"default:true" json:"activo"`
	CreatedAt   time.Time `json:"created_at"`
}

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID               int       `json:"id" db:"id"`
	BuildingID       *int      `json:"building_id" db:"building_id"`
	AssetID          *int      `json:"asset_id" db:"asset_id"`
	SensorID         string    `json:"sensor_id" db:"sensor_id"`
	Message          string    `json:"message" db:"message"`
	AlertType        string    `json:"alert_type" db:"alert_type"`
	Tipo             string    `json:"tipo" db:"tipo"`           // sensor, alerta, mantenimiento, sistema
	Prioridad        string    `json:"prioridad" db:"prioridad"` // low, medium, high, critical
	NotificationMail bool      `json:"notification_mail" db:"notification_mail"`
	NotificationSMS  bool      `json:"notification_sms" db:"notification_sms"`
	Status           string    `json:"status" db:"status"` // pending, sent, failed
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type TipoSensor struct {
	ID          int    `json:"id" db:"id"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Unidad      string `json:"unidad" db:"unidad"`
}

type TipoActivo struct {
	ID          int    `json:"id" db:"id"`
	Nombre      string `json:"nombre" db:"nombre"`
	Descripcion string `json:"descripcion" db:"descripcion"`
	Categoria   string `json:"categoria" db:"categoria"`
}

type NivelPrioridad struct {
	ID          int    `json:"id" db:"id"`
	Nombre      string `json:"nombre" db:"nombre"`
	Valor       int    `json:"valor" db:"valor"`
	Color       string `json:"color" db:"color"`
	Descripcion string `json:"descripcion" db:"descripcion"`
}

// NotificationRequest representa la estructura del JSON de entrada
type NotificationRequest struct {
	ActivoID uint `json:"activo_id" binding:"required"`
}

// SensorAlertRequest representa la estructura para alertas de sensor
type SensorAlertRequest struct {
	SensorID         string `json:"sensor_id" binding:"required"`
	SensorData       string `json:"sensor_data"`
	BuildingID       int    `json:"building_id"`
	AssetID          int    `json:"asset_id"`
	Message          string `json:"message" binding:"required"`
	AlertType        string `json:"alert_type" binding:"required"`
	Tipo             string `json:"tipo" binding:"required"` // sensor, alerta, mantenimiento, sistema
	Prioridad        string `json:"prioridad"`               // low, medium, high, critical (default: medium)
	NotificationMail bool   `json:"notification_mail"`
	NotificationSMS  bool   `json:"notification_sms"`
}

// TechnicianContactRequest representa la estructura para solicitudes de contacto con técnicos
type TechnicianContactRequest struct {
	TechnicianEmail string `json:"technician_email" binding:"required,email"` // Email del técnico al cual enviar
	ActivoID        int    `json:"activo_id" binding:"required"`              // ID del activo relacionado
	UserName        string `json:"user_name" binding:"required"`              // Nombre del usuario que solicita contacto
	UserEmail       string `json:"user_email" binding:"required,email"`       // Email del usuario para respuesta
	Message         string `json:"message" binding:"required"`                // Mensaje personalizado del usuario
	Priority        string `json:"priority"`                                  // Prioridad: low, medium, high, urgent
}

// Ticket representa una solicitud de ingreso, modificación o eliminación
type Ticket struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TipoEntidad   string `gorm:"size:50;not null" json:"tipo_entidad"`               // edificio, activo, tecnico
	TipoOperacion string `gorm:"size:50;not null" json:"tipo_operacion"`             // ingreso, modificacion, eliminacion
	Estado        string `gorm:"size:50;not null;default:no_resuelto" json:"estado"` // resuelto, no_resuelto

	// Usuario solicitante
	UsuarioID    *uint  `gorm:"index" json:"usuario_id"`
	UsuarioEmail string `gorm:"size:255;not null" json:"usuario_email"`

	// Datos para EDIFICIO
	EdificioID        *uint    `gorm:"index" json:"edificio_id,omitempty"`
	EdificioNombre    string   `gorm:"size:255" json:"edificio_nombre,omitempty"`
	EdificioDireccion string   `gorm:"size:255" json:"edificio_direccion,omitempty"`
	EdificioLatitud   *float64 `json:"edificio_latitud,omitempty"`
	EdificioLongitud  *float64 `json:"edificio_longitud,omitempty"`

	// Datos para ACTIVO
	ActivoID          *uint  `gorm:"index" json:"activo_id,omitempty"`
	ActivoNombre      string `gorm:"size:255" json:"activo_nombre,omitempty"`
	ActivoTipo        string `gorm:"size:100" json:"activo_tipo,omitempty"`
	ActivoDescripcion string `gorm:"type:text" json:"activo_descripcion,omitempty"`
	ActivoUbicacion   string `gorm:"size:255" json:"activo_ubicacion,omitempty"`
	ActivoEdificioID  *uint  `json:"activo_edificio_id,omitempty"`

	// Datos para TECNICO
	TecnicoID           *int   `gorm:"index" json:"tecnico_id,omitempty"`
	TecnicoNombre       string `gorm:"size:255" json:"tecnico_nombre,omitempty"`
	TecnicoEmail        string `gorm:"size:255" json:"tecnico_email,omitempty"`
	TecnicoTelefono     string `gorm:"size:20" json:"tecnico_telefono,omitempty"`
	TecnicoEspecialidad string `gorm:"size:100" json:"tecnico_especialidad,omitempty"`
	TecnicoAutorizado   *bool  `json:"tecnico_autorizado,omitempty"`

	// Información adicional
	Justificacion   string `gorm:"type:text" json:"justificacion,omitempty"`
	ComentarioAdmin string `gorm:"type:text" json:"comentario_admin,omitempty"`

	// Metadata
	CreatedAt        time.Time  `json:"created_at"`
	FechaResolucion  *time.Time `json:"fecha_resolucion,omitempty"`
	ResueltoPor      *uint      `json:"resuelto_por,omitempty"`
	ResueltoPorEmail string     `gorm:"size:255" json:"resuelto_por_email,omitempty"`
}

// CreateTicketRequest representa la estructura para crear un ticket
type CreateTicketRequest struct {
	TipoEntidad   string `json:"tipo_entidad" binding:"required,oneof=edificio activo tecnico"`
	TipoOperacion string `json:"tipo_operacion" binding:"required,oneof=ingreso modificacion eliminacion"`
	UsuarioEmail  string `json:"usuario_email" binding:"required,email"`

	// Campos opcionales para edificio
	EdificioID        *uint    `json:"edificio_id,omitempty"`
	EdificioNombre    string   `json:"edificio_nombre,omitempty"`
	EdificioDireccion string   `json:"edificio_direccion,omitempty"`
	EdificioLatitud   *float64 `json:"edificio_latitud,omitempty"`
	EdificioLongitud  *float64 `json:"edificio_longitud,omitempty"`

	// Campos opcionales para activo
	ActivoID          *uint  `json:"activo_id,omitempty"`
	ActivoNombre      string `json:"activo_nombre,omitempty"`
	ActivoTipo        string `json:"activo_tipo,omitempty"`
	ActivoDescripcion string `json:"activo_descripcion,omitempty"`
	ActivoUbicacion   string `json:"activo_ubicacion,omitempty"`
	ActivoEdificioID  *uint  `json:"activo_edificio_id,omitempty"`

	// Campos opcionales para técnico
	TecnicoID           *int   `json:"tecnico_id,omitempty"`
	TecnicoNombre       string `json:"tecnico_nombre,omitempty"`
	TecnicoEmail        string `json:"tecnico_email,omitempty"`
	TecnicoTelefono     string `json:"tecnico_telefono,omitempty"`
	TecnicoEspecialidad string `json:"tecnico_especialidad,omitempty"`
	TecnicoAutorizado   *bool  `json:"tecnico_autorizado,omitempty"`

	Justificacion string `json:"justificacion,omitempty"`
}

// ResolveTicketRequest representa la estructura para resolver un ticket
type ResolveTicketRequest struct {
	ResueltoPorEmail string `json:"resuelto_por_email" binding:"required,email"`
	ComentarioAdmin  string `json:"comentario_admin" binding:"required"`
}

// TicketsPaginationResponse representa la respuesta paginada de tickets
type TicketsPaginationResponse struct {
	Tickets        []Ticket `json:"tickets"`
	TotalTickets   int64    `json:"total_tickets"`
	TotalPaginas   int      `json:"total_paginas"`
	PaginaActual   int      `json:"pagina_actual"`
	ItemsPorPagina int      `json:"items_por_pagina"`
}

// TableName especifica el nombre de la tabla para GORM
func (TipoNotificacion) TableName() string {
	return "tipos_notificacion"
}

func (Edificio) TableName() string {
	return "edificios"
}

func (Usuario) TableName() string {
	return "usuarios"
}

func (Activo) TableName() string {
	return "activos"
}

func (Notificacion) TableName() string {
	return "notificaciones"
}

func (Ticket) TableName() string {
	return "tickets"
}
