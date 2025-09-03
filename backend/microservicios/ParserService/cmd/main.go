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

	// Debug: mostrar qué archivo de config se está usando
	log.Printf("DEBUG: Archivo de configuración usado: %s", viper.ConfigFileUsed())

	// Debug: mostrar valores clave
	log.Printf("DEBUG: mongu.uri = %s", viper.GetString("mongu.uri"))
	log.Printf("DEBUG: influxdb.url = %s", viper.GetString("influxdb.url"))
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
	sensorStatusRepo := repository.NewSensorStatusRepository(mongoDB)

	// Servicios
	activoService := services.NewActivoService(activoRepo)
	sensorService := services.NewSensorService(influxRepo)

	// Configuración del servicio de monitoreo
	notificationURL := viper.GetString("notification.url") // Opcional, por ahora será console.log
	timeoutMinutes := viper.GetInt("sensor.timeout_minutes")
	if timeoutMinutes == 0 {
		timeoutMinutes = 5 // Default: 5 minutos
	}

	monitoringService := services.NewSensorMonitoringService(
		sensorStatusRepo,
		influxRepo,
		notificationURL,
		timeoutMinutes,
	)

	// Handlers
	dataHandler := handlers.NewDataHandler(activoService, sensorService, monitoringService)
	statusHandler := handlers.NewSensorStatusHandler(monitoringService)

	// Iniciar monitoreo automático (verificar cada 2 minutos)
	log.Printf("🚀 Iniciando monitoreo automático de sensores...")
	monitoringService.StartMonitoring(2 * time.Minute)

	// Ruteo con Gin
	r := router.SetupRouter(dataHandler, statusHandler)

	// Servicio web
	service := web.NewService(
		web.Name("iot-service"),
		web.Address(":8090"),
		web.Handler(r),
	)

	log.Printf("🎯 Servidor iniciado en puerto 8090")
	log.Printf("📊 Endpoints de monitoreo disponibles:")
	log.Printf("   GET /api/sensors/status - Estado de todos los sensores")
	log.Printf("   GET /api/sensors/status/:sensor_id - Estado de un sensor")
	log.Printf("   GET /api/sensors/stats - Estadísticas generales")
	log.Printf("   POST /api/sensors/check-disconnected - Verificación manual")
	log.Printf("   GET /api/sensors/health - Health check del monitoreo")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
