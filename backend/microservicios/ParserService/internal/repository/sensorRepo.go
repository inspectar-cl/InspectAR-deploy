package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ParserService/internal/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
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
	point := influxdb2.NewPointWithMeasurement("sensor_reading").
		AddTag("sensor_id", sensorID).
		AddField("valor", valor).
		SetTime(timestamp)

	fmt.Printf("Insertando en InfluxDB: sensor_id=%s, valor=%f, timestamp=%s\n", sensorID, valor, timestamp.Format(time.RFC3339))

	return r.writeAPI.WritePoint(ctx, point)
}

// Obtiene todos los datos históricos de un sensor
// func (r *InfluxRepository) GetSensorData(ctx context.Context, sensorID string, since time.Duration) ([]models.SensorData, error) {
// 	query := fmt.Sprintf(`
// from(bucket: "%s")
//   |> range(start: -%dm)
//   |> filter(fn: (r) => r._measurement == "sensor_reading" and r.sensor_id == "%s")
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

// Obtiene los últimos 100 datos de un sensor
func (r *InfluxRepository) GetSensorData(ctx context.Context, sensorID string) ([]models.SensorData, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> filter(fn: (r) => r["_measurement"] == "sensor_reading" and r["sensor_id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "valor")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 100)
`, r.bucket, sensorID)

	log.Printf("Ejecutando query para sensor %s: %s", sensorID, query)

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		log.Printf("Error en query para sensor %s: %v", sensorID, err)
		return nil, err
	}

	var datos []models.SensorData
	count := 0
	for result.Next() {
		count++
		datos = append(datos, models.SensorData{
			Tiempo: result.Record().Time().Format(time.RFC3339),
			Valor:  result.Record().Value().(float64),
		})
	}
	if result.Err() != nil {
		log.Printf("Error iterando resultados para sensor %s: %v", sensorID, result.Err())
		return nil, result.Err()
	}

	log.Printf("Datos obtenidos para sensor %s: %d registros", sensorID, count)
	return datos, nil
}

// Obtiene los últimos N datos de un sensor con paginación (configurable)
func (r *InfluxRepository) GetSensorDataWindow(ctx context.Context, sensorID string, limit int, offset int) ([]models.SensorData, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> filter(fn: (r) => r["_measurement"] == "sensor_reading" and r["sensor_id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "valor")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: %d, offset: %d)
`, r.bucket, sensorID, limit, offset)

	log.Printf("Ejecutando query window para sensor %s con límite %d y offset %d", sensorID, limit, offset)

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		log.Printf("Error en query window para sensor %s: %v", sensorID, err)
		return nil, err
	}

	var datos []models.SensorData
	count := 0
	for result.Next() {
		count++
		datos = append(datos, models.SensorData{
			Tiempo: result.Record().Time().Format(time.RFC3339),
			Valor:  result.Record().Value().(float64),
		})
	}
	if result.Err() != nil {
		log.Printf("Error iterando resultados window para sensor %s: %v", sensorID, result.Err())
		return nil, result.Err()
	}

	log.Printf("Datos window obtenidos para sensor %s: %d registros (offset: %d)", sensorID, count, offset)
	return datos, nil
}

func (r *InfluxRepository) GetSensorLastData(ctx context.Context, sensorID string) (*models.SensorData, error) {
	log.Printf("Obteniendo último dato de sensor %s", sensorID)
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: 0)
  |> filter(fn: (r) => r["_measurement"] == "sensor_reading" and r["sensor_id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "valor")
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
