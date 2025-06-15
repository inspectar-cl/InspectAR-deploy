package repository

import (
	"context"
	"fmt"
	"time"

	"forms/internal/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	influxdb2write "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxRepository struct {
	writeAPI influxdb2api.WriteAPIBlocking
	queryAPI influxdb2api.QueryAPI
	bucket   string
}

func NewInfluxRepository(client influxdb2.Client, org, bucket string) *InfluxRepository {
	return &InfluxRepository{
		writeAPI: client.WriteAPIBlocking(org, bucket),
		queryAPI: client.QueryAPI(org),
		bucket:   bucket,
	}
}

// Guarda un dato de sensor en InfluxDB
func (r *InfluxRepository) InsertSensorData(ctx context.Context, sensorID string, valor float64, timestamp time.Time) error {
	point := influxdb2write.NewPoint(
		"mediciones",
		map[string]string{
			"sensor_id": sensorID,
		},
		map[string]interface{}{
			"valor": valor,
		},
		timestamp,
	)

	return r.writeAPI.WritePoint(ctx, point)
}

// Obtiene todos los datos históricos de un sensor
func (r *InfluxRepository) GetSensorData(ctx context.Context, sensorID string, since time.Duration) ([]models.SensorData, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: -%dm)
  |> filter(fn: (r) => r._measurement == "mediciones" and r.sensor_id == "%s")
  |> keep(columns: ["_time", "_value"])
`, r.bucket, int(since.Minutes()), sensorID)

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var datos []models.SensorData
	for result.Next() {
		datos = append(datos, models.SensorData{
			Tiempo: result.Record().Time().Format(time.RFC3339),
			Valor:  result.Record().Value().(float64),
		})
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return datos, nil
}
