package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

//Renderizamos la página principal cargar_isbn.html
func BuscarLibroPorISBN(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("templates/libreria/cargar_ISBN.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
	data := DatosDeSeguridad(w, r)
	template.ExecuteTemplate(w, "cargar_ISBN.html", data)
}
// FetchBookInfo obtiene la información de un libro desde OpenLibrary
func fetchBookInfo(isbn string) (*modelos.BookInfo, error) {
	// Asegurarse de que el ISBN no tiene espacios o guiones
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	// URL de la API con formato específico
	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", isbn)

	// Realizar la solicitud HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error de conexión: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error de API: código de estado %d", resp.StatusCode)
	}

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}

	// La respuesta viene en formato {"ISBN:XXXXXXXXXX": {datos}}
	var result map[string]modelos.BookInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %w", err)
	}

	// Comprobar si el ISBN existe en la respuesta
	key := "ISBN:" + isbn
	book, found := result[key]
	if !found {
		return nil, fmt.Errorf("libro con ISBN %s no encontrado", isbn)
	}

	return &book, nil
}
// GuardarLibroPorISBN maneja la solicitud HTTP para buscar y guardar un libro por ISBN
func guardarLibroPorISBN(w http.ResponseWriter, usuarioID interface{}, nombre string, portada string, autor string, isbn string, editorial string) error {
    conectar.Conectar()
    // Preparar la consulta SQL para insertar los datos del libro en la base de datos
    sql := `INSERT INTO todosloslibros (UsuarioId, nombre, portada, autor, isbn, editorial, fecha_creacion)
    VALUES (?, ?, ?, ?, ?, ?, ?);`
    _, err := conectar.Db.Exec(sql,
        usuarioID,
        nombre,
        portada,
        autor,
        isbn,
        editorial,
        time.Now(),
    )
    if err != nil {
        return fmt.Errorf("Error al guardar el libro en la base de datos: %w", err)
    }
    return nil
}
// Respuesta del formulario. Busca el libro con FetchBookInfo y lo incluye en la BD con GuardarLibroPorISBN
func BuscarLibroPorISBN_post(w http.ResponseWriter, r *http.Request) {
    isbn := r.FormValue("isbn")
    book, err := fetchBookInfo(isbn)
    if err != nil {
		Pagina404(w, r, "Error al buscar el ISBN: "+err.Error())
        return
    }

    // url := book.URL
    title := book.Title
    var authorName string
    if len(book.Authors) > 0 {
        authorName = book.Authors[0].Name
    }
    var publisherName string
    if len(book.Publishers) > 0 {
        publisherName = book.Publishers[0].Name
    }else{
		publisherName = "No especifica"
	}
    coverMedium := book.Cover.Medium

	// Obtener el ID del usuario desde los datos de seguridad
	usuarioID := DatosDeSeguridad(w, r)["usuario_id"]

	// Usar la nueva firma de la función
    err = guardarLibroPorISBN(w, usuarioID, title, coverMedium, authorName, isbn, publisherName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Crear un mensaje flash y redirigir al usuario
    utilidades.CrearMensajesFlash(w, r, "success", "Se creó el registro exitosamente")
    http.Redirect(w, r, "/registro-libros/isbn", http.StatusSeeOther)
}




