package crearDatabase

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase() error {
	// Ruta fija o configurada sin variables de entorno irrelevantes
	dbPath := "./golang_consola.db" 

	// Conectar a SQLite (solo necesita la ruta)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Verificar que el archivo init.sql exista
	sqlFile, err := os.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("error al leer init.sql: %v", err)
	}

	// Ejecutar el script SQL
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error al ejecutar init.sql: %v", err)
	}

	fmt.Println("Base de datos inicializada correctamente.")
	return nil
}