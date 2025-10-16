package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func connectPostgres() *sql.DB {
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

func main() {
	fmt.Println("🔄 EJECUTANDO MIGRACIÓN DE BASE DE DATOS")
	fmt.Println("========================================")
	fmt.Println("Agregando columnas de observaciones editables...")

	// Cargar configuración
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Configurar valores por defecto
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "password")
	viper.SetDefault("postgres.dbname", "gestion_db")
	viper.SetDefault("postgres.sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("⚠️  No se pudo leer el archivo de configuración: %v\n", err)
		fmt.Println("Usando valores por defecto...")
	}

	// Conectar a la base de datos
	fmt.Println("🔍 Conectando a PostgreSQL...")
	db := connectPostgres()
	defer db.Close()

	// Leer archivo de migración
	migrationPath := filepath.Join("migrations", "001_add_observaciones_columns.sql")
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		// Intentar ruta alternativa
		migrationPath = filepath.Join("..", "migrations", "001_add_observaciones_columns.sql")
		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			log.Fatalf("❌ Error: No se encontró el archivo de migración")
		}
	}

	fmt.Printf("📋 Leyendo migración: %s\n", migrationPath)
	migrationSQL, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("❌ Error leyendo archivo de migración: %v", err)
	}

	// Ejecutar migración
	fmt.Println("🚀 Ejecutando migración...")
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatalf("❌ Error ejecutando migración: %v", err)
	}

	fmt.Println("✅ MIGRACIÓN COMPLETADA EXITOSAMENTE")

	// Verificar estructura de la tabla
	fmt.Println("\n📊 Verificando estructura de la tabla reportes...")
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable, column_default 
		FROM information_schema.columns 
		WHERE table_name = 'reportes' 
		AND column_name IN (
			'observaciones_analista', 'autor_analista', 'estructura_informe', 
			'metadata_informe', 'version_reporte', 'estado_revision', 
			'fecha_revision', 'revisor', 'observaciones_revision'
		)
		ORDER BY ordinal_position;
	`)
	if err != nil {
		log.Printf("⚠️  Error verificando estructura: %v", err)
	} else {
		fmt.Println("\n🎯 COLUMNAS AGREGADAS:")
		for rows.Next() {
			var columnName, dataType, isNullable string
			var columnDefault sql.NullString

			err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)
			if err != nil {
				continue
			}

			defaultValue := "NULL"
			if columnDefault.Valid {
				defaultValue = columnDefault.String
			}

			fmt.Printf("   ✓ %s (%s, nullable: %s, default: %s)\n",
				columnName, dataType, isNullable, defaultValue)
		}
		rows.Close()
	}

	fmt.Println("\n🎉 Ya puedes usar las nuevas rutas de observaciones editables!")
	fmt.Println("\n📝 Para probar las nuevas funcionalidades ejecuta:")
	fmt.Println("   ./quick_test.sh")
	fmt.Println("   ./test_all_routes.sh")
}
