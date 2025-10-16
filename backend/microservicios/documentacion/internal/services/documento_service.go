package services

import (
	"documentacion/internal/models"
	"documentacion/internal/repository"
	"documentacion/storage"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type DocumentoService struct {
	repo           *repository.DocumentoRepository
	storageService storage.StorageService
}

func NewDocumentoService(repo *repository.DocumentoRepository, storageService storage.StorageService) *DocumentoService {
	return &DocumentoService{
		repo:           repo,
		storageService: storageService,
	}
}

// SubirDocumento sube un archivo y crea el registro en BD (HdU05)
func (s *DocumentoService) SubirDocumento(req *models.CreateDocumentoRequest, file multipart.File, header *multipart.FileHeader) (*models.Documento, error) {
	// Validar extensión
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" && ext != ".docx" {
		return nil, fmt.Errorf("tipo de archivo no permitido: %s", ext)
	}

	// Validar tamaño (50MB max)
	if header.Size > 50*1024*1024 {
		return nil, fmt.Errorf("archivo demasiado grande. Máximo 50MB")
	}

	// Generar nombre único para el archivo
	timestamp := time.Now().Format("20060102_150405")
	nombreArchivo := fmt.Sprintf("%d_%s_%s", req.ActivoID, timestamp, header.Filename)

	// Subir a storage
	rutaArchivo, err := s.storageService.UploadFile(file, nombreArchivo)
	if err != nil {
		return nil, fmt.Errorf("error subiendo archivo: %v", err)
	}

	// Crear registro en BD
	documento, err := s.repo.Create(req, rutaArchivo, header.Size, ext[1:]) // sin el punto
	if err != nil {
		// Si falla la BD, intentar limpiar el archivo subido
		s.storageService.DeleteFile(rutaArchivo)
		return nil, fmt.Errorf("error creando registro: %v", err)
	}

	return documento, nil
}

// ObtenerDocumento obtiene un documento por ID
func (s *DocumentoService) ObtenerDocumento(id int) (*models.Documento, error) {
	return s.repo.GetByID(id)
}

// ObtenerDocumentosPorActivo obtiene todos los documentos de un activo
func (s *DocumentoService) ObtenerDocumentosPorActivo(activoID int) ([]models.Documento, error) {
	return s.repo.GetByActivoID(activoID)
}

// ListarDocumentos obtiene todos los documentos con filtros opcionales
func (s *DocumentoService) ListarDocumentos(activoID *int, soloFichasTecnicas bool) ([]models.Documento, error) {
	return s.repo.GetAll(activoID, soloFichasTecnicas)
}

// ActualizarDocumento actualiza un documento existente
func (s *DocumentoService) ActualizarDocumento(id int, req *models.UpdateDocumentoRequest) (*models.Documento, error) {
	return s.repo.Update(id, req)
}

// ObtenerFichaTecnica obtiene la ficha técnica de un activo (HdU23)
func (s *DocumentoService) ObtenerFichaTecnica(activoID int) (*models.Documento, error) {
	return s.repo.GetFichaTecnicaPorActivo(activoID)
}

// DescargarArchivo obtiene el contenido del archivo para descarga
func (s *DocumentoService) DescargarArchivo(documentoID int) (io.ReadCloser, string, error) {
	// Obtener información del documento
	documento, err := s.repo.GetByID(documentoID)
	if err != nil {
		return nil, "", fmt.Errorf("documento no encontrado: %v", err)
	}

	// Obtener archivo del storage
	file, err := s.storageService.GetFile(documento.RutaArchivo)
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo archivo: %v", err)
	}

	return file, documento.Nombre, nil
}

// BuscarDocumentos busca documentos con filtros (HdU05)
func (s *DocumentoService) BuscarDocumentos(filtros models.DocumentoFiltros) ([]models.Documento, error) {
	return s.repo.BuscarConFiltros(filtros)
}

// ObtenerHistorialMantenimiento construye el historial completo de un activo
func (s *DocumentoService) ObtenerHistorialMantenimiento(activoID int) (*models.HistorialMantenimiento, error) {
	return s.repo.GetHistorialMantenimiento(activoID)
}

// ValidarPermisosUpload valida si un técnico puede subir documentos
func (s *DocumentoService) ValidarPermisosUpload(tecnicoID int) error {
	// Aquí podrías hacer una llamada al microservicio de gestión
	// para verificar que el técnico esté autorizado
	// Por ahora, simplemente permitimos si el ID es válido
	if tecnicoID <= 0 {
		return fmt.Errorf("técnico no válido")
	}
	return nil
}
