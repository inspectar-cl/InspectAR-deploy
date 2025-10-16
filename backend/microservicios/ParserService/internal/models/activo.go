package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sensor struct {
	SensorID string `bson:"sensor_id" json:"sensor_id"`
	Tipo     string `bson:"tipo" json:"tipo"`
	Unidad   string `bson:"unidad" json:"unidad"`
}

type Activo struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActivoID   int                `bson:"activo_id" json:"activo_id"`
	Estado     string             `bson:"estado" json:"estado"`
	EdificioID int                `bson:"edificio_id" json:"edificio_id"`
}
