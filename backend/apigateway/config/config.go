package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       	string
	FrontendURL string
	GestionURL 	string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:        getenv("GATEWAY_PORT", "3500"),
		FrontendURL: getenv("FRONTEND_URL", "http://localhost:3000"),
		GestionURL:  getenv("GESTION_URL", "http://localhost:8092"), // Lo segundo es la segunda ruta que verificará.
	}
}
