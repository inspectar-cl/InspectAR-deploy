package main

import (
	"gestion/api/router"
	"gestion/internal/database"
	"gestion/internal/handlers"
	"gestion/internal/repository"
	"gestion/internal/services"
	"log"

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

	log.Printf("Archivo de configuración usado: %s", viper.ConfigFileUsed())
}

func main() {
	// Cargar configuración
	initConfig()

	// Conexión a PostgreSQL
	db := database.ConnectPostgres()
	defer db.Close()

	// Repositorios
	tecnicoRepo := repository.NewTecnicoRepository(db)
	accionRepo := repository.NewAccionMantenimientoRepository(db)
	activoRepo := repository.NewActivoRepository(db)
	edificioRepo := repository.NewEdificioRepository(db)
	reporteRepo := repository.NewReporteRepository(db)

	// Servicios con las dependencias correctas
	tecnicoService := services.NewTecnicoService(tecnicoRepo)
	accionService := services.NewAccionMantenimientoService(accionRepo)
	reporteService := services.NewReporteService(reporteRepo, activoRepo, accionRepo, edificioRepo)
	solicitudService := services.NewSolicitudService(db)

	// Handlers
	tecnicoHandler := handlers.NewTecnicoHandler(tecnicoService)
	accionHandler := handlers.NewAccionMantenimientoHandler(accionService)
	reporteHandler := handlers.NewReporteHandler(reporteService)
	solicitudHandler := handlers.NewSolicitudHandler(solicitudService)

	// Router (ahora con solicitudHandler)
	r := router.SetupRouter(tecnicoHandler, accionHandler, reporteHandler, solicitudHandler)

	// Servidor web
	port := viper.GetString("server.port")
	if port == "" {
		port = "8092"
	}

	service := web.NewService(
		web.Name("gestion-service"),
		web.Address(":"+port),
		web.Handler(r),
	)

	log.Printf("Iniciando servidor de gestión en puerto %s", port)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
