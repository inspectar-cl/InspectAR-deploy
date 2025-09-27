package main

import (
	// "crypto/tls"
	"log"
	"oauth2/api/router"
	"oauth2/internal/database"
	"oauth2/internal/handler"
	"oauth2/internal/repository"
	"oauth2/internal/services"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go-micro.dev/v4/web"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config") // Ruta dentro del contenedor
	// viper.AddConfigPath("/app/config") // Alternativa
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer config: %v", err)
	}
}

func main() {
	time.Sleep(3 * time.Second)
	// Inicializar la configuración
	initConfig()

	// Conectar a PostgreSQL
	db := database.InitDB()
	defer db.Close()

	// Inicializar dependencias
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	configRepo := repository.NewConfigRepository(db)
	configService := services.NewConfigService(configRepo)

	// Cargar configuración
	config, err := configService.LoadConfig()
	if err != nil {
		log.Fatal("Error al cargar configuración:", err)
	}

	authService := services.NewAuthService(userRepo, tokenRepo, config)
	loginHandler := handler.NewLoginHandler(authService)

	// Configurar rutas con Gin
	r := router.SetupRouter(loginHandler)

	service := web.NewService(
		web.Name("oauth2-service"),
		web.Address(":8094"),
		web.Handler(r),
		// web.TLSConfig(tlsConfig),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
