package repository

import (
	"ParserService/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorStatusRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSensorStatusRepository(db *mongo.Database) *SensorStatusRepository {
	return &SensorStatusRepository{
		db:         db,
		collection: db.Collection("sensor_status"),
	}
}

// UpdateSensorActivity actualiza la actividad de un sensor
func (r *SensorStatusRepository) UpdateSensorActivity(sensorID string, timestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"sensor_id": sensorID}

	// Verificar si el sensor ya existe
	var existingSensor models.SensorStatus
	err := r.collection.FindOne(ctx, filter).Decode(&existingSensor)

	if err == mongo.ErrNoDocuments {
		// Crear nuevo sensor
		newSensor := models.SensorStatus{
			SensorID:     sensorID,
			IsActive:     true,
			LastSeen:     timestamp,
			FirstSeen:    timestamp,
			TotalReports: 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err = r.collection.InsertOne(ctx, newSensor)
		if err != nil {
			log.Printf("Error al crear nuevo sensor %s: %v", sensorID, err)
			return err
		}

		log.Printf("✅ Nuevo sensor registrado: %s", sensorID)
		return nil
	} else if err != nil {
		log.Printf("Error al buscar sensor %s: %v", sensorID, err)
		return err
	}

	// Actualizar sensor existente
	update := bson.M{
		"$set": bson.M{
			"is_active":  true,
			"last_seen":  timestamp,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{
			"total_reports": 1,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error al actualizar sensor %s: %v", sensorID, err)
		return err
	}

	return nil
}

// GetInactiveSensors obtiene sensores que no han reportado en los últimos 5 minutos
func (r *SensorStatusRepository) GetInactiveSensors(timeoutMinutes int) ([]models.SensorStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cutoffTime := time.Now().Add(-time.Duration(timeoutMinutes) * time.Minute)

	filter := bson.M{
		"is_active": true,
		"last_seen": bson.M{"$lt": cutoffTime},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error al buscar sensores inactivos: %v", err)
	}
	defer cursor.Close(ctx)

	var inactiveSensors []models.SensorStatus
	for cursor.Next(ctx) {
		var sensor models.SensorStatus
		if err := cursor.Decode(&sensor); err != nil {
			log.Printf("Error al decodificar sensor: %v", err)
			continue
		}
		inactiveSensors = append(inactiveSensors, sensor)
	}

	return inactiveSensors, nil
}

// MarkSensorAsInactive marca un sensor como inactivo
func (r *SensorStatusRepository) MarkSensorAsInactive(sensorID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"sensor_id": sensorID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error al marcar sensor %s como inactivo: %v", sensorID, err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("sensor %s no encontrado o ya estaba inactivo", sensorID)
	}

	log.Printf("🔴 Sensor marcado como inactivo: %s", sensorID)
	return nil
}

// GetAllSensors obtiene todos los sensores con su estado
func (r *SensorStatusRepository) GetAllSensors() ([]models.SensorStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ordenar por última actividad
	opts := options.Find().SetSort(bson.D{{Key: "last_seen", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error al obtener sensores: %v", err)
	}
	defer cursor.Close(ctx)

	var sensors []models.SensorStatus
	for cursor.Next(ctx) {
		var sensor models.SensorStatus
		if err := cursor.Decode(&sensor); err != nil {
			log.Printf("Error al decodificar sensor: %v", err)
			continue
		}
		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

// GetSensorStatus obtiene el estado de un sensor específico
func (r *SensorStatusRepository) GetSensorStatus(sensorID string) (*models.SensorStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sensor models.SensorStatus
	err := r.collection.FindOne(ctx, bson.M{"sensor_id": sensorID}).Decode(&sensor)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("sensor %s no encontrado", sensorID)
	} else if err != nil {
		return nil, fmt.Errorf("error al obtener sensor %s: %v", sensorID, err)
	}

	return &sensor, nil
}
