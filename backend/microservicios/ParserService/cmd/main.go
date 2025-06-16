package main

import (
	"ParserService/api/router"
	"ParserService/internal/database"
	handlers "ParserService/internal/handler"
	"ParserService/internal/repository"
	"ParserService/internal/services"
	"log"
	"time"

	"github.com/spf13/viper"
	"go-micro.dev/v4/web"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer config: %v", err)
	}
}

func main() {
	time.Sleep(5 * time.Second) // Simulación de espera para asegurar que la configuración se cargue correctamente
	// Cargar configuración
	initConfig()

	// Conexión a MongoDB
	mongoDB := database.ConnectMongo()

	// Conexión a InfluxDB usando el helper
	influxClient := database.ConnectInflux()
	influxOrg := viper.GetString("influxdb.org")
	influxBucket := viper.GetString("influxdb.bucket")

	// Repositorios
	activoRepo := repository.NewActivoRepository(mongoDB)
	influxRepo := repository.NewInfluxRepository(influxClient, influxOrg, influxBucket)

	// Servicios
	activoService := services.NewActivoService(activoRepo)
	sensorService := services.NewSensorService(influxRepo)

	// Handler de datos
	dataHandler := handlers.NewDataHandler(activoService, sensorService)

	// Ruteo con Gin
	r := router.SetupRouter(dataHandler)

	// Servicio web
	service := web.NewService(
		web.Name("iot-service"),
		web.Address(":8090"),
		web.Handler(r),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
