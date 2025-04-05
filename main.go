package main

import (
	"ProyectoWEB/proteccion"
	"ProyectoWEB/rutas"
	"fmt"
	"log"
	"net/http"
	"os"
	"ProyectoWEB/crearDatabase"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)
func iniciarDB(){
	dsn := "usuario:contraseña@tcp(127.0.0.1:3306)/"
	err := crearDatabase.InitializeDatabase(dsn)
	if err != nil {
        fmt.Println("Error al inicializar la base de datos:", err)
        return
    }

    fmt.Println("Aplicación lista para usar.")
}

func main() {
	/* iniciar esta función por primera vez unicamente.*/
	// iniciarDB()
	mux := mux.NewRouter()

	// Rutas públicas
	mux.HandleFunc("/", rutas.Home)
	mux.HandleFunc("/nosotros", rutas.Nosotros)
	mux.HandleFunc("/seguridad/registro", rutas.Seguridad_registro)
	mux.HandleFunc("/seguridad/login", rutas.Seguridad_login)
	mux.HandleFunc("/seguridad/login_post", rutas.Seguridad_login_post).Methods("POST")
	mux.HandleFunc("/seguridad/registro_post", rutas.Seguridad_registro_post).Methods("POST")

	// Crear un subruteador para rutas protegidas
	/* te permite organizar y agrupar rutas relacionadas bajo un prefijo o configuración común.
	Es una forma de subdividir tu enrutador principal para aplicar configuraciones específicas a un grupo de rutas. */
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

	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))) // importante el punto antes del /assets/
	mux.PathPrefix("/assets/").Handler(s)

	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
		return
	}
	server := &http.Server{
		Addr:    "localhost:" + os.Getenv("PORT"),
		Handler: mux,
		// Se recomiendan agregar estas dos variables. Por que? no lo se rick
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Corriendo servidor dede localhost:" + os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())

}
