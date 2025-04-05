package conectar

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // Import para SQLite
)

// Db es la variable global para manejar la conexión a la base de datos.
var Db *sql.DB

// Conectar inicializa la conexión a la base de datos.
func Conectar() {
	// Cargar las variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando las variables de entorno: %v", err)
	}

	// Obtener el tipo de base de datos desde las variables de entorno
	dbType := os.Getenv("DB_TYPE")
	var conection *sql.DB
	var err error

	switch dbType {
	case "mysql":
		// Configuración para MySQL
		dsn := fmt.Sprintf("%s:@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_SERVER"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		conection, err = sql.Open("mysql", dsn)
	case "sqlite3":
		// Configuración para SQLite
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			log.Fatal("La ruta del archivo SQLite no está definida en las variables de entorno (DB_PATH).")
		}
		conection, err = sql.Open("sqlite3", dbPath)
	default:
		log.Fatalf("Tipo de base de datos no soportado: %s", dbType)
	}

	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	if err = conection.Ping(); err != nil {
		log.Fatalf("No se pudo establecer conexión con la base de datos: %v", err)
	}

	log.Println("Conexión exitosa a la base de datos.")
	Db = conection
}

// CerrarConexion cierra la conexión a la base de datos.
func CerrarConexion() {
	if Db != nil {
		if err := Db.Close(); err != nil {
			log.Printf("Error al cerrar la conexión: %v", err)
		} else {
			log.Println("Conexión cerrada exitosamente.")
		}
	}
}
