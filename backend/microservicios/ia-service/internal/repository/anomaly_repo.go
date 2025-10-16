package repository

import (
	"database/sql"
	"fmt"
	"ia-service/internal/models"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type AnomalyRepository struct {
	db *sql.DB
}

// NewAnomalyRepository crea una nueva instancia del repositorio
func NewAnomalyRepository(connectionString string) (*AnomalyRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión a BD: %w", err)
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando conexión a BD: %w", err)
	}

	log.Println("✅ Conexión exitosa a PostgreSQL (IA-db)")

	return &AnomalyRepository{db: db}, nil
}

// SaveAnomaly guarda una nueva anomalía en la base de datos
func (r *AnomalyRepository) SaveAnomaly(anomaly *models.StoreAnomalyRequest) (*models.Anomaly, error) {
	query := `
		INSERT INTO anomalias (
			activo_id, sensor_id, timestamp, anomaly_score, anomaly_likelihood,
			severidad, descripcion, threshold, is_anomaly, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	var savedAnomaly models.Anomaly
	savedAnomaly.ActivoID = anomaly.ActivoID
	savedAnomaly.SensorID = anomaly.SensorID
	savedAnomaly.Timestamp = anomaly.Timestamp
	savedAnomaly.AnomalyScore = anomaly.AnomalyScore
	savedAnomaly.AnomalyLikelihood = anomaly.AnomalyLikelihood
	savedAnomaly.Severidad = anomaly.Severidad
	savedAnomaly.Descripcion = anomaly.Descripcion
	savedAnomaly.Threshold = anomaly.Threshold
	savedAnomaly.IsAnomaly = anomaly.IsAnomaly

	err := r.db.QueryRow(
		query,
		anomaly.ActivoID,
		anomaly.SensorID,
		anomaly.Timestamp,
		anomaly.AnomalyScore,
		anomaly.AnomalyLikelihood,
		anomaly.Severidad,
		anomaly.Descripcion,
		anomaly.Threshold,
		anomaly.IsAnomaly,
		now,
		now,
	).Scan(&savedAnomaly.ID, &savedAnomaly.CreatedAt, &savedAnomaly.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error guardando anomalía: %w", err)
	}

	log.Printf("✅ Anomalía guardada: ID=%d, Sensor=%s, Score=%.4f, Severidad=%s",
		savedAnomaly.ID, savedAnomaly.SensorID, savedAnomaly.AnomalyScore, savedAnomaly.Severidad)

	return &savedAnomaly, nil
}

// GetLatestAnomaly obtiene la anomalía más reciente
func (r *AnomalyRepository) GetLatestAnomaly() (*models.Anomaly, error) {
	query := `
		SELECT 
			id, activo_id, sensor_id, timestamp, anomaly_score, anomaly_likelihood,
			severidad, descripcion, threshold, is_anomaly, created_at, updated_at
		FROM anomalias
		WHERE is_anomaly = true
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var anomaly models.Anomaly
	err := r.db.QueryRow(query).Scan(
		&anomaly.ID,
		&anomaly.ActivoID,
		&anomaly.SensorID,
		&anomaly.Timestamp,
		&anomaly.AnomalyScore,
		&anomaly.AnomalyLikelihood,
		&anomaly.Severidad,
		&anomaly.Descripcion,
		&anomaly.Threshold,
		&anomaly.IsAnomaly,
		&anomaly.CreatedAt,
		&anomaly.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no se encontraron anomalías")
	}

	if err != nil {
		return nil, fmt.Errorf("error obteniendo última anomalía: %w", err)
	}

	return &anomaly, nil
}

// GetAnomaliesBySensor obtiene todas las anomalías de un sensor específico
func (r *AnomalyRepository) GetAnomaliesBySensor(sensorID string, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT 
			id, activo_id, sensor_id, timestamp, anomaly_score, anomaly_likelihood,
			severidad, descripcion, threshold, is_anomaly, created_at, updated_at
		FROM anomalias
		WHERE sensor_id = $1 AND is_anomaly = 1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, sensorID, limit)
	if err != nil {
		return nil, fmt.Errorf("error consultando anomalías por sensor: %w", err)
	}
	defer rows.Close()

	var anomalies []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		err := rows.Scan(
			&anomaly.ID,
			&anomaly.ActivoID,
			&anomaly.SensorID,
			&anomaly.Timestamp,
			&anomaly.AnomalyScore,
			&anomaly.AnomalyLikelihood,
			&anomaly.Severidad,
			&anomaly.Descripcion,
			&anomaly.Threshold,
			&anomaly.IsAnomaly,
			&anomaly.CreatedAt,
			&anomaly.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando anomalía: %w", err)
		}
		anomalies = append(anomalies, anomaly)
	}

	return anomalies, nil
}

// GetAnomaliesByActivo obtiene todas las anomalías de un activo específico
func (r *AnomalyRepository) GetAnomaliesByActivo(activoID, limit, offset int) ([]models.Anomaly, error) {
	query := `
		SELECT 
			id, activo_id, sensor_id, timestamp, anomaly_score, anomaly_likelihood,
			severidad, descripcion, threshold, is_anomaly, created_at, updated_at
		FROM anomalias
		WHERE activo_id = $1 AND is_anomaly = 1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, activoID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error consultando anomalías por activo: %w", err)
	}
	defer rows.Close()

	var anomalies []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		err := rows.Scan(
			&anomaly.ID,
			&anomaly.ActivoID,
			&anomaly.SensorID,
			&anomaly.Timestamp,
			&anomaly.AnomalyScore,
			&anomaly.AnomalyLikelihood,
			&anomaly.Severidad,
			&anomaly.Descripcion,
			&anomaly.Threshold,
			&anomaly.IsAnomaly,
			&anomaly.CreatedAt,
			&anomaly.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando anomalía: %w", err)
		}
		anomalies = append(anomalies, anomaly)
	}

	return anomalies, nil
}

// Close cierra la conexión a la base de datos
func (r *AnomalyRepository) Close() error {
	return r.db.Close()
}
