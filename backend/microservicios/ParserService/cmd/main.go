package main

import (
	"crypto/tls"
	"forms/internal/database"
	"forms/internal/handler"
	"forms/internal/repository"
	"forms/internal/services"
	"log"

	"github.com/go-micro/plugins/v4/registry/consul"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
	"go-micro.dev/v4/registry"
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
	// Cargar configuración
	initConfig()

	// Conexión a MongoDB
	mongoDB := database.ConnectMongo()

	// Conexión a InfluxDB
	influxURL := viper.GetString("influx.url")
	influxToken := viper.GetString("influx.token")
	influxOrg := viper.GetString("influx.org")
	influxBucket := viper.GetString("influx.bucket")

	influxClient := influxdb2.NewClientWithOptions(
		influxURL,
		influxToken,
		influxdb2.DefaultOptions().SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)

	// Repositorios
	activoRepo := repository.NewActivoRepository(mongoDB)
	influxRepo := repository.NewInfluxRepository(influxClient, influxOrg, influxBucket)

	// Servicios
	activoService := services.NewActivoService(activoRepo)
	sensorService := services.NewSensorService(influxRepo)

	// Handler de datos
	dataHandler := handler.NewDataHandler(activoService, sensorService)

	// Registro en Consul
	reg := consul.NewRegistry(registry.Addrs("consul:8500"))

	// Ruteo con Gin
	r := router.SetupRouter(dataHandler)

	// Servicio web
	service := web.NewService(
		web.Name("iot-service"),
		web.Registry(reg),
		web.Address(":8090"),
		web.Handler(r),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
