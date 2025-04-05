package main

import (
	"ProyectoWEB/proteccion"
	"ProyectoWEB/rutas"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetConnection retorna una conexión singleton a la base de datos
func GetConnection() *sql.DB {
	once.Do(func() {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./golang_consola.db"
		}
		
		var err error
		db, err = sql.Open("sqlite3", dbPath+"?_journal=WAL&_timeout=5000")
		if err != nil {
			log.Fatalf("Error al conectar a la base de datos: %v", err)
		}
		
		// Configurar la conexión
		db.SetMaxOpenConns(1) // SQLite funciona mejor con una sola conexión
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(time.Hour)
	})
	return db
}

func iniciarDB() {
	// Verificar si la base de datos ya existe
	dbPath := os.Getenv("DB_NAME")
	if dbPath == "" {
		dbPath = "./golang_consola.db"
	}
	
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// La base de datos no existe, crearla
		log.Println("Creando nueva base de datos...")
		err := os.WriteFile(dbPath, []byte{}, 0644)
		if err != nil {
			log.Fatalf("Error al crear el archivo de base de datos: %v", err)
		}
		
		// Obtener la conexión
		db := GetConnection()
		
		// Leer el archivo init.sql
		sqlFile, err := os.ReadFile("init.sql")
		if err != nil {
			log.Fatalf("Error al leer el archivo init.sql: %v", err)
		}
		
		// Ejecutar el script SQL
		_, err = db.Exec(string(sqlFile))
		if err != nil {
			log.Fatalf("Error al ejecutar el script SQL: %v", err)
		}
		
		log.Println("Base de datos inicializada correctamente.")
	} else {
		log.Println("La base de datos ya existe, omitiendo inicialización.")
	}
}

func main() {
	// Cargar variables de entorno primero
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	// Inicializar la base de datos si es necesario
	iniciarDB()
	
	mux := mux.NewRouter()

	// Rutas públicas
	mux.HandleFunc("/", rutas.Home)
	mux.HandleFunc("/nosotros", rutas.Nosotros)
	mux.HandleFunc("/seguridad/registro", rutas.Seguridad_registro)
	mux.HandleFunc("/seguridad/login", rutas.Seguridad_login)
	mux.HandleFunc("/seguridad/login_post", rutas.Seguridad_login_post).Methods("POST")
	mux.HandleFunc("/seguridad/registro_post", rutas.Seguridad_registro_post).Methods("POST")

	// Crear un subruteador para rutas protegidas
	protected := mux.NewRoute().Subrouter()

	// Aplicar protección y middleware de no-caché a todas las rutas protegidas
	protected.Use(func(next http.Handler) http.Handler {
		return proteccion.NoCacheMiddleware(next)
	})
	protected.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(proteccion.Proteccion(next.ServeHTTP))
	})

	// Rutas protegidas
	protected.HandleFunc("/perfil", rutas.MiPerfil)
	protected.HandleFunc("/perfil-edit", rutas.EditarPerfilForm).Methods("GET")
	protected.HandleFunc("/perfil-edit-actualizar/{id}", rutas.ActualizarPerfil).Methods("POST")

	protected.HandleFunc("/registro-libros/isbn", rutas.BuscarLibroPorISBN)
	protected.HandleFunc("/registro-libros/isbn_post", rutas.BuscarLibroPorISBN_post).Methods("POST")

	protected.HandleFunc("/registro-libros/urlml", rutas.LibroPorMercadoLibre)
	protected.HandleFunc("/registro-libros/urlml_post", rutas.LibroPorMercadoLibre_post).Methods("POST")

	protected.HandleFunc("/registro-libros/imagen", rutas.LibroImagen)

	protected.HandleFunc("/registro-libros/manual", rutas.Formularios)
	protected.HandleFunc("/formulario-post", rutas.Formulario_post).Methods("POST")
	protected.HandleFunc("/formularios-datos", rutas.Todos_mis_libros)
	protected.HandleFunc("/buscar-libros", rutas.BuscarLibros).Methods("GET")
	protected.HandleFunc("/proximos-libros", rutas.Proximos_Libros)

	protected.HandleFunc("/proximos-libros/{id}", rutas.InformacionLibro)
	protected.HandleFunc("/libro/{id}/editar", rutas.EditarLibroForm).Methods("GET")
	protected.HandleFunc("/libro/{id}/actualizar", rutas.ActualizarLibro).Methods("POST")

	protected.HandleFunc("/seguridad/protegida", rutas.Seguridad_protegida)
	protected.HandleFunc("/seguridad/logout", rutas.Seguridad_logout)

	// Configuración de archivos estaticos hacia mux
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.PathPrefix("/assets/").Handler(s)

	// Obtener puerto del archivo .env
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080" // Puerto por defecto
	}

	server := &http.Server{
		Addr:         "localhost:" + puerto,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	
	fmt.Println("Servidor corriendo en localhost:" + puerto)
	log.Fatal(server.ListenAndServe())
}