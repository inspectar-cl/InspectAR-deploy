package repository

import (
	"context"
	"forms/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FormRepository struct {
	collection *mongo.Collection
}

func NewFormRepository(db *mongo.Database) *FormRepository {
	return &FormRepository{
		collection: db.Collection("forms"),
	}
}

func (r *FormRepository) CreateForm(ctx context.Context, form *models.Form) (string, error) {
	result, err := r.collection.InsertOne(ctx, form)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNilDocument
	}
	return id.Hex(), nil
}

func (r *FormRepository) GetFormsByUser(ctx context.Context, user string) ([]models.Form, error) {
	var forms []models.Form
	cursor, err := r.collection.Find(ctx, bson.M{"user": user})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var form models.Form
		if err := cursor.Decode(&form); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}

// Pienso que es mejor obtener formularios especificos sin necesidad de el usuario, ya que igualmente estamos verificando la identidad de este con los tokens
// Esta funcion obtiene un formulario de id especifico de MONGODB segun un usuario
func (r *FormRepository) GetFormByUserAndID(ctx context.Context, user string, id string) (*models.Form, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var form models.Form
	filter := bson.M{"_id": objectID, "user": user}
	err = r.collection.FindOne(ctx, filter).Decode(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

// Esta funcion obtiene un formulario de id especifico de la estructura JSON segun un usuario
func (r *FormRepository) GetFormByUserAndCustomID(ctx context.Context, user string, formID string) (*models.Form, error) {
	filter := bson.M{
		"user":                   user,
		"data.formulario.nombre": formID,
	}

	var form models.Form
	err := r.collection.FindOne(ctx, filter).Decode(&form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *FormRepository) GetFormByIdForm(ctx context.Context, ids []string) ([]models.Form, error) {
	var forms []models.Form
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Si el id no es válido, lo ignora
		}
		var form models.Form
		err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&form)
		if err == nil {
			forms = append(forms, form)
		}
	}
	return forms, nil
}

type ActivoRepository struct {
	collection *mongo.Collection
}

func NewActivoRepository(db *mongo.Database) *ActivoRepository {
	return &ActivoRepository{
		collection: db.Collection("activos"),
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
