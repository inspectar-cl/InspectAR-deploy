package models

type Config struct {
	ID              int    `json:"primaryKey"`
	KeyAccess       string `json:"column:key_access"`
	KeyRefresh      string `json:"column:key_refresh"`
	LifetimeAccess  int    `json:"column:lifetime_access"`
	LifetimeRefresh int    `json:"column:lifetime_refresh"`
}
