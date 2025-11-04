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
	usuarioRepo := repository.NewUsuarioRepository(db)
	tipoFallaRepo := repository.NewTipoFallaRepository(db)
	firmaRepo := repository.NewFirmaRepository(db)
	logRepo := repository.NewLogRepository(db)

	// URL de ParserService desde configuración
	parserURL := viper.GetString("services.parser_url")
	if parserURL == "" {
		parserURL = "http://localhost:8095" // Valor por defecto
	}

	// Servicios con las dependencias correctas
	tecnicoService := services.NewTecnicoService(tecnicoRepo)
	accionService := services.NewAccionMantenimientoService(accionRepo)
	activoService := services.NewActivoService(activoRepo, edificioRepo)
	reporteService := services.NewReporteService(reporteRepo, activoRepo, accionRepo, edificioRepo)
	reporteService.SetFirmaRepo(firmaRepo)     // Configurar repositorio de firmas
	reporteService.SetUsuarioRepo(usuarioRepo) // Configurar repositorio de usuarios
	solicitudService := services.NewSolicitudService(db)
	firmaService := services.NewFirmaService(firmaRepo, usuarioRepo)

	// Handlers
	tecnicoHandler := handlers.NewTecnicoHandler(tecnicoService)
	accionHandler := handlers.NewAccionMantenimientoHandler(accionService)
	activoHandler := handlers.NewActivoHandler(activoService)
	reporteHandler := handlers.NewReporteHandler(reporteService)
	solicitudHandler := handlers.NewSolicitudHandler(solicitudService)
	usuarioHandler := handlers.NewUsuarioHandler(usuarioRepo)
	tipoFallaHandler := handlers.NewTipoFallaHandler(tipoFallaRepo)
	firmaHandler := handlers.NewFirmaHandler(firmaService)
	adminHandler := handlers.NewAdminHandler(activoRepo, edificioRepo, usuarioRepo, logRepo, parserURL)

	// Router (ahora con todos los handlers incluido adminHandler)
	r := router.SetupRouter(tecnicoHandler, accionHandler, activoHandler, reporteHandler, solicitudHandler, usuarioHandler, tipoFallaHandler, firmaHandler, adminHandler)

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
