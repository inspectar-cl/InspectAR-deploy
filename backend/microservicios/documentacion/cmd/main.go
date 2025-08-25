package main

import (
	"database/sql"
	"documentacion/api/router"
	"documentacion/config"
	"documentacion/internal/handlers"
	"documentacion/internal/repository"
	"documentacion/internal/services"
	"documentacion/storage"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Conectar a PostgreSQL
	dbConnectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	log.Printf("Conectando a PostgreSQL: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	if err = db.Ping(); err != nil {
		log.Fatalf("Error verificando conexión a PostgreSQL: %v", err)
	}

	log.Println("Conexión a PostgreSQL establecida exitosamente")

	// Inicializar storage service
	var storageService storage.StorageService
	if cfg.Storage.Type == "minio" {
		storageService, err = storage.NewMinioService(&cfg.Storage.Minio)
		if err != nil {
			log.Printf("Error inicializando MinIO, usando storage local: %v", err)
			storageService, err = storage.NewLocalService(&cfg.Storage.Local)
			if err != nil {
				log.Fatalf("Error inicializando storage local: %v", err)
			}
		} else {
			log.Println("MinIO storage inicializado exitosamente")
		}
	} else {
		storageService, err = storage.NewLocalService(&cfg.Storage.Local)
		if err != nil {
			log.Fatalf("Error inicializando storage local: %v", err)
		}
		log.Println("Storage local inicializado exitosamente")
	}

	// Inicializar repositorios
	documentoRepo := repository.NewDocumentoRepository(db)
	analisisRepo := repository.NewAnalisisRepository(db)

	// Inicializar servicios
	documentoService := services.NewDocumentoService(documentoRepo, storageService)

	// Inicializar servicio de IA (opcional)
	var aiService *services.AIService
	if cfg.Gemini.APIKey != "" {
		aiService, err = services.NewAIService(&cfg.Gemini, analisisRepo, documentoRepo, storageService)
		if err != nil {
			log.Printf("Advertencia: No se pudo inicializar servicio de IA: %v", err)
			log.Println("El microservicio funcionará sin análisis de IA")
		} else {
			log.Println("Servicio de IA (Gemini) inicializado exitosamente")
		}
	} else {
		log.Println("API Key de Gemini no configurada. Servicio de IA deshabilitado")
	}

	// Inicializar handlers
	documentoHandler := handlers.NewDocumentoHandler(documentoService, aiService)

	// Configurar router
	r := router.SetupRouter(documentoHandler)

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Iniciando servidor en %s", serverAddr)
	log.Printf("Microservicio de Documentación listo en puerto %s", cfg.Server.Port)
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
