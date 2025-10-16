package repository

import (
	"apigateway/internal/models"

	"github.com/jmoiron/sqlx"
)

type ConfigRepository struct {
	db *sqlx.DB
}

func NewConfigRepository(db *sqlx.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) GetConfig() (*models.Config, error) {
	var config models.Config
	query := "SELECT key_access, key_refresh, lifetime_access, lifetime_refresh FROM config LIMIT 1"
	err := r.db.QueryRow(query).Scan(&config.KeyAccess, &config.KeyRefresh, &config.LifetimeAccess, &config.LifetimeRefresh)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
