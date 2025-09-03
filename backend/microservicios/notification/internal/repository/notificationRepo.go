package repository

import (
	"notification/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) CreateNotification(notification *models.Notificacion) (*models.Notificacion, error) {
	if err := r.db.Create(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *NotificationRepository) GetNotificationsByActivoID(activoID uint) ([]models.Notificacion, error) {
	var notificaciones []models.Notificacion
	err := r.db.Preload("Activo").Preload("Usuario").Where("activo_id = ?", activoID).Find(&notificaciones).Error
	return notificaciones, err
}

func (r *NotificationRepository) UpdateNotification(notification *models.Notificacion) error {
	return r.db.Save(notification).Error
}

// Métodos para Tipos de Notificación
func (r *NotificationRepository) GetTipoNotificacionByNombre(nombre string) (*models.TipoNotificacion, error) {
	var tipo models.TipoNotificacion
	err := r.db.Where("nombre = ? AND activo = ?", nombre, true).First(&tipo).Error
	return &tipo, err
}

func (r *NotificationRepository) GetAllTiposNotificacion() ([]models.TipoNotificacion, error) {
	var tipos []models.TipoNotificacion
	err := r.db.Where("activo = ?", true).Find(&tipos).Error
	return tipos, err
}

// Métodos adicionales para Usuarios
func (r *NotificationRepository) GetUsuariosByScope(scope string) ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Preload("Edificio").Where("scope = ?", scope).Find(&usuarios).Error
	return usuarios, err
}

func (r *NotificationRepository) GetAllUsuarios() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Preload("Edificio").Find(&usuarios).Error
	return usuarios, err
}

func (r *NotificationRepository) UpdateNotificationStatus(id uint, enviado bool) error {
	return r.db.Model(&models.Notificacion{}).Where("id = ?", id).Update("enviado", enviado).Error
}

func (r *NotificationRepository) GetNotificationByID(id uint) (*models.Notificacion, error) {
	var notificacion models.Notificacion
	err := r.db.Preload("Activo").Preload("Usuario").First(&notificacion, id).Error
	return &notificacion, err
}

// Métodos para Edificios
func (r *NotificationRepository) GetAllEdificios() ([]models.Edificio, error) {
	var edificios []models.Edificio
	err := r.db.Preload("Usuarios").Preload("Activos").Find(&edificios).Error
	return edificios, err
}

func (r *NotificationRepository) GetEdificioByID(id uint) (*models.Edificio, error) {
	var edificio models.Edificio
	err := r.db.Preload("Usuarios").Preload("Activos").First(&edificio, id).Error
	return &edificio, err
}

// Métodos para Usuarios
func (r *NotificationRepository) GetUsuariosByEdificioID(edificioID uint) ([]models.Usuario, error) {
	var usuarios []models.Usuario
	err := r.db.Preload("Edificio").Where("edificio_id = ?", edificioID).Find(&usuarios).Error
	return usuarios, err
}

// Métodos para Activos
func (r *NotificationRepository) GetActivosByEdificioID(edificioID uint) ([]models.Activo, error) {
	var activos []models.Activo
	err := r.db.Preload("Edificio").Where("edificio_id = ?", edificioID).Find(&activos).Error
	return activos, err
}

func (r *NotificationRepository) GetActivoByID(id uint) (*models.Activo, error) {
	var activo models.Activo
	err := r.db.Preload("Edificio").First(&activo, id).Error
	return &activo, err
}
