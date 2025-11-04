package services

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

type QRService struct {
	baseURL string // URL base de la aplicación
}

func NewQRService(baseURL string) *QRService {
	return &QRService{
		baseURL: baseURL,
	}
}

// GenerarQRActivo genera un código QR para un activo
// Retorna el QR en Base64 y la URL codificada
func (s *QRService) GenerarQRActivo(activoID int, codigoActivo string) (string, string, error) {
	// URL que apunta a la información del activo usando el código
	// Formato: http://localhost:8092/qr/EDI01-BOMBA-0001-2025
	url := fmt.Sprintf("%s/qr/%s", s.baseURL, codigoActivo)

	// Generar QR code (256x256 px, nivel de corrección Medium)
	qrCode, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", "", fmt.Errorf("error generando QR: %w", err)
	}

	// Convertir a Base64 para almacenar en BD
	qrBase64 := base64.StdEncoding.EncodeToString(qrCode)

	return qrBase64, url, nil
}

// DecodificarQR decodifica un QR desde Base64 a bytes
func (s *QRService) DecodificarQR(qrBase64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(qrBase64)
}
