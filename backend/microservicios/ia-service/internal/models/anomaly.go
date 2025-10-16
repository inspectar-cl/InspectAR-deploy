package models

import "time"

// Anomaly representa una anomalía detectada en los datos de un sensor
type Anomaly struct {
	ID                int       `json:"id" db:"id"`
	ActivoID          int       `json:"activo_id" db:"activo_id"`
	SensorID          string    `json:"sensor_id" db:"sensor_id"`
	Timestamp         time.Time `json:"timestamp" db:"timestamp"`
	AnomalyScore      float64   `json:"anomaly_score" db:"anomaly_score"`
	AnomalyLikelihood float64   `json:"anomaly_likelihood" db:"anomaly_likelihood"`
	Severidad         string    `json:"severidad" db:"severidad"` // baja, media, alta, critica
	Descripcion       string    `json:"descripcion,omitempty" db:"descripcion"`
	Threshold         float64   `json:"threshold" db:"threshold"`
	IsAnomaly         bool      `json:"is_anomaly" db:"is_anomaly"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// AnomalyInput representa los datos de entrada para crear una anomalía
type AnomalyInput struct {
	ActivoID          int       `json:"activo_id" binding:"required"`
	SensorID          string    `json:"sensor_id" binding:"required"`
	Timestamp         time.Time `json:"timestamp" binding:"required"`
	AnomalyScore      float64   `json:"anomaly_score" binding:"required,min=0,max=1"`
	AnomalyLikelihood float64   `json:"anomaly_likelihood" binding:"required,min=0,max=1"`
	Severidad         string    `json:"severidad" binding:"required,oneof=baja media alta critica"`
	Descripcion       string    `json:"descripcion,omitempty"`
	Threshold         float64   `json:"threshold" binding:"required"`
	IsAnomaly         bool      `json:"is_anomaly"`
}

// SensorData representa los datos recibidos del ParserService
type SensorData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Valor     float64   `json:"valor"`
}

// ParserServiceResponse representa la respuesta del endpoint /lectura/:activo_id/window
type ParserServiceResponse struct {
	ActivoID int `json:"activo_id"`
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Sensores []struct {
		SensorID string       `json:"sensor_id"`
		Tipo     string       `json:"tipo"`
		Lecturas []SensorData `json:"lecturas"`
	} `json:"sensores"`
}
