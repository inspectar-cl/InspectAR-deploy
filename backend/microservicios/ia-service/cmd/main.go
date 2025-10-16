package main

import (
	"fmt"
	"ia-service/api"
	"ia-service/internal/handler"
	"ia-service/internal/repository"
	"ia-service/internal/services"
	"log"
	"os"

	"github.com/spf13/viper"
)

func initConfig() {
	// Obtener el nombre del archivo desde variable de entorno o usar "config" por defecto
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config"
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ Error al leer config: %v", err)
	}

	log.Printf("✅ Archivo de configuración usado: %s", viper.ConfigFileUsed())
}

func main() {
	log.Println("🚀 Iniciando IA-Service...")

	// 1. Cargar configuración
	initConfig()

	// 2. Obtener valores de configuración
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")
	dbSSLMode := viper.GetString("database.sslmode")

	parserServiceURL := viper.GetString("parser_service.url")
	mlEngineURL := viper.GetString("ml_engine.url")
	serverPort := viper.GetString("server.port")
	fetchInterval := viper.GetInt("scheduler.fetch_interval_minutes")
	defaultActivoID := viper.GetInt("scheduler.default_activo_id")

	// 3. Construir cadena de conexión a BD
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	log.Printf("✅ Configuración cargada: Puerto=%s, DB=%s@%s:%s, ParserURL=%s, MLEngine=%s",
		serverPort, dbName, dbHost, dbPort, parserServiceURL, mlEngineURL)

	// 4. Conectar a la base de datos
	repo, err := repository.NewAnomalyRepository(connString)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}
	defer repo.Close()

	// 5. Crear servicio
	anomalyService := services.NewAnomalyService(repo, parserServiceURL, mlEngineURL, fetchInterval, defaultActivoID)

	// 6. Iniciar scheduler en background
	anomalyService.StartScheduler()

	// 7. Crear handler
	anomalyHandler := handler.NewAnomalyHandler(anomalyService)

	// 8. Configurar router
	router := api.SetupRouter(anomalyHandler)

	// 9. Iniciar servidor
	serverAddr := ":" + serverPort
	log.Printf("✅ IA-Service escuchando en el puerto %s", serverPort)
	log.Println("📋 Rutas disponibles:")
	log.Println("   POST   /anomalies/store")
	log.Println("   GET    /status/")
	log.Println("   GET    /anomalies/sensor/:sensor_id")
	log.Println("   GET    /anomalies/activo/:activo_id")
	log.Println("   GET    /health")

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
