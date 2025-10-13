package services

import (
	"fmt"
	"gestion/internal/models"
	"gestion/internal/repository"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FirmaService struct {
	firmaRepo   *repository.FirmaRepository
	usuarioRepo *repository.UsuarioRepository
	storagePath string
}

func NewFirmaService(firmaRepo *repository.FirmaRepository, usuarioRepo *repository.UsuarioRepository) *FirmaService {
	// Directorio donde se guardarán las firmas
	storagePath := os.Getenv("FIRMAS_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage/firmas"
	}

	// Crear el directorio si no existe
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		fmt.Printf("Error al crear directorio de firmas: %v\n", err)
	}

	return &FirmaService{
		firmaRepo:   firmaRepo,
		usuarioRepo: usuarioRepo,
		storagePath: storagePath,
	}
}

// CrearFirmaDesdeArchivo crea una firma a partir de un archivo subido
func (s *FirmaService) CrearFirmaDesdeArchivo(email string, nombreArchivo string, esPredeterminada bool, file multipart.File, header *multipart.FileHeader) (*models.FirmaDigital, error) {
	// Obtener usuario por email
	usuario, err := s.usuarioRepo.GetUsuarioByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}

	usuarioID := usuario.ID

	// Validar tipo de archivo
	tipoMime := header.Header.Get("Content-Type")
	if !s.esFormatoValido(tipoMime) {
		return nil, fmt.Errorf("formato de archivo no válido. Formatos aceptados: PNG, JPEG, JPG, SVG")
	}

	// Generar nombre único para el archivo
	extension := filepath.Ext(header.Filename)
	if extension == "" {
		extension = s.extensionDesdeMime(tipoMime)
	}

	timestamp := time.Now().Unix()
	nombreArchivoUnico := fmt.Sprintf("%d_%d%s", usuarioID, timestamp, extension)
	rutaArchivo := filepath.Join(s.storagePath, nombreArchivoUnico)

	// Guardar el archivo
	dst, err := os.Create(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("error al guardar archivo: %v", err)
	}

	// Obtener información del archivo
	fileInfo, err := os.Stat(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %v", err)
	}

	// Si no se proporcionó nombre, usar el original
	if nombreArchivo == "" {
		nombreArchivo = header.Filename
	}

	// Crear registro en base de datos
	firma := &models.FirmaDigital{
		UsuarioID:        usuarioID,
		NombreArchivo:    nombreArchivo,
		RutaArchivo:      rutaArchivo,
		TipoMime:         tipoMime,
		Formato:          strings.TrimPrefix(extension, "."),
		TamanoBytes:      fileInfo.Size(),
		EsPredeterminada: esPredeterminada,
	}

	err = s.firmaRepo.Create(firma)
	if err != nil {
		// Eliminar archivo si falla la creación en BD
		os.Remove(rutaArchivo)
		return nil, err
	}

	return firma, nil
}

// CrearFirmaDesdeSVG crea una firma a partir de datos SVG (para pizarra de dibujo)
func (s *FirmaService) CrearFirmaDesdeSVG(email string, nombreArchivo string, datosSVG string, esPredeterminada bool) (*models.FirmaDigital, error) {
	// Obtener usuario por email
	usuario, err := s.usuarioRepo.GetUsuarioByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}

	usuarioID := usuario.ID

	// Generar nombre único para el archivo
	timestamp := time.Now().Unix()
	nombreArchivoUnico := fmt.Sprintf("%d_%d.svg", usuarioID, timestamp)
	rutaArchivo := filepath.Join(s.storagePath, nombreArchivoUnico)

	// Guardar el archivo SVG
	err = os.WriteFile(rutaArchivo, []byte(datosSVG), 0644)
	if err != nil {
		return nil, fmt.Errorf("error al guardar archivo SVG: %v", err)
	}

	// Obtener tamaño del archivo
	fileInfo, err := os.Stat(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %v", err)
	}

	// Crear registro en base de datos
	firma := &models.FirmaDigital{
		UsuarioID:        usuarioID,
		NombreArchivo:    nombreArchivo,
		RutaArchivo:      rutaArchivo,
		TipoMime:         "image/svg+xml",
		Formato:          "svg",
		DatosFirma:       []byte(datosSVG),
		TamanoBytes:      fileInfo.Size(),
		EsPredeterminada: esPredeterminada,
	}

	err = s.firmaRepo.Create(firma)
	if err != nil {
		// Eliminar archivo si falla la creación en BD
		os.Remove(rutaArchivo)
		return nil, err
	}

	return firma, nil
}

