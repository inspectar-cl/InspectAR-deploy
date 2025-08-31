package services

import (
	"database/sql"
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
)

type SolicitudService struct {
	solicitudRepo *repository.SolicitudRepository
}

func NewSolicitudService(db *sql.DB) *SolicitudService {
	return &SolicitudService{
		solicitudRepo: repository.NewSolicitudRepository(db),
	}
}

// CreateSolicitud crea una nueva solicitud de servicio técnico
func (s *SolicitudService) CreateSolicitud(req *models.CreateSolicitudRequest) (*models.SolicitudResponse, error) {
	// Validar datos de entrada
	if err := s.ValidateCreateSolicitudRequest(req); err != nil {
		return nil, err
	}

	// Crear la solicitud
	solicitud, err := s.solicitudRepo.CreateSolicitud(req)
	if err != nil {
		return nil, fmt.Errorf("error al crear solicitud: %w", err)
	}

	// Convertir a response
	response := s.convertToResponse(solicitud)
	return response, nil
}

// GetSolicitudByID obtiene una solicitud por su ID
func (s *SolicitudService) GetSolicitudByID(id int) (*models.SolicitudResponse, error) {
	solicitud, err := s.solicitudRepo.GetSolicitudByID(id)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(solicitud)
	return response, nil
}

// GetSolicitudesByFilter obtiene solicitudes con filtros
func (s *SolicitudService) GetSolicitudesByFilter(filter models.SolicitudFilter, page, limit int) (*models.SolicitudListResponse, error) {
	offset := (page - 1) * limit

	solicitudes, total, err := s.solicitudRepo.GetSolicitudesByFilter(filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al obtener solicitudes: %w", err)
	}

	// Convertir a responses
	responses := make([]models.SolicitudResponse, len(solicitudes))
	for i, solicitud := range solicitudes {
		responses[i] = *s.convertToResponse(&solicitud)
	}

	return &models.SolicitudListResponse{
		Solicitudes: responses,
		Total:       int(total),
		Page:        page,
		Limit:       limit,
		TotalPages:  int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}

// EnviarSolicitud marca una solicitud como enviada
func (s *SolicitudService) EnviarSolicitud(id int) error {
	return s.solicitudRepo.MarkSolicitudAsEnviada(id)
}

// ActualizarEstadoSolicitud actualiza el estado de una solicitud
func (s *SolicitudService) ActualizarEstadoSolicitud(id int, estado models.EstadoSolicitud, respuesta string) error {
	return s.solicitudRepo.UpdateSolicitudEstado(id, estado, respuesta)
}

// GetEstadisticasSolicitudes obtiene estadísticas de solicitudes
func (s *SolicitudService) GetEstadisticasSolicitudes(tecnicoID, edificioID *int) (*models.EstadisticasSolicitudes, error) {
	return s.solicitudRepo.GetEstadisticasSolicitudes(tecnicoID, edificioID)
}

// Helper methods

// convertToResponse convierte una solicitud del modelo de base de datos a response
func (s *SolicitudService) convertToResponse(solicitud *models.SolicitudTecnico) *models.SolicitudResponse {
	response := &models.SolicitudResponse{
		ID:               solicitud.ID,
		TecnicoID:        solicitud.TecnicoID,
		ActivoID:         solicitud.ActivoID,
		EdificioID:       solicitud.EdificioID,
		Tipo:             solicitud.Tipo,
		Asunto:           solicitud.Asunto,
		Descripcion:      solicitud.Descripcion,
		Prioridad:        solicitud.Prioridad,
		Estado:           solicitud.Estado,
		FechaCreacion:    solicitud.FechaCreacion,
		FechaEnvio:       solicitud.FechaEnvio,
		FechaRecepcion:   solicitud.FechaRecepcion,
		FechaCompletado:  solicitud.FechaCompletado,
		MedioContacto:    solicitud.MedioContacto,
		TelefonoContacto: solicitud.TelefonoContacto,
		EmailContacto:    solicitud.EmailContacto,
		RespuestaTecnico: solicitud.RespuestaTecnico,
	}

	// Agregar información del técnico si está disponible
	if solicitud.Tecnico.ID != 0 {
		response.Tecnico = models.TecnicoResponse{
			ID:             solicitud.Tecnico.ID,
			Nombre:         solicitud.Tecnico.Nombre,
			Apellido:       solicitud.Tecnico.Apellido,
			NombreCompleto: fmt.Sprintf("%s %s", solicitud.Tecnico.Nombre, solicitud.Tecnico.Apellido),
			Email:          solicitud.Tecnico.Email,
			Telefono:       solicitud.Tecnico.Telefono,
			Especialidad:   solicitud.Tecnico.Especialidad,
		}

		// Agregar información de la empresa si está disponible
		if solicitud.Tecnico.Empresa.ID != 0 {
			response.Tecnico.Empresa = models.EmpresaResponse{
				ID:       solicitud.Tecnico.Empresa.ID,
				Nombre:   solicitud.Tecnico.Empresa.Nombre,
				RUT:      solicitud.Tecnico.Empresa.RUT,
				Telefono: solicitud.Tecnico.Empresa.Telefono,
				Email:    solicitud.Tecnico.Empresa.Email,
			}
		}
	}

	return response
}

// ValidateCreateSolicitudRequest valida los datos de creación de solicitud
func (s *SolicitudService) ValidateCreateSolicitudRequest(req *models.CreateSolicitudRequest) error {
	if req.TecnicoID <= 0 {
		return fmt.Errorf("el ID del técnico es requerido")
	}

	if req.Asunto == "" {
		return fmt.Errorf("el asunto es requerido")
	}

	if req.Descripcion == "" {
		return fmt.Errorf("la descripción es requerida")
	}

	if req.MedioContacto == "" {
		return fmt.Errorf("el medio de contacto es requerido")
	}

	// Validar que se proporcione al menos un medio de contacto
	if req.MedioContacto == "telefono" && req.TelefonoContacto == "" {
		return fmt.Errorf("el teléfono de contacto es requerido cuando el medio de contacto es telefono")
	}

	if req.MedioContacto == "email" && req.EmailContacto == "" {
		return fmt.Errorf("el email de contacto es requerido cuando el medio de contacto es email")
	}

	if req.MedioContacto == "ambos" && (req.TelefonoContacto == "" || req.EmailContacto == "") {
		return fmt.Errorf("tanto el teléfono como el email son requeridos cuando el medio de contacto es ambos")
	}

	return nil
}
