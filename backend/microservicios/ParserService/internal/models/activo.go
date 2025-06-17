package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sensor struct {
	SensorID string `bson:"sensor_id" json:"sensor_id"`
	Tipo    string `bson:"tipo" json:"tipo"`
	Unidad  string `bson:"unidad" json:"unidad"`
}

type Activo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActivoID  string             `bson:"activo_id" json:"activo_id"`
	Nombre    string             `bson:"nombre" json:"nombre"`
	Ubicacion string             `bson:"ubicacion" json:"ubicacion"`
	Estado 	  string			 `bson:"estado" json:"estado"`
	EdificioID string			 `bson:"id_edificio" json:"id_edificio"`
	Sensores  []Sensor           `bson:"sensores" json:"sensores"`
}
