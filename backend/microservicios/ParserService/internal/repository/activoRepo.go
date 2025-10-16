package repository

import (
	"ParserService/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivoRepository struct {
	primary    *mongo.Collection
	sources    []activoSource
	sensorColl *mongo.Collection
}

type activoSource struct {
	coll   *mongo.Collection
	fields map[string]string // mapeo de campos origen->modelo (activo_id, nombre, estado, sensores, id_edificio)
}

func NewActivoRepository(db *mongo.Database) *ActivoRepository {
	// Leer configuración para soportar múltiples colecciones
	primaryName := viper.GetString("activos.primary_collection")
	if primaryName == "" {
		primaryName = "activos"
	}

	// Definir fuentes desde config: activos.sources: [{ name: "activos", fields: { ... } }]
	var cfgSources []map[string]interface{}
	_ = viper.UnmarshalKey("activos.sources", &cfgSources)

	var sources []activoSource
	if len(cfgSources) > 0 {
		for _, s := range cfgSources {
			name, _ := s["name"].(string)
			if name == "" {
				continue
			}
			// fields opcional
			fields := defaultFieldMap()
			if fm, ok := s["fields"].(map[string]interface{}); ok {
				for k, v := range fm {
					if vs, ok := v.(string); ok && vs != "" {
						fields[k] = vs
					}
				}
			}
			sources = append(sources, activoSource{coll: db.Collection(name), fields: fields})
		}
	} else {
		// Alternativa: lista simple de nombres
		var names []string
		_ = viper.UnmarshalKey("activos.source_names", &names)
		if len(names) == 0 {
			names = []string{primaryName}
		}
		for _, name := range names {
			if name == "" {
				continue
			}
			sources = append(sources, activoSource{coll: db.Collection(name), fields: defaultFieldMap()})
		}
	}

	// Colección separada de sensores
	sensorsName := viper.GetString("activos.sensors_collection")
	if sensorsName == "" {
		sensorsName = "sensores"
	}

	return &ActivoRepository{
		primary:    db.Collection(primaryName),
		sources:    sources,
		sensorColl: db.Collection(sensorsName),
	}
}

func (r *ActivoRepository) CreateActivo(ctx context.Context, activo *models.Activo) (string, error) {
	_, err := r.primary.InsertOne(ctx, activo)
	if err != nil {
		return "", err
	}
	// Devolver el activo_id en lugar del ObjectID de MongoDB
	return fmt.Sprint(activo.ActivoID), nil
}

