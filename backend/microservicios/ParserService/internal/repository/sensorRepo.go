package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ParserService/internal/models"

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

	fmt.Printf("Insertando en InfluxDB: sensor_id=%s, valor=%f, timestamp=%s\n", sensorID, valor, timestamp.Format(time.RFC3339))

	return r.writeAPI.WritePoint(ctx, point)
}

// Obtiene todos los datos históricos de un sensor
// func (r *InfluxRepository) GetSensorData(ctx context.Context, sensorID string, since time.Duration) ([]models.SensorData, error) {
// 	query := fmt.Sprintf(`
// from(bucket: "%s")
//   |> range(start: -%dm)
//   |> filter(fn: (r) => r._measurement == "mediciones" and r.sensor_id == "%s")
//   |> keep(columns: ["_time", "_value"])
// `, r.bucket, int(since.Minutes()), sensorID)

// 	result, err := r.queryAPI.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var datos []models.SensorData
// 	for result.Next() {
// 		datos = append(datos, models.SensorData{
// 			Tiempo: result.Record().Time().Format(time.RFC3339),
// 			Valor:  result.Record().Value().(float64),
// 		})
// 	}
// 	if result.Err() != nil {
// 		return nil, result.Err()
// 	}
// 	return datos, nil
// }

// Obtiene todos los datos históricos de un sensor sin filtro de tiempo
func (r *InfluxRepository) GetSensorData(ctx context.Context, sensorID string) ([]models.SensorData, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> filter(fn: (r) => r._measurement == "mediciones" and r.sensor_id == "%s")
  |> keep(columns: ["_time", "_value"])
`, r.bucket, sensorID)

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

func (r *InfluxRepository) GetSensorLastData(ctx context.Context, sensorID string) (*models.SensorData, error) {
	log.Printf("Obteniendo último dato de sensor %s", sensorID)
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> filter(fn: (r) => r._measurement == "mediciones" and r.sensor_id == "%s")
  |> keep(columns: ["_time", "_value"])
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 1)
`, r.bucket, sensorID)

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if result.Next() {
		return &models.SensorData{
			Tiempo: result.Record().Time().Format(time.RFC3339),
			Valor:  result.Record().Value().(float64),
		}, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return nil, nil // No hay datos
}
