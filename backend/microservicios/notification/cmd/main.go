package main

import (
	"log"
	"notification/api/router"
	"notification/internal/database"
	handlers "notification/internal/handler"
	"notification/internal/repository"
	"notification/internal/services"
	"time"

	"github.com/spf13/viper"
	"go-micro.dev/v4/web"
)

func initConfig() {
	// Configurar viper para leer variables de entorno automáticamente
	viper.AutomaticEnv()

	// Intentar leer config.yaml solo como fallback
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No se pudo leer config.yaml, usando solo variables de entorno: %v", err)
	}

	// Verificar si tenemos variables de entorno de correo configuradas
	if viper.GetString("EMAIL_SMTP_HOST") != "" {
		log.Printf("✅ Usando configuración de correo desde variables de entorno")
	} else if viper.GetString("email.smtp_host") != "" {
		log.Printf("⚠️ Usando configuración de correo desde config.yaml")
	} else {
		log.Printf("❌ No se encontró configuración de correo")
	}
}

func main() {
	time.Sleep(10 * time.Second) // Esperar a que la base de datos esté lista
	// Cargar configuración
	initConfig()

	// Conexión a PostgreSQL
	db := database.ConnectPostgres()

	// Repositorios
	notificationRepo := repository.NewNotificationRepository(db)

	// Servicios
	emailService := services.NewEmailService()
	notificationService := services.NewNotificationService(notificationRepo, emailService)

	// Handler de notificaciones
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Ruteo con Gin
	r := router.SetupRouter(notificationHandler)

	// Servicio web
	service := web.NewService(
		web.Name("notification-service"),
		web.Address(":8091"),
		web.Handler(r),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
