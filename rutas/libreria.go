package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func InformacionLibro(w http.ResponseWriter, r *http.Request) {
	conectar.Conectar()
	defer conectar.CerrarConexion()

	vars := mux.Vars(r)
	idLibroStr := vars["id"]

	// Convertir idLibro a entero para evitar inyección SQL
	idLibro, err := strconv.Atoi(idLibroStr)
	if err != nil {
		http.Error(w, "ID de libro no válido", http.StatusBadRequest)
		return
	}

	// Preparar la consulta SQL
	stmt, err := conectar.Db.Prepare("SELECT id, nombre, autor, editorial, Descripcion, portada, puntaje FROM todosloslibros WHERE id = ?")
	if err != nil {
		fmt.Println("Error al preparar la consulta:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var libro modelos.Libro
	err = stmt.QueryRow(idLibro).Scan(&libro.Id, &libro.Nombre, &libro.Autor, &libro.Editorial, &libro.Descripcion, &libro.Url_imagen, &libro.Puntaje)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			fmt.Println("Error al consultar el libro:", err)
			http.Error(w, "Error al obtener datos del libro", http.StatusInternalServerError)
		}
		return
	}

	// Cargar la plantilla HTML
	tmpl := template.Must(template.ParseFiles(
		"templates/libreria/detalles_libro.html",
		"templates/principales/navbarLogin.html",
		"templates/principales/navbarNotLogin.html",
		utilidades.Frontend,
	))

	// Pasar datos a la plantilla
	data := DatosDeSeguridad(w, r)
	data["Datos"] = libro

	err = tmpl.ExecuteTemplate(w, "detalles_libro.html", data)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
	}
}

