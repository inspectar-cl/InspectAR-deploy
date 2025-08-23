package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	FrontURL    string
	SensoresURL string
	ActivoURL   string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	_ = godotenv.Load() // opcional en prod

	return &Config{
		Port:        getenv("GATEWAY_PORT", "3500"),
		FrontURL:    getenv("FRONT_URL", "http://localhost:3000"),
		SensoresURL: getenv("SENSORES_URL", "http://localhost:4000"),
		ActivoURL:   getenv("ACTIVO_URL", "http://localhost:8090"),
	}
}
