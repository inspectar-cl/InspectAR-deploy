package services

import (
	"fmt"
	"log"
	"notification/internal/models"
	"notification/internal/repository"
	"time"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	emailService     *EmailService
}

func NewNotificationService(repo *repository.NotificationRepository, emailSvc *EmailService) *NotificationService {
	return &NotificationService{
		notificationRepo: repo,
		emailService:     emailSvc,
	}
}

func (s *NotificationService) CreateNotification(activoID uint) (*models.Notificacion, error) {
	// Verificar que el activo existe y cargar las relaciones
	activo, err := s.notificationRepo.GetActivoByID(activoID)
	if err != nil {
		return nil, fmt.Errorf("activo no encontrado: %v", err)
	}

	// Crear la notificación
	notification := &models.Notificacion{
		ActivoID: activoID,
		Mensaje:  fmt.Sprintf("Alerta detectada en %s", activo.Nombre),
		Enviado:  false,
	}

	// Guardar la notificación en la base de datos
	createdNotification, err := s.notificationRepo.CreateNotification(notification)
	if err != nil {
		return nil, fmt.Errorf("error al crear notificación: %v", err)
	}

	// Obtener todos los usuarios del mismo edificio que el activo
	usuarios, err := s.notificationRepo.GetUsuariosByEdificioID(activo.EdificioID)
	if err != nil {
		log.Printf("Error al obtener usuarios del edificio %d: %v", activo.EdificioID, err)
		// Continuar aunque no se puedan enviar correos
		return createdNotification, nil
	}

	// Enviar correo electrónico a todos los usuarios del edificio
	if len(usuarios) > 0 {
		err = s.emailService.SendNotificationEmail(activo, usuarios, notification.Mensaje)
		if err != nil {
			log.Printf("Error al enviar correo: %v", err)
			// No retornar error, ya que la notificación se creó correctamente
		} else {
			log.Printf("Correo enviado exitosamente a %d usuarios", len(usuarios))
		}
	}

	return createdNotification, nil
}

func (s *NotificationService) GetNotificationsByActivoID(activoID uint) ([]models.Notificacion, error) {
	return s.notificationRepo.GetNotificationsByActivoID(activoID)
}

func (s *NotificationService) SendNotification(notificationID uint) error {
	// Obtener la notificación
	notificacion, err := s.notificationRepo.GetNotificationByID(notificationID)
	if err != nil {
		return fmt.Errorf("notificación no encontrada: %v", err)
	}

	// Aquí iría la lógica para enviar la notificación
	// Por ejemplo: enviar email, SMS, push notification, etc.

	// Actualizar el estado a enviado
	err = s.notificationRepo.UpdateNotificationStatus(notificationID, true)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %v", err)
	}

	// Actualizar el timestamp de envío
	now := time.Now()
	notificacion.SentAt = &now

	return nil
}

func (s *NotificationService) GetAllEdificios() ([]models.Edificio, error) {
	return s.notificationRepo.GetAllEdificios()
}

func (s *NotificationService) GetEdificioByID(id uint) (*models.Edificio, error) {
	return s.notificationRepo.GetEdificioByID(id)
}

func (s *NotificationService) GetUsuariosByEdificioID(edificioID uint) ([]models.Usuario, error) {
	return s.notificationRepo.GetUsuariosByEdificioID(edificioID)
}

func (s *NotificationService) GetActivosByEdificioID(edificioID uint) ([]models.Activo, error) {
	return s.notificationRepo.GetActivosByEdificioID(edificioID)
}
