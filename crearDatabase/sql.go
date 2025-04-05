package crearDatabase

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDatabase(dsn string) error {
	// Conectar a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Leer el archivo init.sql
	sqlFile, err := ioutil.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("error al leer el archivo init.sql: %v", err)
	}

	// Ejecutar el script SQL
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error al ejecutar el script SQL: %v", err)
	}

	fmt.Println("Base de datos inicializada correctamente.")
	return nil
}
