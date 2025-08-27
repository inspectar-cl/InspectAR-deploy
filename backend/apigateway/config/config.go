package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	GestionURL string
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
		GestionURL:  getenv("GESTION_URL", "http://localhost:8092"), // TODO: Quitar lo segundo en producción (?)
	}
}
