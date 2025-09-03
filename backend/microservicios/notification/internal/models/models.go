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
