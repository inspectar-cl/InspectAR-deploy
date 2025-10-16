package repository

import (
	"database/sql" //Ver esta libreria ("sql.ErrNoRows")
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	DB *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

// Guarda un refresh token en la bd
func (r *TokenRepository) SaveRefreshToken(userID uint, deviceID string, token string, expiresAt time.Time) error {
	_, err := r.DB.Exec("INSERT INTO refresh_tokens (user_id, device_id, token, expires_at) VALUES ($1, $2, $3, $4)", userID, deviceID, token, expiresAt)
	return err
}

// Obtiene un refresh token de la bd por user_id y device_id
func (r *TokenRepository) GetRefreshToken(userID uint, deviceID string) (string, error) {
	var token string
	var expiresAt time.Time

	err := r.DB.QueryRow(
		"SELECT token, expires_at FROM refresh_tokens WHERE user_id = $1 AND device_id = $2",
		userID,
		deviceID,
	).Scan(&token, &expiresAt)

	if err != nil {
		log.Printf("Error al buscar token para usuario %d y dispositivo %s: %v", userID, deviceID, err)
		return "", err
	}

	if time.Now().After(expiresAt) {
		// Si el token está expirado, lo eliminamos y retornamos error
		r.DeleteAllTokensForUser(userID)
		log.Printf("Token expirado para usuario %d y dispositivo %s", userID, deviceID)
		return "", sql.ErrNoRows
	}

	return token, nil
}

// Cuenta la cantidad de tokens de refresh de un usuario
func (r *TokenRepository) GetCountTokens(userID uint) (int, error) {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Si no hay registros, retorna 0
		}
		log.Printf("Error al contar tokens para usuario %d: %v", userID, err)
		return 0, err
	}
	return count, nil
}

// Elimina un refresh token (por cierre de sesion en dispositivo)
func (r *TokenRepository) DeleteRefreshToken(userID uint, deviceID string) error {
	_, err := r.DB.Exec("DELETE FROM refresh_tokens WHERE user_id = $1 AND device_id = $2", userID, deviceID)
	if err != nil {
		log.Printf("Error al eliminar refresh token: %v", err)
	}
	return err
}

// Elimina el refresh token más antiguo de un usuario basado en created_at
func (r *TokenRepository) DeleteRefreshTokenOld(userID uint) error {
	// Obtener el ID del token más antiguo
	var tokenID int
	err := r.DB.QueryRow(`
        SELECT id 
        FROM refresh_tokens 
        WHERE user_id = $1 
        ORDER BY created_at ASC 
        LIMIT 1`, userID).Scan(&tokenID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No hay tokens para eliminar
		}
		return err
	}

	// Eliminar el token más antiguo
	_, err = r.DB.Exec("DELETE FROM refresh_tokens WHERE id = $1", tokenID)
	return err
}

// Elimina todos los refresh tokens de un usuario (puede ser por cierre de sesion en todos los dispositivos)
func (r *TokenRepository) DeleteAllTokensForUser(userID uint) error {
	_, err := r.DB.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", userID)
	return err
}