func (r *ActivoRepository) GetActivo(ctx context.Context, activoID int) (*models.Activo, error) {
	// Buscar en todas las fuentes según mapeo de campos
	for _, src := range r.sources {
		filter := bson.M{src.fields["activo_id"]: activoID}
		var raw bson.M
		err := src.coll.FindOne(ctx, filter).Decode(&raw)
		if err == mongo.ErrNoDocuments {
			continue
		}
		if err != nil {
			// Loguear y seguir buscando en otras fuentes
			log.Printf("WARN: error buscando activo %s en %s: %v", activoID, src.coll.Name(), err)
			continue
		}
		act := mapToActivo(raw, src.fields)
		return &act, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (r *ActivoRepository) AddSensor(ctx context.Context, activoID int, sensor models.Sensor) error {
	// Preferir colección separada de sensores si existe
	if r.sensorColl != nil {
		doc := bson.M{
			"sensor_id": sensor.SensorID,
			"activo_id": activoID,
			"tipo":      sensor.Tipo,
			"unidad":    sensor.Unidad,
		}
		if _, err := r.sensorColl.InsertOne(ctx, doc); err == nil {
			return nil
		} else {
			log.Printf("WARN: fallo insert en colección sensores: %v", err)
		}
	}

	// Fallback: embebido en documento de activo
	filter := bson.M{"activo_id": activoID}
	update := bson.M{"$push": bson.M{"sensores": sensor}}
	if res, err := r.primary.UpdateOne(ctx, filter, update); err == nil && res.ModifiedCount > 0 {
		return nil
	}

	// Intentar en otras fuentes
	for _, src := range r.sources {
		if src.coll.Name() == r.primary.Name() {
			continue
		}
		f := bson.M{src.fields["activo_id"]: activoID}
		u := bson.M{"$push": bson.M{src.fields["sensores"]: sensor}}
		if res, err := src.coll.UpdateOne(ctx, f, u); err == nil && res.ModifiedCount > 0 {
			return nil
		}
	}
	return mongo.ErrNoDocuments
}

func (r *ActivoRepository) GetAllActivos(ctx context.Context) ([]models.Activo, error) {
	var activos []models.Activo
	seen := map[string]bool{}

	for _, src := range r.sources {
		cursor, err := src.coll.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("WARN: error listando en %s: %v", src.coll.Name(), err)
			continue
		}
		func() {
			defer cursor.Close(ctx)
			for cursor.Next(ctx) {
				var raw bson.M
				if err := cursor.Decode(&raw); err != nil {
					log.Printf("WARN: error decodificando doc en %s: %v", src.coll.Name(), err)
					continue
				}
				act := mapToActivo(raw, src.fields)
				// Evitar duplicados por ActivoID
				if act.ActivoID == 0 || seen[fmt.Sprint(act.ActivoID)] {
					continue
				}
				seen[fmt.Sprint(act.ActivoID)] = true
				activos = append(activos, act)
			}
		}()
	}

	log.Printf("Total activos found (merge): %d", len(activos))
	return activos, nil
}

func (r *ActivoRepository) ExistActivo(ctx context.Context, activoID int) (bool, error) {
	// Verificar en todas las fuentes para evitar duplicados
	for _, src := range r.sources {
		count, err := src.coll.CountDocuments(ctx, bson.M{src.fields["activo_id"]: activoID})
		if err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r *ActivoRepository) ActualizarEstado(ctx context.Context, activoID int, nuevoEstado string) error {
	// Intentar en la colección primaria
	filter := bson.M{"activo_id": activoID}
	update := bson.M{"$set": bson.M{"estado": nuevoEstado}}

	log.Printf("Actualizando estado del activo %s a '%s'", activoID, nuevoEstado)

	if res, err := r.primary.UpdateOne(ctx, filter, update); err == nil && res.ModifiedCount > 0 {
		return nil
	}

	// Intentar en el resto de fuentes usando el mapeo de campos
	for _, src := range r.sources {
		if src.coll.Name() == r.primary.Name() {
			continue
		}
		f := bson.M{src.fields["activo_id"]: activoID}
		u := bson.M{"$set": bson.M{src.fields["estado"]: nuevoEstado}}
		if res, err := src.coll.UpdateOne(ctx, f, u); err == nil && res.ModifiedCount > 0 {
			return nil
		}
	}
	return mongo.ErrNoDocuments
}

// GetActivosByEdificio obtiene todos los activos de un edificio específico
func (r *ActivoRepository) GetActivosByEdificio(ctx context.Context, edificioID int) ([]models.Activo, error) {
	var activos []models.Activo
	seen := map[string]bool{}

	log.Printf("DEBUG: Buscando activos para edificio_id=%d", edificioID)

	for _, src := range r.sources {
		// Buscar por edificio_id
		// En MongoDB el campo se llama "edificio_id" y está guardado como int
		filter := bson.M{src.fields["edificio_id"]: edificioID}

		log.Printf("DEBUG: Buscando en colección %s con filtro: %v (edificio_id=%d)", src.coll.Name(), filter, edificioID)

		cursor, err := src.coll.Find(ctx, filter)
		if err != nil {
			log.Printf("WARN: error buscando activos por edificio en %s: %v", src.coll.Name(), err)
			continue
		}

		func() {
			defer cursor.Close(ctx)
			for cursor.Next(ctx) {
				var raw bson.M
				if err := cursor.Decode(&raw); err != nil {
					log.Printf("WARN: error decodificando activo en %s: %v", src.coll.Name(), err)
					continue
				}
				act := mapToActivo(raw, src.fields)
				// Evitar duplicados por ActivoID
				if act.ActivoID == 0 || seen[fmt.Sprint(act.ActivoID)] {
					continue
				}
				seen[fmt.Sprint(act.ActivoID)] = true
				activos = append(activos, act)
			}
		}()
	}

	log.Printf("Total activos encontrados para edificio %d: %d", edificioID, len(activos))
	return activos, nil
}

// GetSensoresByActivo obtiene todos los sensores de un activo desde la colección de sensores
func (r *ActivoRepository) GetSensoresByActivo(ctx context.Context, activoID int) ([]models.Sensor, error) {
	if r.sensorColl == nil {
		return []models.Sensor{}, nil
	}

	cursor, err := r.sensorColl.Find(ctx, bson.M{"activo_id": activoID})
	if err != nil {
		log.Printf("ERROR: No se pudieron obtener sensores para activo %d: %v", activoID, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var sensores []models.Sensor
	for cursor.Next(ctx) {
		var sm bson.M
		if err := cursor.Decode(&sm); err != nil {
			log.Printf("WARN: error decodificando sensor: %v", err)
			continue
		}

		sensor := models.Sensor{
			SensorID: asString(sm["sensor_id"]),
			Tipo:     asString(sm["tipo"]),
			Unidad:   asString(sm["unidad"]),
		}

		if sensor.SensorID != "" {
			sensores = append(sensores, sensor)
		}
	}

	log.Printf("Sensores encontrados para activo %d: %d", activoID, len(sensores))
	return sensores, nil
}

// --- Helpers ---

func defaultFieldMap() map[string]string {
	return map[string]string{
		"activo_id":   "activo_id",
		"estado":      "estado",
		"edificio_id": "edificio_id",
	}
}

func asString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case primitive.ObjectID:
		return t.Hex()
	case int32:
		return fmt.Sprint(int(t))
	case int64:
		return fmt.Sprint(t)
	case float64:
		return fmt.Sprint(int64(t))
	default:
		return ""
	}
}

func asInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float64:
		return int(t)
	case string:
		// intentar extraer dígitos si viene con prefijo tipo "AC-1001"
		var num int
		_, err := fmt.Sscanf(t, "%d", &num)
		if err == nil {
			return num
		}
		// extraer últimos números
		var d int
		for i := 0; i < len(t); i++ {
			if t[i] >= '0' && t[i] <= '9' {
				// intentar parsear desde primer dígito
				_, err2 := fmt.Sscanf(t[i:], "%d", &d)
				if err2 == nil {
					return d
				}
				break
			}
		}
		return 0
	default:
		return 0
	}
}

func mapToActivo(raw bson.M, fm map[string]string) models.Activo {
	var id primitive.ObjectID
	if rid, ok := raw["_id"].(primitive.ObjectID); ok {
		id = rid
	}

	activo := models.Activo{
		ID:         id,
		ActivoID:   asInt(raw[fm["activo_id"]]),
		Estado:     asString(raw[fm["estado"]]),
		EdificioID: asInt(raw[fm["edificio_id"]]),
	}

	return activo
}

// enrichSensors ya no es necesaria, se elimina
func (r *ActivoRepository) enrichSensors(ctx context.Context, act *models.Activo) {
	// Los activos ya no contienen sensores
	return
}
