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

// Notificacion representa una notificación en el sistema
type Notificacion struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ActivoID  uint       `gorm:"not null;index" json:"activo_id"`
	Activo    Activo     `gorm:"foreignKey:ActivoID" json:"activo,omitempty"`
	UsuarioID *uint      `gorm:"index" json:"usuario_id"`
	Usuario   *Usuario   `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Mensaje   string     `gorm:"type:text;not null" json:"mensaje"`
	Enviado   bool       `gorm:"default:false" json:"enviado"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at"`
}

// NotificationRequest representa la estructura del JSON de entrada
type NotificationRequest struct {
	ActivoID uint `json:"activo_id" binding:"required"`
}

// TableName especifica el nombre de la tabla para GORM
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