// ObtenerFirma obtiene una firma por ID
func (s *FirmaService) ObtenerFirma(id int) (*models.FirmaDigital, error) {
	return s.firmaRepo.GetByID(id)
}

// ObtenerFirmasUsuario obtiene todas las firmas de un usuario
func (s *FirmaService) ObtenerFirmasUsuario(usuarioID int) ([]models.FirmaDigital, error) {
	return s.firmaRepo.GetByUsuarioID(usuarioID)
}

// ObtenerFirmasUsuarioPorEmail obtiene todas las firmas de un usuario por email
func (s *FirmaService) ObtenerFirmasUsuarioPorEmail(email string) ([]models.FirmaDigital, error) {
	return s.firmaRepo.GetByUsuarioEmail(email)
}

// ObtenerFirmaPredeterminada obtiene la firma predeterminada de un usuario
func (s *FirmaService) ObtenerFirmaPredeterminada(usuarioID int) (*models.FirmaDigital, error) {
	return s.firmaRepo.GetDefaultByUsuario(usuarioID)
}

// ObtenerFirmaPredeterminadaPorEmail obtiene la firma predeterminada de un usuario por email
func (s *FirmaService) ObtenerFirmaPredeterminadaPorEmail(email string) (*models.FirmaDigital, error) {
	return s.firmaRepo.GetDefaultByUsuarioEmail(email)
}

// ActualizarFirma actualiza los datos de una firma
func (s *FirmaService) ActualizarFirma(id int, nombreArchivo string, esPredeterminada bool) (*models.FirmaDigital, error) {
	// Obtener firma existente
	firma, err := s.firmaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if nombreArchivo != "" {
		firma.NombreArchivo = nombreArchivo
	}
	firma.EsPredeterminada = esPredeterminada

	err = s.firmaRepo.Update(firma)
	if err != nil {
		return nil, err
	}

	return firma, nil
}

// EliminarFirma elimina una firma
func (s *FirmaService) EliminarFirma(id int) error {
	// Obtener firma para eliminar archivo físico
	firma, err := s.firmaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar de base de datos
	err = s.firmaRepo.Delete(id)
	if err != nil {
		return err
	}

	// Eliminar archivo físico
	if err := os.Remove(firma.RutaArchivo); err != nil {
		fmt.Printf("Advertencia: no se pudo eliminar archivo físico: %v\n", err)
	}

	return nil
}

// EstablecerComoPredeterminada marca una firma como predeterminada
func (s *FirmaService) EstablecerComoPredeterminada(id int, usuarioID int) error {
	return s.firmaRepo.SetAsDefault(id, usuarioID)
}

// EstablecerComoPredeterminadaPorEmail marca una firma como predeterminada usando email
func (s *FirmaService) EstablecerComoPredeterminadaPorEmail(id int, email string) error {
	// Obtener usuario por email
	usuario, err := s.usuarioRepo.GetUsuarioByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado: %v", err)
	}

	return s.firmaRepo.SetAsDefault(id, usuario.ID)
}

// esFormatoValido valida si el tipo MIME es aceptado
func (s *FirmaService) esFormatoValido(tipoMime string) bool {
	formatosValidos := []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
		"image/svg+xml",
	}

	for _, formato := range formatosValidos {
		if tipoMime == formato {
			return true
		}
	}
	return false
}

// extensionDesdeMime obtiene la extensión a partir del tipo MIME
func (s *FirmaService) extensionDesdeMime(tipoMime string) string {
	switch tipoMime {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".bin"
	}
}
