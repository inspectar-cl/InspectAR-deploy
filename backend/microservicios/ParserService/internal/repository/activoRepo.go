package repository

import (
	"ParserService/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivoRepository struct {
	collection *mongo.Collection
}

func NewActivoRepository(db *mongo.Database) *ActivoRepository {
	return &ActivoRepository{
		collection: db.Collection("iot_db"),
	}
}

func (r *ActivoRepository) CreateActivo(ctx context.Context, activo *models.Activo) (string, error) {
	result, err := r.collection.InsertOne(ctx, activo)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNilDocument
	}
	return id.Hex(), nil
}

func (r *ActivoRepository) GetActivo(ctx context.Context, activoID string) (*models.Activo, error) {
	var activo models.Activo
	err := r.collection.FindOne(ctx, bson.M{"activo_id": activoID}).Decode(&activo)
	if err != nil {
		return nil, err
	}
	return &activo, nil
}

func (r *ActivoRepository) AddSensor(ctx context.Context, activoID string, sensor models.Sensor) error {
	filter := bson.M{"activo_id": activoID}
	update := bson.M{"$push": bson.M{"sensores": sensor}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ActivoRepository) GetAllActivos(ctx context.Context) ([]models.Activo, error) {
	var activos []models.Activo
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var activo models.Activo
		if err := cursor.Decode(&activo); err != nil {
			return nil, err
		}
		activos = append(activos, activo)
	}
	log.Printf("Total activos found: %d", len(activos))
	return activos, nil
}

func (r *ActivoRepository) ExistActivo(ctx context.Context, activoID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"activo_id": activoID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ActivoRepository) ActualizarEstado(ctx context.Context, activoID string, nuevoEstado string) error {
	filter := bson.M{"activo_id": activoID}
	update := bson.M{"$set": bson.M{"estado": nuevoEstado}}

	log.Printf("Actualizando estado del activo %s a '%s'", activoID, nuevoEstado)

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}