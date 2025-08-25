package storage

import (
	"documentacion/config"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalService struct {
	uploadDir string
}

func NewLocalService(cfg *config.LocalConfig) (*LocalService, error) {
	// Crear directorio si no existe
	err := os.MkdirAll(cfg.UploadDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creando directorio de uploads: %v", err)
	}

	return &LocalService{
		uploadDir: cfg.UploadDir,
	}, nil
}

func (s *LocalService) UploadFile(file multipart.File, filename string) (string, error) {
	// Crear ruta completa
	rutaCompleta := filepath.Join(s.uploadDir, filename)

	// Crear archivo en el sistema de archivos
	dst, err := os.Create(rutaCompleta)
	if err != nil {
		return "", fmt.Errorf("error creando archivo: %v", err)
	}
	defer dst.Close()

	// Copiar contenido
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("error copiando archivo: %v", err)
	}

	return filename, nil
}

func (s *LocalService) GetFile(rutaArchivo string) (io.ReadCloser, error) {
	rutaCompleta := filepath.Join(s.uploadDir, rutaArchivo)
	
	file, err := os.Open(rutaCompleta)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %v", err)
	}

	return file, nil
}

func (s *LocalService) DeleteFile(rutaArchivo string) error {
	rutaCompleta := filepath.Join(s.uploadDir, rutaArchivo)
	
	err := os.Remove(rutaCompleta)
	if err != nil {
		return fmt.Errorf("error eliminando archivo: %v", err)
	}

	return nil
}
