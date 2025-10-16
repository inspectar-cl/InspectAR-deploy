package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Gemini   GeminiConfig   `mapstructure:"gemini"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Limits   LimitsConfig   `mapstructure:"limits"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type StorageConfig struct {
	Type  string      `mapstructure:"type"`
	Minio MinioConfig `mapstructure:"minio"`
	Local LocalConfig `mapstructure:"local"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Secure    bool   `mapstructure:"secure"`
}

type LocalConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
}

type GeminiConfig struct {
	APIKey    string `mapstructure:"api_key"`
	ProjectID string `mapstructure:"project_id"`
	Location  string `mapstructure:"location"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type LimitsConfig struct {
	MaxFileSizeMB     int      `mapstructure:"max_file_size_mb"`
	AllowedExtensions []string `mapstructure:"allowed_extensions"`
}

func LoadConfig() (*Config, error) {
	// Obtener el nombre del archivo de configuración desde variable de entorno
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config"
	} else {
		// Si se especifica un archivo completo, remover la extensión
		if strings.HasSuffix(configFile, ".yaml") || strings.HasSuffix(configFile, ".yml") {
			configFile = strings.TrimSuffix(configFile, ".yaml")
			configFile = strings.TrimSuffix(configFile, ".yml")
		}
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Variables de entorno
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DOC")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error leyendo archivo de configuración: %v", err)
		return nil, err
	}

	log.Printf("Archivo de configuración usado: %s", viper.ConfigFileUsed())

	// Expandir variables de entorno en valores específicos
	viper.Set("gemini.api_key", os.ExpandEnv(viper.GetString("gemini.api_key")))
	viper.Set("gemini.project_id", os.ExpandEnv(viper.GetString("gemini.project_id")))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
