package repository

import (
	"errors"
	"log"
	"apigateway/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password, scope FROM users WHERE email = $1"
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Scope)
	if err != nil {
		log.Printf("Usuario no encontrado %v", err)
		return nil, errors.New("usuario no encontrado")
	}
	log.Println("Usuario encontrado")
	return &user, nil
}
