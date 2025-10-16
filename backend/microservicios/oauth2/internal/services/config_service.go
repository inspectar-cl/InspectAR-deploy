package services

import (
	"log"
	"oauth2/internal/models"
	"oauth2/internal/repository"
)

type ConfigService struct {
	Repo *repository.ConfigRepository
}

func NewConfigService(repo *repository.ConfigRepository) *ConfigService {
	return &ConfigService{Repo: repo}
}

func (s *ConfigService) LoadConfig() (*models.Config, error) {
	config, err := s.Repo.GetConfig()
	if err != nil {
		return nil, err
	}
	log.Println("Configuración cargada correctamente desde la BD")
	return config, nil
}
