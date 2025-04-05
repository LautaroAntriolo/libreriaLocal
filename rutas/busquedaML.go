package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func LibroPorMercadoLibre(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/libreria/por_url_ml.html",
		"templates/principales/navbarLogin.html",
		"templates/principales/navbarNotLogin.html",
		utilidades.Frontend,
	))

	data := DatosDeSeguridad(w, r)
	err := tmpl.ExecuteTemplate(w, "por_url_ml.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LibroPorMercadoLibre_post(w http.ResponseWriter, r *http.Request) {
    conectar.Conectar()
    defer conectar.CerrarConexion()

    tmpl := template.Must(template.ParseFiles(
        "templates/libreria/por_url_ml.html",
        "templates/principales/navbarLogin.html",
        "templates/principales/navbarNotLogin.html",
        utilidades.Frontend,
    ))

    url := r.FormValue("url")
    if url == "" {
        http.Error(w, "URL no proporcionada", http.StatusBadRequest)
        return
    }

    bookData, err := getBookData(url)
    if err != nil {
        log.Printf("Error al obtener datos del libro: %v", err)
        http.Error(w, "Error al obtener datos del libro", http.StatusInternalServerError)
        return
    }

    // Obtener el ID del usuario desde la sesión
    data := DatosDeSeguridad(w, r)
    usuarioID := data["usuario_id"] // Asegúrate de que esto coincide con cómo guardas el ID en tu sistema

    // Guardar el libro en la base de datos
    if err := guardarLibroPorMercadoLibre(w, usuarioID, bookData); err != nil {
        log.Printf("Error al guardar el libro: %v", err)
        http.Error(w, "Error al guardar el libro", http.StatusInternalServerError)
        return
    }

    data["Libro"] = bookData
    data["Mensaje"] = "Libro guardado exitosamente"

    err = tmpl.ExecuteTemplate(w, "por_url_ml.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func getBookData(url string) (modelos.BookData, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(url)
	if err != nil {
		return modelos.BookData{}, fmt.Errorf("error al hacer la solicitud HTTP: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return modelos.BookData{}, fmt.Errorf("código de estado %d: %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return modelos.BookData{}, fmt.Errorf("error al parsear HTML: %v", err)
	}

	book := modelos.BookData{
		URL: url,
	}

	// Extraer título principal
	book.Title = strings.TrimSpace(doc.Find("h1.ui-pdp-title").Text())

	// Buscar en la sección de características destacadas
	doc.Find("section.ui-vpp-highlighted-specs table.andes-table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr.andes-table__row").Each(func(j int, row *goquery.Selection) {
			header := strings.TrimSpace(row.Find("th.andes-table__header").Text())
			value := strings.TrimSpace(row.Find("td.andes-table__column").Text())

			switch {
			case strings.Contains(header, "Título") || strings.Contains(header, "Nombre"):
				book.Title = value
			case strings.Contains(header, "Autor"):
				authors := strings.Split(value, ",")
				for _, a := range authors {
					book.Authors = append(book.Authors, modelos.Author{Name: strings.TrimSpace(a)})
				}
			case strings.Contains(header, "Editorial") || strings.Contains(header, "Sello"):
				publishers := strings.Split(value, ",")
				for _, p := range publishers {
					book.Publishers = append(book.Publishers, modelos.Publisher{Name: strings.TrimSpace(p)})
				}
			case strings.Contains(header, "ISBN"):
				book.Isbn = value
			case strings.Contains(header, "Número de páginas") || strings.Contains(header, "Páginas"):
				pages, err := strconv.Atoi(strings.TrimSpace(value))
				if err == nil {
					book.NumberOfPages = pages
				}
			}
		})
	})

	// Extraer imagen de portada
	imgSrc, exists := doc.Find(".ui-pdp-gallery__figure img").First().Attr("src")
	if exists {
		book.Cover = modelos.Cover{
			Small:  imgSrc,
			Medium: imgSrc,
			Large:  imgSrc,
		}
	}

	return book, nil
}
func guardarLibroPorMercadoLibre(w http.ResponseWriter, usuarioID interface{}, libro modelos.BookData) error {
    conectar.Conectar()
    defer conectar.CerrarConexion()

    // Convertir slice de autores a string separado por comas
    var autoresStr strings.Builder
    for i, autor := range libro.Authors {
        if i > 0 {
            autoresStr.WriteString(", ")
        }
        autoresStr.WriteString(autor.Name)
    }

    // Convertir slice de editoriales a string separado por comas
    var editorialesStr strings.Builder
    for i, editorial := range libro.Publishers {
        if i > 0 {
            editorialesStr.WriteString(", ")
        }
        editorialesStr.WriteString(editorial.Name)
    }

    // Preparar la consulta SQL para insertar los datos del libro
    sql := `INSERT INTO todosloslibros 
            (UsuarioId, nombre, portada, autor, isbn, editorial, fecha_creacion, pagina)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

    _, err := conectar.Db.Exec(sql,
        usuarioID,
        libro.Title,
        libro.Cover.Small, // Usamos la imagen pequeña como portada
        autoresStr.String(),
        libro.Isbn,
        editorialesStr.String(),
        time.Now(),
        libro.NumberOfPages,
    )

    if err != nil {
        return fmt.Errorf("error al guardar el libro en la base de datos: %w", err)
    }
    return nil
}