func Formularios(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/principales/Formulario.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
	data := DatosDeSeguridad(response, request)
	template.ExecuteTemplate(response, "Formulario.html", data)
}
func Todos_mis_libros(response http.ResponseWriter, request *http.Request) {
	// Cargar la plantilla HTML
	tmpl := template.Must(template.ParseFiles("templates/principales/FormularioDatos.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))

	// Conectar a la base de datos
	conectar.Conectar()
	defer conectar.CerrarConexion() // Asegurar cierre de conexión

	sql := "SELECT UsuarioID,nombre, url_imagen, autor, isbn, editorial, descripcion, critica, leido, puntaje FROM libros WHERE UsuarioID = ? ORDER BY nombre DESC"
	// Definir slice para almacenar los resultados
	var libros []modelos.Libro
	usuarioID := DatosDeSeguridad(response, request)["usuario_id"]
	// Ejecutar la consulta
	datos, err := conectar.Db.Query(sql, usuarioID)
	if err != nil {
		http.Error(response, "Error en la consulta a la base de datos", http.StatusInternalServerError)
		fmt.Println("Error en la consulta SQL:", err)
		return
	}
	defer datos.Close() // Cerrar después de obtener los datos

	// Recorrer los resultados
	for datos.Next() {
		var dato modelos.Libro
		err := datos.Scan(&dato.UsuarioID, &dato.Nombre, &dato.Url_imagen, &dato.Autor, &dato.Isbn, &dato.Editorial, &dato.Descripcion, &dato.Critica, &dato.Leido, &dato.Puntaje)
		if err != nil {
			fmt.Println("Error al escanear los datos:", err)
			continue // Saltar este registro si hay un error
		}
		libros = append(libros, dato)
	}
	// fmt.Printf("Libros obtenidos: %+v\n", libros)

	// Verificar si hubo errores en la iteración
	if err = datos.Err(); err != nil {
		fmt.Println("Error al iterar sobre los datos:", err)
	}
	data := DatosDeSeguridad(response, request)
	data["Datos"] = libros

	// Renderizar la plantilla
	err = tmpl.ExecuteTemplate(response, "FormularioDatos.html", data)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
	}

}
func BuscarLibros(response http.ResponseWriter, request *http.Request) {
    // Obtener el término de búsqueda del query string
    q := request.URL.Query().Get("q")
	likePattern := "%" + q + "%"
	fmt.Printf("Término: [%s], Patrón: [%s]\n", q, likePattern)
    fmt.Println("Término de búsqueda bruto:", q)
    
    // Cargar la plantilla HTML
    tmpl := template.Must(template.ParseFiles("templates/principales/FormularioDatos.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
    
    // Conectar a la base de datos
    conectar.Conectar()
    defer conectar.CerrarConexion()
    
    // Obtener datos de seguridad (incluyendo usuario_id)
    data := DatosDeSeguridad(response, request)
    usuarioID := data["usuario_id"]
    
    var libros []modelos.Libro    
    
    // Evitar usar fmt.Sprintf para construir el patrón
    sql := "SELECT UsuarioID, nombre, url_imagen, autor, isbn, editorial, descripcion, critica, leido, puntaje FROM libros WHERE UsuarioID = ? AND (nombre LIKE ? OR autor LIKE ?)"
    
    // Debug manual para ver exactamente qué se está enviando a la base de datos
    fmt.Println("SQL query:", sql)
    fmt.Println("Usuario ID:", usuarioID)
    fmt.Println("Patrón LIKE:", likePattern)

	searchPattern := "%Los%"
    
    // Ejecutar la consulta
    datos, err := conectar.Db.Query(sql, usuarioID, searchPattern, searchPattern)
    if err != nil {
        fmt.Println("Error SQL:", err)
        http.Error(response, "Error en la base de datos", http.StatusInternalServerError)
        return
    }
    defer datos.Close()
    
    // Escanear resultados
    for datos.Next() {
        var libro modelos.Libro
        err := datos.Scan(
            &libro.UsuarioID, 
            &libro.Nombre, 
            &libro.Url_imagen, 
            &libro.Autor, 
            &libro.Isbn, 
            &libro.Editorial, 
            &libro.Descripcion, 
            &libro.Critica, 
            &libro.Leido, 
            &libro.Puntaje)
        
        if err != nil {
            fmt.Println("Error al escanear:", err)
            continue
        }
        
        libros = append(libros, libro)
    }
    
    // Mostrar resultados para debugging
    fmt.Println("Libros encontrados:", len(libros))
    
    // Asignar resultados a la plantilla
    data["Datos"] = libros
    
    // Renderizar plantilla
    err = tmpl.ExecuteTemplate(response, "FormularioDatos.html", data)
    if err != nil {
        fmt.Println("Error de plantilla:", err)
    }
}
func Formulario_post(response http.ResponseWriter, request *http.Request) {
	conectar.Conectar()

	sql := `INSERT INTO todosloslibros (UsuarioId, nombre, portada, autor, isbn, editorial, descripcion, Comentarios, leido, fecha_creacion, fecha_lectura, puntaje)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?);`

	puntajeStr := request.FormValue("puntaje")
	var puntaje interface{}
	if puntajeStr == "" {
		puntaje = nil
	} else {
		p, err := strconv.Atoi(puntajeStr)
		if err != nil || p < 0 || p > 10 {
			http.Error(response, "Puntaje inválido. Debe estar entre 0 y 10.", http.StatusBadRequest)
			return
		}
		puntaje = p
	}

	isbn := request.FormValue("isbn")
	var isbnDB interface{}
	if isbn == "" {
		isbnDB = "SIN_ISBN"
	} else {
		isbnDB = isbn
	}

	leido := request.FormValue("leido")
	var fecha_lectura time.Time
	var leidoDB int
	if leido == "si" {
		leidoDB = 1
		fecha_lectura = time.Now()

	} else {
		leidoDB = 0
		fecha_lectura = time.Time{}
	}

	usuarioID := DatosDeSeguridad(response, request)["usuario_id"]

	_, err := conectar.Db.Exec(sql,
		usuarioID,
		request.FormValue("nombre"),      //ok
		request.FormValue("url"),         //ok
		request.FormValue("autor"),       //ok
		isbnDB,                           //ok
		request.FormValue("editorial"),   //ok
		request.FormValue("descripcion"), //ok
		request.FormValue("comentarios"), //ok
		leidoDB,                          //ok
		time.Now(),                       // fecha_creacion//ok
		fecha_lectura,
		puntaje, //ok
	)

	if err != nil {
		fmt.Fprintln(response, err)
	}

	utilidades.CrearMensajesFlash(response, request, "success", "Se creó el registro exitosamente")
	http.Redirect(response, request, "/formularios-datos", http.StatusSeeOther)
}
func Proximos_Libros(response http.ResponseWriter, request *http.Request) {
    // Obtener el ID del usuario logueado desde la sesión
    session, _ := utilidades.Store.Get(request, "session-name")
    usuarioID, ok := session.Values["usuario_id"].(string)
    if !ok || usuarioID == "" {
        http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
        return
    }

    // Crear un FuncMap con las funciones personalizadas
    funcMap := template.FuncMap{
        "contarLibrosPorAutor": func(autor string) int {
            var cantidad int
            err := conectar.Db.QueryRow(`
                SELECT COUNT(*) 
                FROM todosloslibros 
                WHERE autor = ? AND UsuarioId = ?
            `, autor, usuarioID).Scan(&cantidad)
            if err != nil {
                fmt.Println("Error al contar libros por autor:", err)
                return 0
            }
            return cantidad
        },
        "buscarMasLibros": func(autor string) []struct {
            Id     int
            Nombre string
        } {
            var libros []struct {
                Id     int
                Nombre string
            }
            rows, err := conectar.Db.Query(`
                SELECT id, nombre 
                FROM todosloslibros 
                WHERE autor = ? AND UsuarioId = ? 
                LIMIT 5
            `, autor, usuarioID)
            if err == nil {
                defer rows.Close()
                for rows.Next() {
                    var libro struct {
                        Id     int
                        Nombre string
                    }
                    rows.Scan(&libro.Id, &libro.Nombre)
                    libros = append(libros, libro)
                }
            }
            return libros
        },
    }

    // Cargar la plantilla HTML con las funciones personalizadas
    tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
        "templates/libreria/LibrosQueQuiero.html",
        "templates/principales/navbarLogin.html",
        "templates/principales/navbarNotLogin.html",
        utilidades.Frontend))

    // Conectar a la base de datos
    conectar.Conectar()
    defer conectar.CerrarConexion()

    // Consulta SQL modificada para filtrar por usuario_id
    sql := `
        SELECT id, nombre, portada, isbn, editorial, comentarios, autor 
        FROM todosloslibros 
        WHERE UsuarioId = ?
        ORDER BY fecha_creacion DESC
    `

    var libros []modelos.Libro
    datos, err := conectar.Db.Query(sql, usuarioID)
    if err != nil {
        http.Error(response, "Error en la consulta a la base de datos", http.StatusInternalServerError)
        fmt.Println("Error en la consulta SQL:", err)
        return
    }
    defer datos.Close()

    for datos.Next() {
        var dato modelos.Libro
        err := datos.Scan(&dato.Id, &dato.Nombre, &dato.Url_imagen, &dato.Isbn, &dato.Editorial, &dato.Descripcion, &dato.Autor)
        if err != nil {
            fmt.Println("Error al escanear los datos:", err)
            continue
        }
        libros = append(libros, dato)
    }

    if err = datos.Err(); err != nil {
        fmt.Println("Error al iterar sobre los datos:", err)
    }

    data := DatosDeSeguridad(response, request)
    data["Datos"] = libros

    err = tmpl.ExecuteTemplate(response, "LibrosQueQuiero.html", data)
    if err != nil {
        fmt.Println("Error al ejecutar la plantilla:", err)
    }
}
func EditarLibroForm(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "ID de libro inválido", http.StatusBadRequest)
        return
    }

    libro, err := obtenerLibroPorID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Cargar la plantilla HTML
    tmpl := template.Must(template.ParseFiles(
        "templates/libreria/editar.html",
        "templates/principales/navbarLogin.html",
        "templates/principales/navbarNotLogin.html",
        utilidades.Frontend))

    // Pasar datos a la plantilla
    data := DatosDeSeguridad(w, r)
    data["Datos"] = libro

    // Corregido: usar el nombre correcto de la plantilla
    err = tmpl.ExecuteTemplate(w, "editar.html", data)
    if err != nil {
        fmt.Println("Error al ejecutar la plantilla:", err)
        http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
    }
}
func obtenerLibroPorID(id int) (*modelos.Libro, error) {
    conectar.Conectar()
    defer conectar.CerrarConexion()

    // Corregido según la estructura real de la tabla
    query := `SELECT id, UsuarioID, nombre, portada, autor, isbn, 
             editorial, descripcion, comentarios, leido, puntaje 
             FROM todosloslibros WHERE id = ?`

    var libro modelos.Libro
    err := conectar.Db.QueryRow(query, id).Scan(
        &libro.Id, &libro.UsuarioID, &libro.Nombre, &libro.Url_imagen, 
        &libro.Autor, &libro.Isbn, &libro.Editorial,
        &libro.Comentarios, &libro.Critica, &libro.Leido, &libro.Puntaje,
    )
    if err != nil {
        return nil, fmt.Errorf("error al obtener libro: %w", err)
    }

    return &libro, nil
}
// actualizarLibro procesa el formulario de edición y actualiza el libro
func ActualizarLibro(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "ID de libro inválido", http.StatusBadRequest)
        return
    }

    // Recopilar datos del formulario
    isbn := r.FormValue("isbn")
    nombre := r.FormValue("nombre")
    autor := r.FormValue("autor")
    editorial := r.FormValue("editorial")
    url_imagen := r.FormValue("url_imagen")
    critica := r.FormValue("critica")
    descripcion := r.FormValue("descripcion")

    // En tu DB el campo leido es enum('si','no')
    leido := r.FormValue("leido")
    if leido == "1" {
        leido = "si"
    } else {
        leido = "no"
    }

    // Convertir campos numéricos
    puntaje, _ := strconv.Atoi(r.FormValue("puntaje"))

    // Actualizar libro en la base de datos
    err = actualizarLibroEnDB(id, isbn, nombre, autor, editorial,
        url_imagen, critica, descripcion, leido, puntaje)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Redirigir a la página de detalles del libro
    utilidades.CrearMensajesFlash(w, r, "success", "Libro actualizado exitosamente")
    http.Redirect(w, r, fmt.Sprintf("/proximos-libros/%d", id), http.StatusSeeOther)
}
// actualizarLibroEnDB actualiza los datos de un libro en la base de datos
func actualizarLibroEnDB(id int, isbn, nombre, autor, editorial,
    portada, critica, descripcion string, leido string,
    puntaje int) error {

    conectar.Conectar()
    defer conectar.CerrarConexion()

    // Corregido según la estructura real de la tabla
    query := `UPDATE todosloslibros 
             SET isbn = ?, nombre = ?, autor = ?, editorial = ?, 
             portada = ?, comentarios = ?, descripcion = ?, 
             leido = ?, puntaje = ? 
             WHERE id = ?`

    // Preparar los punteros para campos que pueden ser NULL
    var isbnPtr *string
    if isbn != "" {
        isbnPtr = &isbn
    }

    _, err := conectar.Db.Exec(query, 
        isbnPtr, nombre, autor, editorial,
        portada, critica, descripcion, 
        leido, puntaje, id)

    if err != nil {
        return fmt.Errorf("error al actualizar libro: %w", err)
    }

    return nil
}




