package models

import "time"

// Anomaly representa una anomalía detectada en el sistema
type Anomaly struct {
	ID                int       `json:"id" db:"id"`
	ActivoID          int       `json:"activo_id" db:"activo_id"`
	SensorID          string    `json:"sensor_id" db:"sensor_id"`
	Timestamp         time.Time `json:"timestamp" db:"timestamp"`
	AnomalyScore      float64   `json:"anomaly_score" db:"anomaly_score"`
	AnomalyLikelihood float64   `json:"anomaly_likelihood" db:"anomaly_likelihood"`
	Severidad         string    `json:"severidad" db:"severidad"`
	Descripcion       string    `json:"descripcion,omitempty" db:"descripcion"`
	Threshold         float64   `json:"threshold" db:"threshold"`
	IsAnomaly         int       `json:"is_anomaly" db:"is_anomaly"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// StoreAnomalyRequest representa la solicitud para almacenar una anomalía
type StoreAnomalyRequest struct {
	ActivoID          int       `json:"activo_id" binding:"required"`
	SensorID          string    `json:"sensor_id" binding:"required"`
	Timestamp         time.Time `json:"timestamp" binding:"required"`
	AnomalyScore      float64   `json:"anomaly_score"`
	AnomalyLikelihood float64   `json:"anomaly_likelihood"`
	Severidad         string    `json:"severidad"`
	Descripcion       string    `json:"descripcion"`
	Threshold         float64   `json:"threshold"`
	IsAnomaly         int       `json:"is_anomaly"`
}

// SensorData representa los datos de un sensor desde ParserService
type SensorData struct {
	IDSensor string      `json:"id_sensor"`
	Datos    []DataPoint `json:"datos"`
}

// DataPoint representa un punto de dato desde ParserService
type DataPoint struct {
	Tiempo string  `json:"tiempo"` // El ParserService envía "tiempo" como string
	Valor  float64 `json:"valor"`
}

// ParserServiceResponse representa la respuesta del endpoint /lectura/:activo_id/window
type ParserServiceResponse struct {
	ActivoID int          `json:"activo_id"`
	Page     int          `json:"page"`
	Limit    int          `json:"limit"`
	Offset   int          `json:"offset"`
	Sensores []SensorData `json:"sensores"`
}

// MLSensorRecord representa un registro de sensor para el ML Engine
// con todos los campos de sensores en un único timestamp
type MLSensorRecord struct {
	Timestamp      string  `json:"timestamp"`
	AACRMotPV      float64 `json:"A_ACR_Mot.PV"`
	AACRMotSV      float64 `json:"A_ACR_Mot.SV"`
	ADescBomPV     float64 `json:"A_Desc_Bom.PV"`
	ADescBomSV     float64 `json:"A_Desc_Bom.SV"`
	FlujoCaudal    float64 `json:"Flujo_Caudal"`
	NivelPresion   float64 `json:"Nivel_Presion"`
	PotenciaActiva float64 `json:"Potencia_Activa"`
	TempAgua       float64 `json:"Temp_Agua"`
	TempAmbiente   float64 `json:"Temp_Ambiente"`
	Vibracion      float64 `json:"Vibracion"`
}

// MLPredictionRequest representa la solicitud al ML Engine
type MLPredictionRequest struct {
	Pump    string           `json:"pump"`
	Records []MLSensorRecord `json:"records"`
}

// MLAnomalyResult representa el resultado de una predicción individual
type MLAnomalyResult struct {
	Timestamp         string  `json:"timestamp"`
	AnomalyScore      float64 `json:"anomaly_score"`
	AnomalyLikelihood float64 `json:"anomaly_likelihood"`
	Threshold         float64 `json:"threshold"`
	IsAnomaly         int     `json:"is_anomaly"`
	Severity          string  `json:"severity"`
	Description       string  `json:"description"`
}

// MLPredictionResponse representa la respuesta del ML Engine
type MLPredictionResponse struct {
	Pump    string            `json:"pump"`
	Results []MLAnomalyResult `json:"results"`
	Message string            `json:"message,omitempty"`
}
