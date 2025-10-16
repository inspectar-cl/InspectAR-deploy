package models

import (
	"time"
)

// SensorStatus representa el estado de un sensor en el sistema
type SensorStatus struct {
	SensorID     string    `json:"sensor_id" bson:"sensor_id"`
	IsActive     bool      `json:"is_active" bson:"is_active"`
	LastSeen     time.Time `json:"last_seen" bson:"last_seen"`
	FirstSeen    time.Time `json:"first_seen" bson:"first_seen"`
	TotalReports int64     `json:"total_reports" bson:"total_reports"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// SensorStatusUpdate representa los datos para actualizar el estado de un sensor
type SensorStatusUpdate struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
}

// DisconnectedSensor representa un sensor que se ha desconectado
type DisconnectedSensor struct {
	SensorID       string    `json:"sensor_id"`
	LastSeen       time.Time `json:"last_seen"`
	DisconnectedAt time.Time `json:"disconnected_at"`
	TotalReports   int64     `json:"total_reports"`
}
