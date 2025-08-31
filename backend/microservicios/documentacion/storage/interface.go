package storage

import (
	"io"
	"mime/multipart"
)

// StorageService define la interfaz para el almacenamiento de archivos
type StorageService interface {
	UploadFile(file multipart.File, filename string) (rutaArchivo string, err error)
	GetFile(rutaArchivo string) (io.ReadCloser, error)
	DeleteFile(rutaArchivo string) error
}
