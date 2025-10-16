package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectPostgres() *sql.DB {
	host := viper.GetString("postgres.host")
	port := viper.GetInt("postgres.port")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")
	dbname := viper.GetString("postgres.dbname")
	sslmode := viper.GetString("postgres.sslmode")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	log.Printf("Conectando a PostgreSQL: %s:%d/%s", host, port, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error abriendo conexión a PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error conectando a PostgreSQL: %v", err)
	}

	log.Printf("Conectado exitosamente a PostgreSQL")
	return db
}
