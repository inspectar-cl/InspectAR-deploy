package database

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
)

func ConnectInflux() influxdb2.Client {
	uri := viper.GetString("influxdb.url")
	token := viper.GetString("influxdb.token")
	if uri == "" {
		uri = "http://localhost:8086" // valor por defecto
	}
	if token == "" {
		log.Fatal("No se encontró el token de InfluxDB en la configuración")
	}
	client := influxdb2.NewClient(uri, token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Probar conexión
	_, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("Error conectando a InfluxDB: %v", err)
	}
	return client
}
