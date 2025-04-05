package conectar

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // Solo SQLite
)

var Db *sql.DB

func Conectar() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	if os.Getenv("DB_TYPE") != "sqlite3" {
		log.Fatal("Solo se soporta SQLite3. Verifica tu variable DB_TYPE en .env")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH no está definido en .env")
	}

	var err error
	Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}

	if err := Db.Ping(); err != nil {
		log.Fatalf("Error en Ping: %v", err)
	}

	log.Println("✅ Conexión exitosa a SQLite3")
}

func CerrarConexion() {
	if Db != nil {
		Db.Close()
	}
}