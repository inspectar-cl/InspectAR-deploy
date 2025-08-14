package database

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Database {
	uri := viper.GetString("mongu.uri")
	dbName := viper.GetString("mongu.database")
	
	// Debug: imprimir lo que está leyendo
	log.Printf("DEBUG: MongoDB URI leído del config: %s", uri)
	log.Printf("DEBUG: MongoDB Database leído del config: %s", dbName)
	
	log.Printf("Intentando conectar a MongoDB en: %s", uri)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Conectado exitosamente a MongoDB, base de datos: %s", dbName)
	return client.Database(dbName)
}
