package storage

import (
	"context"
	"documentacion/config"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(cfg *config.MinioConfig) (*MinioService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("error inicializando cliente MinIO: %v", err)
	}

	// Verificar/crear bucket
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("error verificando bucket: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("error creando bucket: %v", err)
		}
	}

	return &MinioService{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (s *MinioService) UploadFile(file multipart.File, filename string) (string, error) {
	ctx := context.Background()
	
	// Subir archivo
	_, err := s.client.PutObject(ctx, s.bucket, filename, file, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("error subiendo archivo a MinIO: %v", err)
	}

	return filename, nil
}

func (s *MinioService) GetFile(rutaArchivo string) (io.ReadCloser, error) {
	ctx := context.Background()
	
	object, err := s.client.GetObject(ctx, s.bucket, rutaArchivo, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo archivo de MinIO: %v", err)
	}

	return object, nil
}

func (s *MinioService) DeleteFile(rutaArchivo string) error {
	ctx := context.Background()
	
	err := s.client.RemoveObject(ctx, s.bucket, rutaArchivo, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error eliminando archivo de MinIO: %v", err)
	}

	return nil
}
