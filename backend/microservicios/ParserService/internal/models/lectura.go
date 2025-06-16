package models

type LecturaRequest struct {
	ActivoID string `json:"activo_id"`
	SensorID string `json:"sensor_id"`
	Tipo     string `json:"tipo"`
	Unidad   string `json:"unidad"`
}

type SensorDataset struct {
	SensorID string       `json:"sensor_id"`
	Datos    []SensorData `json:"datos"`
}

type LecturaSensorRequest struct {
	SensorID  string  `json:"sensor_id"`
	Valor     float64 `json:"valor"`
	Timestamp string  `json:"timestamp"`
}
